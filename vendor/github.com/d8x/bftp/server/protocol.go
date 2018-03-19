package server

import (
	"bufio"
	"encoding/gob"
	"github.com/d8x/bftp/common"
	"github.com/d8x/bftp/storage"
	"io"
	"log"
	"net"
	"path/filepath"
)

type Session struct {
	ConnParams *common.ConnParameters
	Conn       net.Conn
	Id         uint64
	Transfer   *common.Transfer
	Storage storage.Storage
}

func NewSession(id uint64, params *common.ConnParameters, conn net.Conn, storage storage.Storage) *Session {
	session := &Session{}
	session.Id = id
	session.ConnParams = params
	session.Conn = conn
	session.Storage = storage
	return session
}

func (s *Session) Negotiate(protoVersion string) error {
	log.Println("Starting negotiate with client")
	neg := new(common.Negotiate)
	dec := gob.NewDecoder(s.Conn)
	err := dec.Decode(&neg)
	if err != nil {
		log.Println("Could not decode negotatiation struct")
		return err
	}
	log.Println("Got protocol revision: ", neg.ProtoVersion)
	if neg.ProtoVersion != protoVersion {
		return common.ProtocolVersionMismatch
	}
	negc := new(common.Negotiate)
	negc.ProtoVersion = protoVersion
	enc := gob.NewEncoder(s.Conn)
	err = enc.Encode(&negc)
	if err != nil {
		log.Println("Could not encode negotiation struct")
		return err
	}
	return nil
}

func (s *Session) Authenticate(password string) error {
	log.Println("Starting authentication")
	auth := new(common.Authenticate)
	dec := gob.NewDecoder(s.Conn)
	err := dec.Decode(&auth)
	if err != nil {
		log.Println("Could not decode authentication struct")
		return err
	}
	authRespond := new(common.Authenticate)

	receivedPass, err := common.Decrypt(auth.CiperText, []byte(common.KEY))
	if password != string(receivedPass) {
		log.Println("SRV: Authentication failed")
		authRespond.CiperText = []byte("failed")
		dec := gob.NewEncoder(s.Conn)
		err = dec.Encode(&authRespond)
		return common.AuthenticationFailed
	}
	authRespond.CiperText = []byte("passed")
	decRespond := gob.NewEncoder(s.Conn)
	err = decRespond.Encode(&authRespond)
	log.Println("SRV: Authentication passed!")
	return nil
}

func (s *Session) TransferHandler(saveset string) error {
	log.Println("Dispatching incoming connection")
	transfer := new(common.Transfer)
	tr := gob.NewDecoder(s.Conn)
	err := tr.Decode(&transfer)
	if err != nil {
		if err == io.EOF {
			log.Println("Client closed connection")
			return err
		}
		log.Println("Cannot decode incoming transfer type, responding with empty struct. Error: ", err)
	}
	log.Println("Got incoming transfer type connection: ", transfer.TransferType)
	s.Transfer = transfer
	switch transfer.TransferType {
	case common.TGET:
		transfer.Buffer = common.BUFFERSIZE
		trc := gob.NewEncoder(s.Conn)
		err = trc.Encode(&transfer)
		if err != nil {
			log.Println("Could not send transfer type response to client")
		}
		return s.handleTGETOperation()

	case common.TPUT:
		transfer.Buffer = common.BUFFERSIZE
		trp := gob.NewEncoder(s.Conn)
		err = trp.Encode(&transfer)
		if err != nil {
			log.Println("Could not send transfer type response to client")
		}
		return s.handleTPUTOperation(saveset)

	default:
		trEmpty := new(common.Transfer)
		trE := gob.NewEncoder(s.Conn)
		err = trE.Encode(&trEmpty)
		if err != nil {
			log.Println("Could not encode empty struct in non handled transfer type")
			return err
		}
	}

	return nil
}

func (s *Session) downloadFile(name , localFilePath string, size int64, saveset string) error {
	log.Println("Starting downloading file to path: ", localFilePath)
	file, err := s.Storage.CreateFile(filepath.Join(saveset, localFilePath), name)
	if err != nil {
		log.Println("Cannot create localfile to write")
		return err
		//	TODO Respond with failed transfer, error on server side
	}
	defer file.Close()
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
	return nil
}

func (s *Session) uploadFile(localFilePath string, size int64) error {
	log.Println("Starting sending file: ", localFilePath)
	// Creating storage in current path, just for testing
	file, err := s.Storage.OpenFile(localFilePath)
	if err != nil {
		log.Println("Cannot create localfile to write")
		return err
	}
	defer file.Close()
	writer := bufio.NewWriter(s.Conn)
	reader := bufio.NewReader(file)
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

func (s *Session) handleTGETOperation() error {
	log.Println("Handling incomming TGET transfer type")
	fileT := new(common.FileTransfer)
	fileTEmpty := new(common.FileTransfer)

	//Decoding file path to be transfered
	fileTDec := gob.NewDecoder(s.Conn)
	err := fileTDec.Decode(&fileT)
	if err != nil {
		log.Print("Could not decode FileTransfer struct, error: ", err)
		fileTEnc := gob.NewEncoder(s.Conn)
		if err := fileTEnc.Encode(&fileTEmpty); err != nil {
			log.Println("Could not encode empty FileTransfer struct")
			return err
		}
		return err
	}
	log.Printf("Checking if file %s exists", fileT.FullPath)
	if !storage.CheckIfFileExists(fileT.FullPath) {
		log.Printf("File: %s does not exist", fileT.FullPath)
		fileTEncNotExist := gob.NewEncoder(s.Conn)
		if err := fileTEncNotExist.Encode(&fileTEmpty); err != nil {
			log.Println("Could not encode empty FileTransfer struct")
			return err
		}
	}
	//TODO Refactor me

	//Sending size of file being transfered
	fileTEnc := gob.NewEncoder(s.Conn)
	fileT.FileSize = storage.GetFileSize(fileT.FullPath)
	if err := fileTEnc.Encode(&fileT); err != nil {
		log.Println("Could not encode empty FileTransfer struct")
		return err
	}
	log.Println("Handling transfer with sending file to client, file: ", fileT.FullPath)

	//Sending file
	err = s.uploadFile(fileT.FullPath, fileT.FileSize)
	if err != nil {
		log.Println("Could not send file, err: ", err.Error())
	}

	//Receiving file acknowledge
	fileSize := new(common.FileAcknowledge)
	fileSizeEncoder := gob.NewDecoder(s.Conn)
	if err := fileSizeEncoder.Decode(&fileSize); err != nil {
		log.Println("Could not decode FileAcknowledge struct")
		return err
	}
	log.Println("Received file acknowledge, file size: ", fileSize.Size)

	return nil
}

func (s *Session) handleTPUTOperation(saveset string) error {
	log.Println("Handling incomming TPUT transfer type")
	fileTransferPutInfo := new(common.FileTransfer)
	fileTEmpty := new(common.FileTransfer)
	fileTDec := gob.NewDecoder(s.Conn)
	err := fileTDec.Decode(&fileTransferPutInfo)
	if err != nil {
		log.Print("Coult not decode FileTransfer struct, error: ", err)
		fileTEnc := gob.NewEncoder(s.Conn)
		if err := fileTEnc.Encode(&fileTEmpty); err != nil {
			log.Println("Could not encode empty FileTransfer struct")
			return err
		}
	}
	// Sending acknowledge
	// TODO Make checks like: disk space
	fileAEnc := gob.NewEncoder(s.Conn)
	if err := fileAEnc.Encode(&fileTransferPutInfo); err != nil {
		log.Println("Could not send acknowledge")
		return err
	}



	// Downloading file
	err = s.downloadFile(fileTransferPutInfo.Name, fileTransferPutInfo.FullPath, fileTransferPutInfo.FileSize, saveset)
	if err != nil {
		log.Println("Cannot upload file, err: ", err.Error())
		return err
	}
	log.Println("Received file, sending acknowledge")

	//Sending file acknowledge
	fileSize := storage.GetFileSize(fileTransferPutInfo.FullPath)

	fileSizeAckn:= new(common.FileAcknowledge)
	fileSizeEncoder := gob.NewEncoder(s.Conn)
	fileSizeAckn.Size = fileSize
	if err := fileSizeEncoder.Encode(&fileSizeAckn); err != nil {
		log.Println("Could not encode FileAcknowledge struct")
		return err
	}
	log.Println("Sent file size acknowledge")
	return nil
}
