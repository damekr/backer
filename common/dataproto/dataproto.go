package dataproto

import (
	"encoding/gob"
	log "github.com/Sirupsen/logrus"
	"io"
	net "net"
	"os"
	"github.com/damekr/backer/baclnt/backup"
	"path"
	"encoding/hex"
	md5 "crypto/md5"
	"errors"
)

const (
	DELIMITER = "\r\n"
	BUFFERSIZE = 1024
)

func init() {
	log.Debug("Initializing transfer protocol")
}


type Transfer struct {
	From	string
	Connection net.Conn
}



type ReturnMessage struct {
	Status bool
	Message	string
}


func New(from string, conn net.Conn) *Transfer {
	return &Transfer{
		From:         from,
		Connection: 		conn,
	}
}


func sendDataTypeHeader(transferType string, conn net.Conn) error {
	log.Debug("Marshaling data for header transfer")
	enc := gob.NewEncoder(conn)
	err := enc.Encode(transferType)
	if err != nil {
		log.Errorf("Cannot encode transfer header with type: %s", transferType)
		return err
	}
	return nil
}

func receiveDataTypeHeader(conn net.Conn) (string, error) {
	log.Debug("Receiving transfer type header")
	var transferType string
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&transferType)
	if err != nil {
		log.Error("Cannot decode transfer header data")
		return "", err
	}
	return transferType, nil
}


func sendDelimiter(conn net.Conn) error {
	log.Debug("Sending delimiter")
	s, err := conn.Write([]byte(DELIMITER))
	if err != nil {
		log.Error("Error while sending delimiter: ", err.Error())
		return err
	}
	log.Debug("Sent delimiter size: ", s)
	return nil
}

func receiveDelimiter(conn net.Conn) error {
	log.Debug("Receiving delimiter")
	buff := make([]byte, len(DELIMITER))
	s, err := conn.Read(buff)
	log.Debug("Read delimiter size: ", s)
	if err != nil {
		if err == io.EOF {
			log.Debug("Client closed connection")
		} else {
			log.Error("Error while reading delimiter: ", err.Error())
			return err
		}
	}
	log.Debug("Received delimiter size: ", s)
	return nil
}


func (t *Transfer) SendTypeHeader(transferType string) error {
	log.Debugf("Sending transfer type: %s to server", transferType)
	err := sendDataTypeHeader(transferType, t.Connection)
	if err != nil {
		log.Error("Encoding transfer header failed")
		return err
	}
	err = sendDelimiter(t.Connection)
	if err != nil {
		log.Error("Error when sending delimiter after data type header")
	}
	return nil
}

func (t *Transfer) ReceiveTypeHeader() (string, error){
	log.Debug("Receiving transfer type header")
	transferType, err := receiveDataTypeHeader(t.Connection)
	if err != nil {
		log.Error("Could not decode transfer type header, closing connection")
		t.Connection.Close()
		return "", err
	}
	log.Debug("Received transfer type: ", transferType)
	err = receiveDelimiter(t.Connection)
	if err != nil{
		return "", err
	}
	return transferType, nil
}

func (t *Transfer) SendFile(fileLocation string) error {
	log.Debug("Sending file ", path.Base(fileLocation))
	err := t.sendFileHeader(fileLocation)
	fileBuffer := make([]byte, BUFFERSIZE)
	file, err := openFile(fileLocation)
	if err != nil {
		log.Errorf("Couldn't read file %s, skipping", fileLocation)
		return nil
	}
	defer file.Close()
	for {
		_, err := file.Read(fileBuffer)
		if err == io.EOF {
			break
		}
		t.Connection.Write(fileBuffer)
	}
	log.Debugf("File: %s has been sent", path.Base(fileLocation))

	log.Debug("Sending delimiter after file send")
	err = sendDelimiter(t.Connection)
	if err != nil {
		log.Error("Error when sending delimiter after file transfer")
	}
	return nil
}



func (t *Transfer) ReceiveFile(saveFullDirectory string) error {
	log.Debug("Starting receiving file, writing it to directory: ", saveFullDirectory)
	fileInfo, err := receiveFileInfoHeader(t.Connection)
	if err != nil {
		log.Error("Error while reading file header, error: ", err.Error())
	}
	receivingFile, err := os.Create(path.Join(saveFullDirectory, fileInfo.Name))
	if err != nil {
		log.Error("Cannot create file to writing in: ", saveFullDirectory)
	}
	defer receivingFile.Close()
	log.Debug("Created file to write: ", receivingFile.Name())
	var receivedBytes int64
	for {
		if (fileInfo.Size - receivedBytes) < BUFFERSIZE {
			if fileInfo.Size == 0 {
				break
			}
			io.CopyN(receivingFile, t.Connection, fileInfo.Size - receivedBytes)
			t.Connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileInfo.Size))
			break
		}
		io.CopyN(receivingFile, t.Connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	log.Debug("Reading delimiter after file transfer")
	err = receiveDelimiter(t.Connection)
	if err != nil {
		log.Error("Error while reading delimiter in file receiving")
	}
	if checkFileChecksum(receivingFile.Name(), fileInfo.Checksum) != nil {
		log.Error("File does not match checksum")
		return errors.New("Checksum does not match")
	}
	log.Debugf("File %s has been correctly received", receivingFile.Name())
	return nil
}


func (t *Transfer) sendFileHeader(fileLocation string) error {
	log.Debug("Reading file info header")
	fileHeader, err := backup.ReadFileHeader(fileLocation)
	if err != nil {
		log.Error("File does not exist")
		return err
	}
	log.Debug("Sending file info header")
	err = sendFileInfoHeader(fileHeader, t.Connection)
	if err != nil {
		log.Error("Encoding and sending file type header failed")
		return err
	}
	log.Debug("Sending delimiter after file header")
	err = sendDelimiter(t.Connection)
	if err != nil {
		log.Error("Could not send delimiter")
		return err
	}
	return nil
}

func (t *Transfer) receiveFileHeader() (backup.FileTransferInfo, error) {
	log.Debug("Receiving file info header")
	fileInfo, err := receiveFileInfoHeader(t.Connection)
	if err == io.EOF {
		log.Debug("No more files to receive, waiting for client to close connection")
		return fileInfo, err
	} else if err != nil {
		log.Error("Error while receiving file info header, error: ", err.Error())
		return fileInfo, err
	}
	log.Debug("Received file info header: ", fileInfo)
	log.Debug("Reading delimiter after file header")
	err = receiveDelimiter(t.Connection)
	if err != nil {
		log.Error("Could not receive delimiter after file header")
		return fileInfo, err
	}
	return fileInfo, nil
}


func receiveFileInfoHeader(conn net.Conn) (backup.FileTransferInfo, error) {
	log.Debug("Reading file header")
	var fileInfo backup.FileTransferInfo
	dec := createDecoder(conn)
	err := dec.Decode(&fileInfo)
	if err != nil {
		log.Error("Could not decode file info header")
		return fileInfo, err
	}
	return fileInfo, err
}

func createEncoder(conn net.Conn) *gob.Encoder {
	log.Debug("Creating encoder")
	return gob.NewEncoder(conn)
}

func createDecoder(conn net.Conn) *gob.Decoder {
	log.Debug("Creating decoder")
	return gob.NewDecoder(conn)
}
func sendFileInfoHeader(fileInfo *backup.FileTransferInfo, conn net.Conn) error {
	log.Debugf("Sending file header:  %#v", fileInfo)
	enc := createEncoder(conn)
	err := enc.Encode(fileInfo)
	if err != nil {
		log.Error("Could not encode file info header")
		return err
	}
	return nil
}

func openFile(fileLocation string) (*os.File, error) {
	// TODO Check if file exists
	file, err := os.Open(fileLocation)
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return file, nil
}

func checkFileChecksum(fileLocation, checksum string) error {
	log.Debugf("Checking received %s file checksum", checksum)
	file, err := os.Open(fileLocation)
	if err != nil {
		return err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String := hex.EncodeToString(hashInBytes)
	if returnMD5String != checksum {
		log.Errorf("Calculation of checksum failed - was: %s is: %s", checksum, returnMD5String)
	} else {
		log.Debugf("Calculation of checksum of file: %s passsed", file.Name())
	}
	return nil
}
