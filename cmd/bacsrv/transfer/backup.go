package transfer

import (
	"bufio"
	"encoding/gob"
	"path/filepath"
	"time"

	"github.com/damekr/backer/cmd/bacsrv/db"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/sirupsen/logrus"
)

type BackupSession struct {
	MainSession *MainSession
	Database    db.BackupsDB
}

var logBackup = logrus.WithFields(logrus.Fields{"prefix": "transfer:backup"})

func CreateBackupSession(mainSession *MainSession) *BackupSession {
	return &BackupSession{
		MainSession: mainSession,
		Database:    db.DB(),
	}
}

func (b *BackupSession) HandleBackupSession(savesetLocation string, objectsNumber int) error {
	logBackup.Debugln("Handling incoming TPUT transfer type")

	// Receiving dirs structure if any
	dirsMetadata, err := b.receiveDirsStructure()
	if err != nil {
		log.Errorln("Cannot receive dirs metadata structure, err: ", err)
	}
	for _, v := range *dirsMetadata {
		err = b.MainSession.Storage.CreateDir(savesetLocation, v)
		if err != nil {
			log.Errorln("Could not create dir in storage, err: ", err)
		}
	}

	// Sending acknowledge message
	err = b.MainSession.sendEmptyAckMessage()
	if err != nil {
		return err
	}

	// Saving dirs metadata
	b.MainSession.Metadata.DirsMetadata = *dirsMetadata

	for i := 0; i < objectsNumber; i++ {
		logBackup.Debugln("Receiving object: ", i)

		// Getting file metadata
		fileMetadata, err := b.receiveFileMetadata()
		if err != nil {
			logBackup.Errorln("Could not decode file metadata, err: ", err.Error())
		}
		log.Debugln("Received metadata: ", fileMetadata)
		// Sending acknowledge
		err = b.sendFileMetaDataAcknowledge(fileMetadata)
		if err != nil {
			logBackup.Errorln("Could not send file metadata as an acknowledge")
		}
		// Downloading file
		err = b.downloadFile(*fileMetadata, savesetLocation)
		if err != nil {
			logBackup.Errorln("Cannot download file, err: ", err.Error())
			return err
		}
		logBackup.Debugln("Received file, sending acknowledge")
		// Sending file acknowledge

		fileSizeAckn := new(bftp.FileAcknowledge)
		fileSizeAckn.Size = fileMetadata.FileSize
		err = b.sendFileTransferAcknowledge(fileSizeAckn)
		if err != nil {
			logBackup.Error(err)
		}
	}
	if err := b.createMetadata(); err != nil {
		logBackup.Errorln("Could not create session metadata, err: ", err)
	}

	return nil
}

func (b *BackupSession) receiveDirsStructure() (*[]bftp.DirMetadata, error) {
	log.Debugln("Receiving dir structure")
	dirsStructure := new([]bftp.DirMetadata)
	dirsMetadataDecoder := gob.NewDecoder(b.MainSession.Conn)
	err := dirsMetadataDecoder.Decode(&dirsStructure)
	if err != nil {
		return nil, err
	}
	log.Debugln("Received dirs metadata struct: ", dirsStructure)
	return dirsStructure, nil
}

func (b *BackupSession) receiveFileMetadata() (*bftp.FileMetadata, error) {
	fileMetadata := new(bftp.FileMetadata)
	fileTEmpty := new(bftp.FileMetadata)
	fileTDec := gob.NewDecoder(b.MainSession.Conn)
	err := fileTDec.Decode(&fileMetadata)
	if err != nil {
		logBackup.Errorln("Could not decode FileMetadata struct, error: ", err)
		fileTEnc := gob.NewEncoder(b.MainSession.Conn)
		if err := fileTEnc.Encode(&fileTEmpty); err != nil {
			logBackup.Errorln("Could not encode empty FileMetadata struct")
			return nil, err
		}
	}
	return fileMetadata, nil
}

func (b *BackupSession) sendFileMetaDataAcknowledge(metadata *bftp.FileMetadata) error {
	fileAEnc := gob.NewEncoder(b.MainSession.Conn)
	if err := fileAEnc.Encode(&metadata); err != nil {
		logBackup.Errorln("Could not send acknowledge")
		return err
	}
	return nil
}

func (b *BackupSession) sendFileTransferAcknowledge(acknowledge *bftp.FileAcknowledge) error {
	fileSizeEncoder := gob.NewEncoder(b.MainSession.Conn)
	if err := fileSizeEncoder.Encode(&acknowledge); err != nil {
		logBackup.Errorln("Could not encode FileAcknowledge struct")
		return err
	}
	logBackup.Println("Sent file size acknowledge")
	return nil
}

func (b *BackupSession) downloadFile(fileMetadata bftp.FileMetadata, savesetLocation string) error {
	logBackup.Debugln("Starting downloading file:", fileMetadata.NameWithPath)
	file, err := b.MainSession.Storage.CreateFile(savesetLocation, fileMetadata.NameWithPath)
	if err != nil {
		logBackup.Errorln("Cannot create localfile to write, err: ", err)
		return err
		//	TODO Respond with failed transfer, error on server side
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	reader := bufio.NewReader(b.MainSession.Conn)
	var readFromConnection int64
	var wroteToFile int64
	timeStartBackup := time.Now()
	buffer := make([]byte, b.MainSession.Transfer.Buffer)
	if fileMetadata.FileSize < int64(b.MainSession.Transfer.Buffer) {
		logBackup.Debugln("Shrinking buffer to filesize: ", fileMetadata.FileSize)
		buffer = make([]byte, fileMetadata.FileSize)
	}
	for {
		read, err := reader.Read(buffer)
		if err != nil {
			logBackup.Errorln("Could not read from connection reader, error: ", err)
			break
		}
		readFromConnection += int64(read)
		wrote, err := writer.Write(buffer[:read])
		if err != nil {
			logBackup.Errorln("Could not write to file writer: ", err)
			break
		}
		wroteToFile += int64(wrote)
		if wroteToFile == fileMetadata.FileSize {
			logBackup.Debugln("Wrote all data to file")
			break
		}
	}
	writer.Flush()
	fileMetadata.LocationOnServer = filepath.Join(savesetLocation, fileMetadata.NameWithPath)
	timeFinishBackup := time.Since(timeStartBackup)
	fileMetadata.BackupTime = timeFinishBackup.String()
	b.MainSession.Metadata.FilesMetadata = append(b.MainSession.Metadata.FilesMetadata, fileMetadata)
	logBackup.Infoln("Backup duration: ", timeFinishBackup.String())
	fileSizeInMb := wroteToFile / 1000 / 1000
	logBackup.Infof("Average speed: %f MiB/s", float64(fileSizeInMb)/timeFinishBackup.Seconds())
	return nil
}

func (b *BackupSession) createMetadata() error {
	if err := b.Database.CreateAssetMetadata(b.MainSession.Metadata); err != nil {
		return err
	}
	return nil
}
