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

func NewSession(id uint64, params *bftp.ConnParameters, conn net.Conn) *MainSession {
	return &MainSession{
		Id:         id,
		ConnParams: params,
		Conn:       conn,
	}
}

func (s *MainSession) StartBackup(backupObjects fs.BackupObjects) error {
	backupTransfer := &bftp.Transfer{
		TransferType:  bftp.TPUT,
		ObjectsNumber: len(backupObjects.Files),
	}
	err := s.sendTransferType(backupTransfer)
	if err != nil {
		log.Errorln("Cannot send transfer type struct, err: ", err)
	}

	// Acknowledge message
	ack, err := s.receiveTransferTypeAck()
	if err != nil {
		log.Errorln("Cannot receive ack message about transfer type, err: ", err)
		return err
	}

	// Server response an empty string when does not support requesting operation
	if ack.TransferType == "" {
		log.Println("Server does not support such operation")
		return bftp.ServerDoesNotSupportSuchOperation
	}

	log.Println("Server accepts connection type")
	s.Transfer = ack

	// Creating new filesystem object to handle backup session - TODO Consider put it in main session
	fileSystem := fs.NewLocalFileSystem()

	backupSession := CreateBackupSession(s, fileSystem)

	// Sending full dirs structure
	err = backupSession.sendDirsMetadata(backupObjects.Dirs)
	if err != nil {
		return err
	}

	// Receiving acknowledge
	err = backupSession.receiveEmptyAckMessage()
	if err != nil {
		return err
	}

	for fileNumber, path := range backupObjects.Files {
		log.Debugf("Sending file %s, number %s", path, fileNumber)
		backupSession.putFile(path, path)
	}
	return nil
}

func (s *MainSession) StartRestore(assetID int, options bftp.RestoreOptions) error {
	restoreTransfer := &bftp.Transfer{
		TransferType: bftp.TGET,
		AssetID:      assetID,
	}

	log.Println("Opening transfer with server, type: ", restoreTransfer.TransferType)
	err := s.sendTransferType(restoreTransfer)
	if err != nil {
		log.Errorln("Cannot send transfer type, err: ", err)
		return err
	}

	log.Println("Waiting for response...")
	ack, err := s.receiveTransferTypeAck()
	if err != nil {
		log.Errorln("Cannot receive ack message about transfer type, err: ", err)
		return err
	}

	// Server response an empty string when does not support requesting operation
	if ack.TransferType == "" {
		log.Println("Server does not support such operation")
		return bftp.ServerDoesNotSupportSuchOperation
	}

	log.Println("Server accepts connection type")
	s.Transfer = ack
	log.Debugln("Creating restore session")

	// Creating new filesystem object to handle backup session - TODO consider place it in main session
	fileSystem := fs.NewLocalFileSystem()

	restoreSession := CreateRestoreSession(s, fileSystem)

	// Getting assets metadata
	assetsMetadata, err := restoreSession.receiveAssetMetadata(options)
	if err != nil {
		log.Errorln("Cannot receive asset metadata")
	}
	log.Println("ASSET METADATA: ", assetsMetadata)

	//// Creating local dirs bases on metadata
	//err = restoreSession.createDirs(restoreMetadata.DirsMetadata)
	//if err != nil {
	//	log.Errorln("At least one error occured during creating dirs, err: ", err)
	//}
	//
	//for _, v := range restoreMetadata.FilesMetadata {
	//	log.Debugf("Downloading file from server path: %s, to local path: %s", v.LocationOnServer, v.FullPath)
	//	err = restoreSession.GetFile(v.LocationOnServer, v.FullPath)
	//	if err != nil {
	//		log.Errorln("Cannot download file, err: ", err)
	//	}
	//	log.Debugf("File %s has been downloaded", v.FullPath)
	//}

	return nil
}

func (s *MainSession) sendTransferType(transferType *bftp.Transfer) error {
	// Sending transfer type and number of objects
	log.Debugln("Opening transfer with server, type: ", transferType.TransferType)
	tr := gob.NewEncoder(s.Conn)
	err := tr.Encode(transferType)
	if err != nil {
		return err
	}
	return nil
}

func (s *MainSession) receiveTransferTypeAck() (*bftp.Transfer, error) {
	log.Println("Waiting for response...")
	tranType := new(bftp.Transfer)
	trInc := gob.NewDecoder(s.Conn)
	err := trInc.Decode(&tranType)
	if err != nil {
		return nil, err
	}
	return tranType, nil
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
