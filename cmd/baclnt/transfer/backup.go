package transfer

import (
	"bufio"
	"encoding/gob"
	"os"
	"path"

	"github.com/damekr/backer/cmd/baclnt/fs"
	"github.com/damekr/backer/pkg/bftp"
)

type BackupSession struct {
	MainSession *MainSession
}

func CreateBackupSession(mainSession *MainSession) *BackupSession {
	return &BackupSession{
		MainSession: mainSession,
	}
}

func (b *BackupSession) PutFile(fileLocalPath, fileRemotePath string) error {
	// TODO Create FileSystem once per backup and then read files.
	localfs := fs.NewFS(path.Dir(fileLocalPath))

	fileMetadata, err := localfs.ReadFileMetadata(fileRemotePath)

	if localfs.CheckIfFileExists(fileLocalPath) {
		log.Println("Size of sending file: ", fileMetadata.FileSize)
	} else {
		return bftp.FileDoesNotExist
	}

	// Sending file info
	err = b.sendFileMetadata(fileMetadata)
	if err != nil {
		log.Errorln("Could not send file metadata info, err: ", err.Error())
	}

	// Receiving the same as acknowledge
	acknMetadata, err := b.receiveFileMetadataAcknowledge()
	if err != nil {
		log.Errorln("Could not receive file meta data acknowledge, err: ", err.Error())
	}
	log.Debugln("Received file metadata acknowledge: ", acknMetadata)

	file, err := localfs.OpenFile(path.Base(fileLocalPath))
	if err != nil {
		log.Println("Cannot open localfile, err: ", err.Error())
	}
	defer file.Close()

	// Uploading file using current session
	log.Println("Starting sending file: ", fileLocalPath)
	err = b.uploadFile(file, fileMetadata.FileSize)
	if err != nil {
		log.Println("Could not upload file, error: ", err)
	}
	log.Println("Uploaded file, waiting for acknowledge")

	//Receiving file acknowledge
	fileTransferAckn, err := b.receiveFileTransferAcknowledge()
	if err != nil {
		log.Errorln("Could not receive fileTransferAcknowledge")
	}
	log.Println("Received file acknowledge, file size: ", fileTransferAckn.Size)

	return nil
}

func (b *BackupSession) sendFileMetadata(metadata *bftp.FileMetadata) error {
	ftran := gob.NewEncoder(b.MainSession.Conn)
	err := ftran.Encode(&metadata)
	if err != nil {
		log.Error("Could not encode FileMetadata struct")
		return err
	}
	return nil
}

func (b *BackupSession) receiveFileMetadataAcknowledge() (*bftp.FileMetadata, error) {
	fileRecTransfer := new(bftp.FileMetadata)
	frect := gob.NewDecoder(b.MainSession.Conn)
	err := frect.Decode(&fileRecTransfer)
	if err != nil {
		log.Errorln("Could not decode FileMetadata struct from server, error: ", err)
		return nil, err
	}
	return fileRecTransfer, nil
}

func (b *BackupSession) uploadFile(file *os.File, size int64) error {
	log.Println("Starting uploading file, size: ", size)
	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(b.MainSession.Conn)
	var readFromFile int64
	var wroteToConnection int64
	buffer := make([]byte, b.MainSession.Transfer.Buffer)
	if size < int64(b.MainSession.Transfer.Buffer) {
		log.Println("Shrinking buffer to filesize: ", size)
		buffer = make([]byte, size)
	}
	for {
		read, err := reader.Read(buffer)
		if err != nil {
			log.Print("CLNT - Upload: Could not read from file reader, error: ", err)
			return err
		}
		readFromFile += int64(read)
		wrote, err := writer.Write(buffer[:read])
		if err != nil {
			log.Println("CLNT - Upload: Could not write to connection writter, error: ", err)
			break
		}
		wroteToConnection += int64(wrote)
		if wroteToConnection == size {
			log.Print("Wrote all data to connection")
			break
		}
	}
	writer.Flush()
	return nil
}

func (b *BackupSession) receiveFileTransferAcknowledge() (*bftp.FileAcknowledge, error) {
	fileTransferAcknowledge := new(bftp.FileAcknowledge)
	fileSizeEncoder := gob.NewDecoder(b.MainSession.Conn)
	if err := fileSizeEncoder.Decode(&fileTransferAcknowledge); err != nil {
		log.Println("Could not decode FileAcknowledge struct")
		return nil, err
	}
	return fileTransferAcknowledge, nil
}
