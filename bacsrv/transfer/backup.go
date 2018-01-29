package transfer

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"path/filepath"
	"time"

	"github.com/damekr/backer/bacsrv/db"
	"github.com/damekr/backer/bacsrv/storage"
	"github.com/damekr/backer/common"
)

type BackupSession struct {
	MainSession *MainSession
	Database    db.DB
}

func CreateBackupSession(mainSession *MainSession) *BackupSession {
	return &BackupSession{
		MainSession: mainSession,
		Database:    db.Get(),
	}
}

func (b *BackupSession) receiveFileMetadata() (*common.FileMetadata, error) {
	fileMetadata := new(common.FileMetadata)
	fileTEmpty := new(common.FileMetadata)
	fileTDec := gob.NewDecoder(b.MainSession.Conn)
	err := fileTDec.Decode(&fileMetadata)
	if err != nil {
		log.Print("Could not decode FileMetadata struct, error: ", err)
		fileTEnc := gob.NewEncoder(b.MainSession.Conn)
		if err := fileTEnc.Encode(&fileTEmpty); err != nil {
			log.Println("Could not encode empty FileMetadata struct")
			return nil, err
		}
	}
	return fileMetadata, nil
}

func (b *BackupSession) sendFileMetaDataAcknowledge(metadata *common.FileMetadata) error {
	fileAEnc := gob.NewEncoder(b.MainSession.Conn)
	if err := fileAEnc.Encode(&metadata); err != nil {
		log.Println("Could not send acknowledge")
		return err
	}
	return nil
}

func (b *BackupSession) sendFileTransferAcknowledge(acknowledge *common.FileAcknowledge) error {
	fileSizeEncoder := gob.NewEncoder(b.MainSession.Conn)
	if err := fileSizeEncoder.Encode(&acknowledge); err != nil {
		log.Println("Could not encode FileAcknowledge struct")
		return err
	}
	log.Println("Sent file size acknowledge")
	return nil
}

func (b *BackupSession) HandleBackupSession(savesetLocation string, objectsNumber int) error {
	log.Println("Handling incoming TPUT transfer type")
	for i := 0; i < objectsNumber; i++ {
		log.Debugln("Receiving object: ", i)
		// Getting file metadata
		fileMetadata, err := b.receiveFileMetadata()
		if err != nil {
			log.Errorln("Could not decode file metadata, err: ", err.Error())
		}

		// Sending acknowledge
		// TODO Make checks like: disk space
		err = b.sendFileMetaDataAcknowledge(fileMetadata)
		if err != nil {
			log.Errorln("Could not send file metadata as an acknowledge")
		}

		// Downloading file
		err = b.downloadFile(fileMetadata.Name, fileMetadata.FullPath, fileMetadata.FileSize, savesetLocation)
		if err != nil {
			log.Println("Cannot upload file, err: ", err.Error())
			return err
		}
		log.Println("Received file, sending acknowledge")

		//Sending file acknowledge
		fileSize := storage.GetFileSize(fileMetadata.FullPath)
		fileSizeAckn := new(common.FileAcknowledge)
		fileSizeAckn.Size = fileSize
		err = b.sendFileTransferAcknowledge(fileSizeAckn)
		if err != nil {
			log.Error(err)
		}
	}
	if err := b.createMetadata(); err != nil {
		log.Errorln("Could not create session metadata, err: ", err)
	}

	return nil
}

func (b *BackupSession) downloadFile(name, localFilePath string, size int64, savesetLocation string) error {
	log.Println("Starting downloading file to path: ", localFilePath)
	file, err := b.MainSession.Storage.CreateFile(savesetLocation, localFilePath)
	if err != nil {
		log.Println("Cannot create localfile to write")
		return err
		//	TODO Respond with failed transfer, error on server side
	}
	defer file.Close()
	fileMetadata := FileMetaData{
		FileWithPath: localFilePath,
	}
	writer := bufio.NewWriter(file)
	reader := bufio.NewReader(b.MainSession.Conn)
	var readFromConnection int64
	var wroteToFile int64
	timeStartBackup := time.Now()
	buffer := make([]byte, b.MainSession.Transfer.Buffer)
	if size < int64(b.MainSession.Transfer.Buffer) {
		log.Println("Shrinking buffer to filesize: ", size)
		buffer = make([]byte, size)
	}
	for {
		read, err := reader.Read(buffer)
		if err != nil {
			log.Println("Could not read from connection reader, error: ", err)
			break
		}
		readFromConnection += int64(read)
		wrote, err := writer.Write(buffer[:read])
		if err != nil {
			log.Println("Could not write to file writter: ", err)
			break
		}
		wroteToFile += int64(wrote)
		if wroteToFile == size {
			log.Println("Wrote all data to file")
			break
		}
	}
	writer.Flush()
	timeFinishBackup := time.Since(timeStartBackup)
	fileMetadata.BackupTime = timeFinishBackup.String()
	b.MainSession.Metadata.FilesMetadata = append(b.MainSession.Metadata.FilesMetadata, fileMetadata)
	log.Println("Backup duration: ", timeFinishBackup.String())
	fileSizeInMb := wroteToFile / 1000 / 1000
	log.Printf("Average speed: %f MiB/s", float64(fileSizeInMb)/timeFinishBackup.Seconds())
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
