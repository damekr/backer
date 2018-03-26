package transfer

import (
	"encoding/gob"

	"net"

	"github.com/damekr/backer/cmd/baclnt/config"
	"github.com/damekr/backer/cmd/baclnt/fs"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "transfer"})

type MainSession struct {
	ConnParams *bftp.ConnParameters
	Conn       net.Conn
	Id         uint64
	Transfer   *bftp.Transfer
	Storage    *fs.LocalFileSystem
}

type RestoreFileMetadata struct {
	PathOnServer string
	PathOnClient string
}

func NewSession(id uint64, params *bftp.ConnParameters, conn net.Conn) *MainSession {
	return &MainSession{
		Id:         id,
		ConnParams: params,
		Conn:       conn,
	}
}

func (s *MainSession) StartBackup(filesPaths []string) error {
	backupTransfer := bftp.Transfer{
		TransferType:  bftp.TPUT,
		ObjectsNumber: len(filesPaths),
	}
	// Sending transfer type
	log.Println("Opening transfer with server, type: ", backupTransfer.TransferType)
	tr := gob.NewEncoder(s.Conn)
	err := tr.Encode(&backupTransfer)
	if err != nil {
		log.Println("Could not encode transfer struct, error: ", err)
		return err
	}
	// Acknowledge message
	log.Println("Waiting for response...")
	tranType := new(bftp.Transfer)
	trInc := gob.NewDecoder(s.Conn)
	err = trInc.Decode(&tranType)
	if err != nil {
		log.Println("Could not decode response of transfer type")
		return err
	}
	// Server response an empty string when does not support requesting operation
	if tranType.TransferType == "" {
		log.Println("Server does not support such operation")
		return bftp.ServerDoesNotSupportSuchOperation
	}
	log.Println("Server accepts connection type")
	s.Transfer = tranType

	// Creating new filesystem object to handle backup session - TODO Consider put it in main session
	fileSystem := fs.NewLocalFileSystem()

	backupSession := CreateBackupSession(s, fileSystem)
	for fileNumber, path := range filesPaths {
		log.Debugf("Sending file %s, number %s", path, fileNumber)
		backupSession.PutFile(path, path)
	}
	return nil
}

func (s *MainSession) StartRestore(restoreFileMetadata []RestoreFileMetadata) error {
	restoreTransfer := bftp.Transfer{
		TransferType:  bftp.TGET,
		ObjectsNumber: len(restoreFileMetadata),
	}
	log.Println("Opening transfer with server, type: ", restoreTransfer.TransferType)
	tr := gob.NewEncoder(s.Conn)
	err := tr.Encode(&restoreTransfer)
	if err != nil {
		log.Println("Could not encode transfer struct, error: ", err)
		return err
	}
	log.Println("Waiting for response...")
	tranType := new(bftp.Transfer)
	trInc := gob.NewDecoder(s.Conn)
	err = trInc.Decode(&tranType)
	if err != nil {
		log.Println("Could not decode response of transfer type")
		return err
	}
	// Server response an empty string when does not support requesting operation
	if tranType.TransferType == "" {
		log.Println("Server does not support such operation")
		return bftp.ServerDoesNotSupportSuchOperation
	}
	log.Println("Server accepts connection type")
	s.Transfer = tranType
	log.Debugln("Creating restore session")

	// Creating new filesystem object to handle backup session - TODO consider place it in main session
	fileSystem := fs.NewLocalFileSystem()

	restoreSession := CreateRestoreSession(s, fileSystem)
	for _, v := range restoreFileMetadata {
		log.Debugf("Downloading file from server path: %s, to local path: %s", v.PathOnServer, v.PathOnServer)
		err = restoreSession.GetFile(v.PathOnServer, v.PathOnClient)
		if err != nil {
			log.Errorln("Cannot download file, err: ", err)
		}
		log.Debugf("File %s has been downloaded", v.PathOnClient)
	}

	return nil
}

func (s *MainSession) CloseSession() error {
	log.Debugln("Closing connection...")
	return s.Conn.Close()
}

func (s *MainSession) Negotiate(protoVersion string) error {
	log.Print("Starting Negotiate with server")
	neg := new(bftp.Negotiate)
	neg.ProtoVersion = protoVersion
	neg.ClientName = config.MainConfig.ExternalName
	enc := gob.NewEncoder(s.Conn)
	err := enc.Encode(&neg)
	if err != nil {
		log.Println("Could not encode negotiation struct")
		return err
	}
	negs := new(bftp.Negotiate)
	dec := gob.NewDecoder(s.Conn)
	err = dec.Decode(&negs)
	if err != nil {
		log.Println("Could not decode negotiation struct, error: ", err)
		return err
	}

	if negs.ProtoVersion != protoVersion {
		log.Println("Server sent protocol version: ", negs.ProtoVersion)
		return bftp.ProtocolVersionMismatch
	}
	log.Println("Got protocol revision: ", negs.ProtoVersion)
	return nil
}

func (s *MainSession) Authenticate(password string) error {
	log.Print("Starting authentication")
	auth := new(bftp.Authenticate)
	cipherText, err := bftp.Encrypt([]byte(password), []byte(bftp.KEY))
	if err != nil {
		log.Println("Could not encrypt and get cipher, error: ", err)
		return err
	}
	auth.CiperText = cipherText
	dec := gob.NewEncoder(s.Conn)
	err = dec.Encode(&auth)
	if err != nil {
		log.Println("Could not decode authentication struct")
		return err
	}
	authRespond := new(bftp.Authenticate)
	decRespond := gob.NewDecoder(s.Conn)
	err = decRespond.Decode(&authRespond)
	if string(authRespond.CiperText) != bftp.AuthenticationPassed {
		log.Println("Authentication failed!")
		return bftp.AuthenticationFailedError
	}
	log.Println("Authentication passed!")
	return nil
}
