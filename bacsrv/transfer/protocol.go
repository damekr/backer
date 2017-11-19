package transfer

import (
	"bufio"
	"encoding/gob"
	"log"
	"net"
	"os"
	"path"

	"github.com/d8x/bftp/storage"
)

type Session struct {
	ConnParams *ConnParameters
	Conn       net.Conn
	Id         uint64
	Transfer   *Transfer
	Storage    *storage.Local
}

func NewSession(id uint64, params *ConnParameters, conn net.Conn) *Session {
	session := &Session{}
	session.Id = id
	session.ConnParams = params
	session.Conn = conn
	return session
}

func (s *Session) negotiate(protoVersion string) error {
	log.Print("Starting negotiate with server")
	neg := new(Negotiate)
	neg.ProtoVersion = protoVersion
	enc := gob.NewEncoder(s.Conn)
	err := enc.Encode(&neg)
	if err != nil {
		log.Println("Could not encode negotiation struct")
		return err
	}
	negs := new(Negotiate)
	dec := gob.NewDecoder(s.Conn)
	err = dec.Decode(&negs)
	if err != nil {
		log.Println("Could not decode negotiation struct, error: ", err)
		return err
	}

	if negs.ProtoVersion != protoVersion {
		log.Println("Server sent protocol version: ", negs.ProtoVersion)
		return ProtocolVersionMismatch
	}
	log.Println("Got protocol revision: ", negs.ProtoVersion)
	return nil
}

func (s *Session) authenticate(password string) error {
	log.Print("Starting authentication")
	auth := new(Authenticate)
	cipherText, err := Encrypt([]byte(password), []byte(KEY))
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
	authRespond := new(Authenticate)
	decRespond := gob.NewDecoder(s.Conn)
	err = decRespond.Decode(&authRespond)
	if string(authRespond.CiperText) != "passed" {
		log.Println("CLNT: Authentication failed!")
		return AuthenticationFailed
	}
	log.Println("CLNT: Authentication passed!")
	return nil
}

func (s *Session) setupTransfer(transfer *Transfer) error {
	log.Println("Opening transfer with server, type: ", transfer.TransferType)
	tr := gob.NewEncoder(s.Conn)
	err := tr.Encode(&transfer)
	if err != nil {
		log.Println("Could not encode transfer struct, error: ", err)
		return err
	}
	log.Println("Waiting for response...")
	tranType := new(Transfer)
	trInc := gob.NewDecoder(s.Conn)
	err = trInc.Decode(&tranType)
	if err != nil {
		log.Println("Could not decode response of transfer type")
		return err
	}
	//Server response an empty string when does not support requesting operation
	if tranType.TransferType == "" {
		log.Println("Server does not support such operation")
		return ServerDoesNotSupportSuchOperation
	}
	log.Println("Server accepts connection type")
	s.Transfer = tranType
	return nil
}

func (s *Session) GetFile(fileRemotePath, fileLocalPath string) error {
	log.Printf("Downloading file: %s to: %s", fileRemotePath, fileLocalPath)
	transfer := new(Transfer)
	// Transfer type to get file is TGET
	transfer.TransferType = TGET
	transfer.Buffer = BUFFERSIZE
	err := s.setupTransfer(transfer)
	if err != nil {
		log.Fatal("Error when sending transfer type, error: ", err)
	}

	// Sending remote file path using FileTransfer struct
	log.Println("Sending request for file located in remote path: ", fileRemotePath)
	fileTransfer := new(FileTransfer)
	fileTransfer.FullPath = fileRemotePath
	ftran := gob.NewEncoder(s.Conn)
	err = ftran.Encode(&fileTransfer)
	if err != nil {
		log.Println("Could not encode FileTransfer struct")
		return err
	}

	// Receiving file size with the same struct as above
	log.Println("Waiting for accepting transfering file")
	fileRecTransfer := new(FileTransfer)
	frect := gob.NewDecoder(s.Conn)
	err = frect.Decode(&fileRecTransfer)
	if err != nil {
		log.Println("Could not decode FileTransfer struct from server, err: ", err)
		return err
	}

	if fileRecTransfer.FullPath == "" {
		//	FilePath is empty so does not exist or an error on server side
		return FileDoesNotExist
	}
	log.Println("Creating local file: ", fileLocalPath)
	localStorage := storage.NewLocalStorage(path.Dir(fileLocalPath))
	file, err := localStorage.CreateFile(path.Base(fileLocalPath))
	if err != nil {
		log.Println("Cannot create localfile to write")
		return err
	}
	defer file.Close()

	// Downloading file part
	log.Println("Server has such file starting downloading file to local path: ", fileLocalPath)
	err = s.downloadFile(file, fileRecTransfer.FileSize)
	if err != nil {
		log.Println("Could not download file, error: ", err)
		return err
	}
	return nil
}

func (s *Session) PutFile(fileLocalPath, fileRemotePath string) error {
	transfer := new(Transfer)
	// Transfer type to put is TPUT
	transfer.TransferType = TPUT
	err := s.setupTransfer(transfer)
	if err != nil {
		log.Println("Error when sending transfer type, error: ", err)
	}
	fileTransfer := new(FileTransfer)
	if storage.CheckIfFileExists(fileLocalPath) {
		fileTransfer.FileSize = storage.GetFileSize(fileLocalPath)
		log.Println("Size of sending file: ", fileTransfer.FileSize)
	} else {
		return FileDoesNotExist
	}

	// Sending file info
	fileTransfer.FullPath = fileRemotePath
	ftran := gob.NewEncoder(s.Conn)
	err = ftran.Encode(&fileTransfer)
	if err != nil {
		log.Fatal("Could not encode FileTransfer struct")
	}

	// Receiving the same as acknowledge
	fileRecTransfer := new(FileTransfer)
	frect := gob.NewDecoder(s.Conn)
	err = frect.Decode(&fileRecTransfer)
	if err != nil {
		log.Println("Could not decode FileTransfer struct from server, error: ", err)
	}

	// Creating local storage for reading from
	localStorage := storage.NewLocalStorage(path.Dir(fileLocalPath))
	file, err := localStorage.OpenFile(path.Base(fileLocalPath))
	if err != nil {
		log.Println("Cannot create localfile to write")
		return err
	}
	defer file.Close()

	// Uploading file using current session
	log.Println("Starting sending file: ", fileLocalPath)
	err = s.uploadFile(file, fileTransfer.FileSize)
	if err != nil {
		log.Println("Could not upload file, error: ", err)
	}
	return nil
}

func (s *Session) downloadFile(file *os.File, size int64) error {
	log.Println("Starting downloading file, size: ", size)
	writer := bufio.NewWriter(file)
	reader := bufio.NewReader(s.Conn)
	var readFromConnection int64
	var wroteToFile int64
	buffer := make([]byte, s.Transfer.Buffer)
	if size < int64(s.Transfer.Buffer) {
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

func (s *Session) uploadFile(file *os.File, size int64) error {
	log.Println("Starting uploading file, size: ", size)
	reader := bufio.NewReader(file)
	writer := bufio.NewWriter(s.Conn)
	var readFromFile int64
	var wroteToConnection int64
	buffer := make([]byte, s.Transfer.Buffer)
	if size < int64(s.Transfer.Buffer) {
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
