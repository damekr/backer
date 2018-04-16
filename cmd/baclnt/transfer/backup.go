package transfer

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/damekr/backer/cmd/baclnt/fs"
	"github.com/damekr/backer/pkg/bftp"
)

type BackupSession struct {
	MainSession *MainSession
	FileSystem  fs.FileSystem
}

func CreateBackupSession(mainSession *MainSession, fileSystem fs.FileSystem) *BackupSession {
	return &BackupSession{
		MainSession: mainSession,
		FileSystem:  fileSystem,
	}
}

func (b *BackupSession) sendDirsMetadata(dirPaths []string) error {
	log.Debugln("Sending dirs metadata")
	dirsMetadata := b.readDirsMetadata(dirPaths)
	dirsMetadataEncoder := gob.NewEncoder(b.MainSession.Conn)
	err := dirsMetadataEncoder.Encode(&dirsMetadata)
	if err != nil {
		log.Error("Could not encode DirMetadata struct, err: ", err)
		return err
	}
	return nil
}

func (b *BackupSession) readDirsMetadata(dirPaths []string) []*bftp.DirMetadata {
	log.Debugln("Reading paths structure")
	var dirsMetadata []*bftp.DirMetadata
	for _, v := range dirPaths {
		dirMetadata, err := b.FileSystem.ReadDirMetadata(v)
		if err != nil {
			log.Errorln("Cannot read dir metadata, err: ", err)
		} else {
			dirsMetadata = append(dirsMetadata, dirMetadata)
		}
	}
	return dirsMetadata
}

func (b *BackupSession) receiveEmptyAckMessage() error {
	log.Debugln("Receiving empty Ack Message")
	emptyAck := new(bftp.EmtpyAck)
	dirsMetadataDecoder := gob.NewDecoder(b.MainSession.Conn)
	err := dirsMetadataDecoder.Decode(&emptyAck)
	if err != nil {
		log.Errorln("Cannot receive empty ack message, err: ", err)
		return err
	}
	return nil
}

func (b *BackupSession) putFile(fileLocalPath, fileRemotePath string) error {

	fileMetadata, err := b.FileSystem.ReadFileMetadata(fileRemotePath)

	if b.FileSystem.CheckIfFileExists(fileLocalPath) {
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

	fileReader, err := b.FileSystem.ReadFile(fileLocalPath)
	if err != nil {
		log.Println("Cannot open file, err: ", err.Error())
	}
	defer fileReader.Close()

	// Uploading file using current backup session
	log.Println("Starting sending file: ", fileLocalPath)
	err = b.uploadFile(fileReader, fileMetadata)
	if err != nil {
		log.Debugln("Could not upload file, error: ", err)
	}
	log.Debugln("Uploaded file, waiting for acknowledge")

	// Receiving file acknowledge
	fileTransferAckn, err := b.receiveFileTransferAcknowledge()
	if err != nil {
		log.Errorln("Could not receive fileTransferAcknowledge")
	}
	log.Debugln("Received file acknowledge, file size: ", fileTransferAckn.Size)

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

func (b *BackupSession) uploadFile(fileReader io.Reader, metadata *bftp.FileMetadata) error {
	log.Println("Starting uploading file, size: ", metadata.FileSize)
	writer := bufio.NewWriter(b.MainSession.Conn)
	var readFromFile int64
	var wroteToConnection int64
	buffer := make([]byte, b.MainSession.Transfer.Buffer)
	if metadata.FileSize < int64(b.MainSession.Transfer.Buffer) {
		log.Println("Shrinking buffer to file size: ", metadata.FileSize)
		buffer = make([]byte, metadata.FileSize)
	}
	for {
		read, err := fileReader.Read(buffer)
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
		if wroteToConnection == metadata.FileSize {
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
