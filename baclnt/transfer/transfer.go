package transfer

import (
	"encoding/gob"

	"net"

	"github.com/damekr/backer/baclnt/config"
	"github.com/damekr/backer/baclnt/fs"
	"github.com/damekr/backer/common"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "transfer"})

type MainSession struct {
	ConnParams *common.ConnParameters
	Conn       net.Conn
	Id         uint64
	Transfer   *common.Transfer
	Storage    *fs.FileSystem
}

func NewSession(id uint64, params *common.ConnParameters, conn net.Conn) *MainSession {
	return &MainSession{
		Id:         id,
		ConnParams: params,
		Conn:       conn,
	}
}

func (s *MainSession) CloseSession() error {
	return s.Conn.Close()
}

func (s *MainSession) Negotiate(protoVersion string) error {
	log.Print("Starting Negotiate with server")
	neg := new(common.Negotiate)
	neg.ProtoVersion = protoVersion
	neg.ClientName = config.MainConfig.ExternalName
	enc := gob.NewEncoder(s.Conn)
	err := enc.Encode(&neg)
	if err != nil {
		log.Println("Could not encode negotiation struct")
		return err
	}
	negs := new(common.Negotiate)
	dec := gob.NewDecoder(s.Conn)
	err = dec.Decode(&negs)
	if err != nil {
		log.Println("Could not decode negotiation struct, error: ", err)
		return err
	}

	if negs.ProtoVersion != protoVersion {
		log.Println("Server sent protocol version: ", negs.ProtoVersion)
		return common.ProtocolVersionMismatch
	}
	log.Println("Got protocol revision: ", negs.ProtoVersion)
	return nil
}

func (s *MainSession) Authenticate(password string) error {
	log.Print("Starting authentication")
	auth := new(common.Authenticate)
	cipherText, err := common.Encrypt([]byte(password), []byte(common.KEY))
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
	authRespond := new(common.Authenticate)
	decRespond := gob.NewDecoder(s.Conn)
	err = decRespond.Decode(&authRespond)
	if string(authRespond.CiperText) != common.AuthenticationPassed {
		log.Println("Authentication failed!")
		return common.AuthenticationFailedError
	}
	log.Println("Authentication passed!")
	return nil
}

func (s *MainSession) StartBackup(filesPaths []string) error {
	backupTransfer := common.Transfer{
		TransferType:  common.TPUT,
		ObjectsNumber: len(filesPaths),
	}
	//Sending transfer type
	log.Println("Opening transfer with server, type: ", backupTransfer.TransferType)
	tr := gob.NewEncoder(s.Conn)
	err := tr.Encode(&backupTransfer)
	if err != nil {
		log.Println("Could not encode transfer struct, error: ", err)
		return err
	}
	//Acknowledge
	log.Println("Waiting for response...")
	tranType := new(common.Transfer)
	trInc := gob.NewDecoder(s.Conn)
	err = trInc.Decode(&tranType)
	if err != nil {
		log.Println("Could not decode response of transfer type")
		return err
	}
	//Server response an empty string when does not support requesting operation
	if tranType.TransferType == "" {
		log.Println("Server does not support such operation")
		return common.ServerDoesNotSupportSuchOperation
	}
	log.Println("Server accepts connection type")
	s.Transfer = tranType
	backupSession := CreateBackupSession(s)
	for fileNumber, path := range filesPaths {
		log.Debugf("Sending file %s, number %s", path, fileNumber)
		backupSession.PutFile(path, path)
	}

	return nil
}

func (s *MainSession) StartRestore(filesPaths []string) error {
	restoreTransfer := common.Transfer{
		TransferType: common.TGET,
	}
	log.Println("Opening transfer with server, type: ", restoreTransfer.TransferType)
	tr := gob.NewEncoder(s.Conn)
	err := tr.Encode(&restoreTransfer)
	if err != nil {
		log.Println("Could not encode transfer struct, error: ", err)
		return err
	}
	log.Println("Waiting for response...")
	tranType := new(common.Transfer)
	trInc := gob.NewDecoder(s.Conn)
	err = trInc.Decode(&tranType)
	if err != nil {
		log.Println("Could not decode response of transfer type")
		return err
	}
	//Server response an empty string when does not support requesting operation
	if tranType.TransferType == "" {
		log.Println("Server does not support such operation")
		return common.ServerDoesNotSupportSuchOperation
	}
	log.Println("Server accepts connection type")
	s.Transfer = tranType

	return nil
}
