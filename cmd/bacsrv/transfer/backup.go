package transfer

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/damekr/backer/cmd/bacsrv/db"
	"github.com/damekr/backer/cmd/bacsrv/storage"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/sirupsen/logrus"
)

type BackupSession struct {
	MainSession *MainSession
	Database    db.DB
}

var logBackup = logrus.WithFields(logrus.Fields{"prefix": "transfer:backup"})

func CreateBackupSession(mainSession *MainSession) *BackupSession {
	return &BackupSession{
		MainSession: mainSession,
		Database:    db.Get(),
	}
}

func (b *BackupSession) HandleBackupSession(savesetLocation string, objectsNumber int) error {
	logBackup.Debugln("Handling incoming TPUT transfer type")
	for i := 0; i < objectsNumber; i++ {
		logBackup.Debugln("Receiving object: ", i)
		// Getting file metadata
		fileMetadata, err := b.receiveFileMetadata()
		if err != nil {
			logBackup.Errorln("Could not decode file metadata, err: ", err.Error())
		}
		log.Debugln("Received metadata: ", fileMetadata)
		// Sending acknowledge
		// TODO Make checks like: disk space
		err = b.sendFileMetaDataAcknowledge(fileMetadata)
		if err != nil {
			logBackup.Errorln("Could not send file metadata as an acknowledge")
		}

		// Downloading file
		err = b.downloadFile(*fileMetadata, savesetLocation)
		if err != nil {
			logBackup.Errorln("Cannot upload file, err: ", err.Error())
			return err
		}
		logBackup.Debugln("Received file, sending acknowledge")

		//Sending file acknowledge
		fileSize := storage.GetFileSize(fileMetadata.FullPath)
		fileSizeAckn := new(bftp.FileAcknowledge)
		fileSizeAckn.Size = fileSize
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
	logBackup.Debugln("Starting downloading file to path: ", fileMetadata.FullPath)
	file, err := b.MainSession.Storage.CreateFile(savesetLocation, fileMetadata.FullPath)
	if err != nil {
		logBackup.Errorln("Cannot create localfile to write")
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
	fileMetadata.LocationOnServer = filepath.Join(savesetLocation, fileMetadata.FullPath)
	fileMetadata.OriginalFileLocation = fileMetadata.FullPath
	timeFinishBackup := time.Since(timeStartBackup)
	fileMetadata.BackupTime = timeFinishBackup.String()
	b.MainSession.Metadata.FilesMetadata = append(b.MainSession.Metadata.FilesMetadata, fileMetadata)
	logBackup.Infoln("Backup duration: ", timeFinishBackup.String())
	fileSizeInMb := wroteToFile / 1000 / 1000
	logBackup.Infof("Average speed: %f MiB/s", float64(fileSizeInMb)/timeFinishBackup.Seconds())
	return nil
}

func (b *BackupSession) createMetadata() error {
	jsonData, err := json.Marshal(b.MainSession.Metadata)
	if err != nil {
		return err
	} else {
		if err := b.Database.WriteBackupMetadata(jsonData, filepath.Base(b.MainSession.Metadata.SavesetPath), b.MainSession.Metadata.ClientName); err != nil {
			return err
		}
	}
	return nil
}
