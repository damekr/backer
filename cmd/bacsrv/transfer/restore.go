package transfer

import (
	"bufio"
	"encoding/gob"

	"github.com/damekr/backer/cmd/bacsrv/storage"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/sirupsen/logrus"
)

type RestoreSession struct {
	MainSession *MainSession
}

var logRestore = logrus.WithFields(logrus.Fields{"prefix": "transfer:restore"})

func CreateRestoreSession(mainSession *MainSession) *RestoreSession {
	return &RestoreSession{
		MainSession: mainSession,
	}
}

//TODO Add logic to read file from specific backup(json file)

func (r *RestoreSession) HandleRestoreSession(objectsNumber int) error {
	logRestore.Debugln("Handling incomming TGET transfer type")
	for i := 0; i < objectsNumber; i++ {

		fileT := new(bftp.FileMetadata)
		fileTEmpty := new(bftp.FileMetadata)

		//Decoding file path to be transfered
		fileTDec := gob.NewDecoder(r.MainSession.Conn)
		err := fileTDec.Decode(&fileT)
		if err != nil {
			logRestore.Errorln("Could not decode FileMetadata struct, error: ", err)
			fileTEnc := gob.NewEncoder(r.MainSession.Conn)
			if err := fileTEnc.Encode(&fileTEmpty); err != nil {
				logRestore.Errorln("Could not encode empty FileMetadata struct, err: ", err)
				return err
			}
			return err
		}
		logRestore.Debugf("Checking if file %s exists", fileT.FullPath)
		if !storage.CheckIfFileExists(fileT.FullPath) {
			logRestore.Debugf("File: %s does not exist", fileT.FullPath)
			fileTEncNotExist := gob.NewEncoder(r.MainSession.Conn)
			if err := fileTEncNotExist.Encode(&fileTEmpty); err != nil {
				logRestore.Errorln("Could not encode empty FileMetadata struct, err: ", err)
				return err
			}
		}

		//Sending size of file being transfered
		fileTEnc := gob.NewEncoder(r.MainSession.Conn)
		fileT.FileSize = storage.GetFileSize(fileT.FullPath)
		if err := fileTEnc.Encode(&fileT); err != nil {
			logRestore.Errorln("Could not encode empty FileMetadata struct, err: ", err)
			return err
		}
		logRestore.Debugln("Handling transfer with sending file to client, file: ", fileT.FullPath)

		//Sending file
		err = r.uploadFile(fileT.FullPath, fileT.FileSize)
		if err != nil {
			logRestore.Errorln("Could not send file, err: ", err.Error())
		}

		//Receiving file acknowledge
		fileSize := new(bftp.FileAcknowledge)
		fileSizeEncoder := gob.NewDecoder(r.MainSession.Conn)
		if err := fileSizeEncoder.Decode(&fileSize); err != nil {
			logRestore.Errorln("Could not decode FileAcknowledge struct, err: ", err)
			return err
		}
		logRestore.Debugln("Received file acknowledge, file size: ", fileSize.Size)
	}
	return nil
}

func (r *RestoreSession) uploadFile(localFilePath string, size int64) error {
	logRestore.Debugln("Starting sending file: ", localFilePath)
	file, err := r.MainSession.Storage.OpenFile(localFilePath)
	if err != nil {
		logRestore.Errorln("Cannot create localfile to write, err: ", err)
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(r.MainSession.Conn)
	reader := bufio.NewReader(file)
	var readFromFile int64
	var wroteToConnection int64
	buffer := make([]byte, r.MainSession.Transfer.Buffer)
	if size < int64(r.MainSession.Transfer.Buffer) {
		logRestore.Debugln("Shrinking buffer to filesize: ", size)
		buffer = make([]byte, size)
	}
	for {
		read, err := reader.Read(buffer)
		if err != nil {
			logRestore.Errorln("SRV: Could not read from file reader, error: ", err)
			break
		}
		readFromFile += int64(read)
		wrote, err := writer.Write(buffer[:read])
		if err != nil {
			logRestore.Errorln("SRV: Could not write to connection buffer, error: ", err)
			break
		}
		wroteToConnection += int64(wrote)
		if wroteToConnection == size {
			logRestore.Errorln("Wrote all data to connection")
			break
		}
	}
	writer.Flush()
	return nil
}
