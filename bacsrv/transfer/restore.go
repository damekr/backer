package transfer

import (
	"bufio"
	"encoding/gob"

	"github.com/damekr/backer/bacsrv/storage"
	"github.com/damekr/backer/common"
)

type RestoreSession struct {
	MainSession *MainSession
}

func CreateRestoreSession(mainSession *MainSession) *RestoreSession {
	return &RestoreSession{
		MainSession: mainSession,
	}
}

//TODO Add logic to read file from specific backup(json file)

func (r *RestoreSession) HandleRestoreSession() error {
	log.Println("Handling incomming TGET transfer type")
	fileT := new(common.FileMetadata)
	fileTEmpty := new(common.FileMetadata)

	//Decoding file path to be transfered
	fileTDec := gob.NewDecoder(r.MainSession.Conn)
	err := fileTDec.Decode(&fileT)
	if err != nil {
		log.Print("Could not decode FileMetadata struct, error: ", err)
		fileTEnc := gob.NewEncoder(r.MainSession.Conn)
		if err := fileTEnc.Encode(&fileTEmpty); err != nil {
			log.Println("Could not encode empty FileMetadata struct")
			return err
		}
		return err
	}
	log.Printf("Checking if file %s exists", fileT.FullPath)
	if !storage.CheckIfFileExists(fileT.FullPath) {
		log.Printf("File: %s does not exist", fileT.FullPath)
		fileTEncNotExist := gob.NewEncoder(r.MainSession.Conn)
		if err := fileTEncNotExist.Encode(&fileTEmpty); err != nil {
			log.Println("Could not encode empty FileMetadata struct")
			return err
		}
	}

	//Sending size of file being transfered
	fileTEnc := gob.NewEncoder(r.MainSession.Conn)
	fileT.FileSize = storage.GetFileSize(fileT.FullPath)
	if err := fileTEnc.Encode(&fileT); err != nil {
		log.Println("Could not encode empty FileMetadata struct")
		return err
	}
	log.Println("Handling transfer with sending file to client, file: ", fileT.FullPath)

	//Sending file
	err = r.uploadFile(fileT.FullPath, fileT.FileSize)
	if err != nil {
		log.Println("Could not send file, err: ", err.Error())
	}

	//Receiving file acknowledge
	fileSize := new(common.FileAcknowledge)
	fileSizeEncoder := gob.NewDecoder(r.MainSession.Conn)
	if err := fileSizeEncoder.Decode(&fileSize); err != nil {
		log.Println("Could not decode FileAcknowledge struct")
		return err
	}
	log.Println("Received file acknowledge, file size: ", fileSize.Size)

	return nil
}

func (r *RestoreSession) uploadFile(localFilePath string, size int64) error {
	log.Println("Starting sending file: ", localFilePath)
	file, err := r.MainSession.Storage.OpenFile(localFilePath)
	if err != nil {
		log.Println("Cannot create localfile to write")
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(r.MainSession.Conn)
	reader := bufio.NewReader(file)
	var readFromFile int64
	var wroteToConnection int64
	buffer := make([]byte, r.MainSession.Transfer.Buffer)
	if size < int64(r.MainSession.Transfer.Buffer) {
		log.Println("Shrinking buffer to filesize: ", size)
		buffer = make([]byte, size)
	}
	for {
		read, err := reader.Read(buffer)
		if err != nil {
			log.Println("SRV: Could not read from file reader, error: ", err)
			break
		}
		readFromFile += int64(read)
		wrote, err := writer.Write(buffer[:read])
		if err != nil {
			log.Println("SRV: Could not write to connection buffer, error: ", err)
			break
		}
		wroteToConnection += int64(wrote)
		if wroteToConnection == size {
			log.Println("Wrote all data to connection")
			break
		}
	}
	writer.Flush()
	return nil
}
