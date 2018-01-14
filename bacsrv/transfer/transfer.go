package transfer

import (
	"encoding/gob"
	"io"
	"net"
	"time"

	"github.com/damekr/backer/bacsrv/storage"
	"github.com/damekr/backer/common"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "transfer"})

type MainSession struct {
	ConnParams *common.ConnParameters
	Conn       net.Conn
	Id         uint64
	Transfer   *common.Transfer
	Storage    storage.Storage
	Metadata   SessionMetaData
}

type SessionMetaData struct {
	ClientName    string `json:"clientName"`
	BackupID      int    `json:"backupID"`
	BucketPath    string `json:"bucketLocation"`
	SavesetPath   string `json:"savesetLocation"`
	FilesMetadata []FileMetaData
}

type FileMetaData struct {
	FileWithPath string `json:"fileWithPath"`
	BackupTime   string `json:"backupTime"`
}

func NewSession(id uint64, params *common.ConnParameters, conn net.Conn, storage storage.Storage) *MainSession {
	return &MainSession{
		Id:         id,
		ConnParams: params,
		Conn:       conn,
		Storage:    storage,
	}
}

func (s *MainSession) Negotiate(protoVersion string) error {
	log.Println("Starting negotiate with client")
	neg := new(common.Negotiate)
	dec := gob.NewDecoder(s.Conn)
	err := dec.Decode(&neg)
	if err != nil {
		log.Println("Could not decode negotiation struct")
		return err
	}
	log.Println("Got protocol revision: ", neg.ProtoVersion)
	if neg.ProtoVersion != protoVersion {
		return common.ProtocolVersionMismatch
	}
	negc := new(common.Negotiate)
	negc.ProtoVersion = protoVersion
	enc := gob.NewEncoder(s.Conn)
	err = enc.Encode(&negc)
	if err != nil {
		log.Println("Could not encode negotiation struct")
		return err
	}
	s.Metadata.ClientName = neg.ClientName
	return nil
}

func (s *MainSession) Authenticate(password string) error {
	log.Println("Starting authentication")
	auth := new(common.Authenticate)
	dec := gob.NewDecoder(s.Conn)
	err := dec.Decode(&auth)
	if err != nil {
		log.Println("Could not decode authentication struct")
		return err
	}
	authRespond := new(common.Authenticate)
	receivedPass, err := common.Decrypt(auth.CiperText, []byte(common.KEY))
	if password != string(receivedPass) {
		log.Println("Authentication failed")
		authRespond.CiperText = []byte(common.AuthenticationFailed)
		dec := gob.NewEncoder(s.Conn)
		err = dec.Encode(&authRespond)
		return common.AuthenticationFailedError
	}
	authRespond.CiperText = []byte(common.AuthenticationPassed)
	decRespond := gob.NewEncoder(s.Conn)
	err = decRespond.Encode(&authRespond)
	log.Println("Authentication passed!")
	return nil
}

func (s *MainSession) SessionDispatcher() error {
	//Receiving type of session (TPUT, TGET)

	log.Println("Dispatching incoming connection")
	transfer := new(common.Transfer)
	tr := gob.NewDecoder(s.Conn)
	err := tr.Decode(&transfer)
	if err != nil {
		if err == io.EOF {
			log.Println("Client closed connection")
			return err
		}
		log.Println("Cannot decode incoming transfer type, responding with empty struct. Error: ", err)
	}
	log.Println("Got incoming transfer type connection: ", transfer.TransferType)
	s.Transfer = transfer
	switch transfer.TransferType {
	case common.TGET:
		transfer.Buffer = common.BUFFERSIZE
		trc := gob.NewEncoder(s.Conn)
		err = trc.Encode(&transfer)
		if err != nil {
			log.Println("Could not send transfer type response to client")
		}
		restoreSession := CreateRestoreSession(s)
		log.Debugln("Handling restore session")
		return restoreSession.HandleRestoreSession()

	case common.TPUT:
		transfer.Buffer = common.BUFFERSIZE
		trp := gob.NewEncoder(s.Conn)
		err = trp.Encode(&transfer)
		if err != nil {
			log.Println("Could not send transfer type response to client")
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
		trEmpty := new(common.Transfer)
		trE := gob.NewEncoder(s.Conn)
		err = trE.Encode(&trEmpty)
		if err != nil {
			log.Println("Could not encode empty struct in non handled transfer type")
			return err
		}
	}

	return nil
}
