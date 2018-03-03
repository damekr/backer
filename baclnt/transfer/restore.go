package transfer

import (
	"bufio"
	"encoding/gob"
	"os"
	"path"

	"github.com/damekr/backer/baclnt/fs"
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

func (r *RestoreSession) GetFile(fileRemotePath, fileLocalPath string) error {
	log.Printf("Downloading file: %r to: %r", fileRemotePath, fileLocalPath)

	// Sending remote file path using FileMetadata struct
	fileMetadata := new(common.FileMetadata)
	fileMetadata.FullPath = fileRemotePath
	err := r.sendRequestForFile(fileMetadata)
	if err != nil {
		log.Errorln("Could not send file metadata for restore, err: ", err.Error())
	}

	// Receiving file size with the same struct as above
	receivedFileMetadata, err := r.receiveFileMetadata()
	if err != nil {
		log.Errorln("Could not receive file metadata, err: ", err.Error())
	}

	if receivedFileMetadata.FullPath == "" {
		//	FilePath is empty so does not exist or an error on server side
		return common.FileDoesNotExist
	}
	log.Println("Creating local file: ", fileLocalPath)
	localStorage := fs.NewFS(path.Dir(fileLocalPath))
	file, err := localStorage.CreateFile(path.Base(fileLocalPath))
	if err != nil {
		log.Println("Cannot create localfile to write, err: ", err)
		return err
	}
	defer file.Close()

	// Downloading file part
	log.Println("Server has such file starting downloading file to local path: ", fileLocalPath)
	err = r.downloadFile(file, receivedFileMetadata.FileSize)
	if err != nil {
		log.Println("Could not download file, error: ", err)
		return err
	}

	//Sending file acknowledge
	fileTransferAcknowledge, err := r.sendTransferAcknowledge(fileLocalPath)
	if err != nil {
		log.Errorln("Could not send transfer acknowledge, err: ", err)
	}
	log.Println("File transfer acknowledge: ", fileTransferAcknowledge)

	return nil
}

func (r *RestoreSession) sendRequestForFile(metadata *common.FileMetadata) error {
	//TODO When backup metadata will be created, then send here also backup id
	log.Println("Sending request for file located in remote path: ", metadata.FullPath)
	ftran := gob.NewEncoder(r.MainSession.Conn)
	err := ftran.Encode(&metadata)
	if err != nil {
		log.Errorln("Could not encode FileMetadata struct")
		return err
	}
	return nil
}

func (r *RestoreSession) receiveFileMetadata() (*common.FileMetadata, error) {
	log.Println("Waiting for accepting transfering file")
	fileRecTransfer := new(common.FileMetadata)
	frect := gob.NewDecoder(r.MainSession.Conn)
	err := frect.Decode(&fileRecTransfer)
	if err != nil {
		log.Println("Could not decode FileMetadata struct from server, err: ", err)
		return nil, err
	}
	return fileRecTransfer, nil
}

//TODO it is the same in backup also, consider put it into main session struct
func (r *RestoreSession) sendTransferAcknowledge(filePath string) (*common.FileAcknowledge, error) {
	fileSize := fs.GetFileSize(filePath)
	fileSizeAckn := new(common.FileAcknowledge)
	fileSizeEncoder := gob.NewEncoder(r.MainSession.Conn)
	fileSizeAckn.Size = fileSize
	log.Debugln("Filetransfer acknowledge struct: ", fileSizeAckn)
	if err := fileSizeEncoder.Encode(&fileSizeAckn); err != nil {
		log.Warningln("Could not encode FileAcknowledge struct, err: ", err)
		return nil, err
	}
	log.Println("Sent file size acknowledge")
	return fileSizeAckn, nil
}

func (r *RestoreSession) downloadFile(file *os.File, size int64) error {
	log.Println("Starting downloading file, size: ", size)
	writer := bufio.NewWriter(file)
	reader := bufio.NewReader(r.MainSession.Conn)
	var readFromConnection int64
	var wroteToFile int64
	buffer := make([]byte, r.MainSession.Transfer.Buffer)
	if size < int64(r.MainSession.Transfer.Buffer) {
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
			log.Println("Could not write to file writter, error: ", err)
			break
		}
		wroteToFile += int64(wrote)
		if wroteToFile == size {
			log.Println("Wrote all data to file")
			break
		}

	}
	writer.Flush()
	return nil
}
