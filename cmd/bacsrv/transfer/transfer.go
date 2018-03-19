package transfer

import (
	"encoding/gob"
	"io"
	"net"
	"time"

	"github.com/damekr/backer/cmd/bacsrv/db"
	"github.com/damekr/backer/cmd/bacsrv/storage"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "transfer"})

type MainSession struct {
	ConnParams *bftp.ConnParameters
	Conn       net.Conn
	Id         uint64
	Transfer   *bftp.Transfer
	Storage    storage.Storage
	Metadata   SessionMetaData
}

type SessionMetaData struct {
	ClientName    string `json:"clientName"`
	BackupID      int    `json:"backupID"`
	BucketPath    string `json:"bucketLocation"`
	SavesetPath   string `json:"savesetLocation"`
	FilesMetadata []db.FileMetaData
}

func NewSession(id uint64, params *bftp.ConnParameters, conn net.Conn, storage storage.Storage) *MainSession {
	return &MainSession{
		Id:         id,
		ConnParams: params,
		Conn:       conn,
		Storage:    storage,
	}
}

func (s *MainSession) Negotiate(protoVersion string) error {
	log.Debugln("Starting negotiate with client")
	neg := new(bftp.Negotiate)
	dec := gob.NewDecoder(s.Conn)
	err := dec.Decode(&neg)
	if err != nil {
		log.Errorln("Could not decode negotiation struct")
		return err
	}
	log.Debugln("Got protocol revision: ", neg.ProtoVersion)
	if neg.ProtoVersion != protoVersion {
		return bftp.ProtocolVersionMismatch
	}
	negc := new(bftp.Negotiate)
	negc.ProtoVersion = protoVersion
	enc := gob.NewEncoder(s.Conn)
	err = enc.Encode(&negc)
	if err != nil {
		log.Errorln("Could not encode negotiation struct")
		return err
	}
	s.Metadata.ClientName = neg.ClientName
	return nil
}

func (s *MainSession) Authenticate(password string) error {
	log.Debugln("Starting authentication")
	auth := new(bftp.Authenticate)
	dec := gob.NewDecoder(s.Conn)
	err := dec.Decode(&auth)
	if err != nil {
		log.Debugln("Could not decode authentication struct")
		return err
	}
	authRespond := new(bftp.Authenticate)
	receivedPass, err := bftp.Decrypt(auth.CiperText, []byte(bftp.KEY))
	if password != string(receivedPass) {
		log.Errorln("Authentication failed")
		authRespond.CiperText = []byte(bftp.AuthenticationFailed)
		dec := gob.NewEncoder(s.Conn)
		err = dec.Encode(&authRespond)
		return bftp.AuthenticationFailedError
	}
	authRespond.CiperText = []byte(bftp.AuthenticationPassed)
	decRespond := gob.NewEncoder(s.Conn)
	err = decRespond.Encode(&authRespond)
	log.Debugln("Authentication passed!")
	return nil
}

func (s *MainSession) SessionDispatcher(createSessionMetadata bool) error {
	//Receiving type of session (TPUT, TGET)
	//TODO Might be useful to have configurable session metadata, but now it's not used. For backup create always
	log.Println("Dispatching incoming connection")
	transfer := new(bftp.Transfer)
	tr := gob.NewDecoder(s.Conn)
	err := tr.Decode(&transfer)
	if err != nil {
		if err == io.EOF {
			log.Debugln("Client closed connection")
			return err
		}
		log.Errorln("Cannot decode incoming transfer type, responding with empty struct. Error: ", err)
	}
	log.Debugln("Got incoming transfer type connection: ", transfer.TransferType)
	s.Transfer = transfer
	switch transfer.TransferType {
	case bftp.TGET:
		transfer.Buffer = bftp.BUFFERSIZE
		trc := gob.NewEncoder(s.Conn)
		err = trc.Encode(&transfer)
		if err != nil {
			log.Errorln("Could not send transfer type response to client")
		}
		restoreSession := CreateRestoreSession(s)
		log.Debugln("Handling restore session")
		return restoreSession.HandleRestoreSession(transfer.ObjectsNumber)

	case bftp.TPUT:
		transfer.Buffer = bftp.BUFFERSIZE
		trp := gob.NewEncoder(s.Conn)
		err = trp.Encode(&transfer)
		if err != nil {
			log.Errorln("Could not send transfer type response to client")
		}
		backupSession := CreateBackupSession(s)
		log.Debugln("Handling backup session")
		bucket, err := s.Storage.CreateBucket(s.Metadata.ClientName)
		if err != nil {
			log.Errorln("Could not create bucket, err: ", err.Error())
			return err
		}
		s.Metadata.BucketPath = bucket
		saveset, err := s.Storage.CreateSaveset(bucket)
		if err != nil {
			log.Errorln("Could not create saveset, err: ", err.Error())
		}
		s.Metadata.SavesetPath = saveset
		s.Metadata.BackupID = time.Now().Nanosecond()
		return backupSession.HandleBackupSession(saveset, transfer.ObjectsNumber)

	default:
		trEmpty := new(bftp.Transfer)
		trE := gob.NewEncoder(s.Conn)
		err = trE.Encode(&trEmpty)
		if err != nil {
			log.Errorln("Could not encode empty struct in non handled transfer type")
			return err
		}
	}

	return nil
}
