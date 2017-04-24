package transfer

import (
	"io"
	"net"
	"os"
	// "strconv"
	"errors"

	"crypto/md5"
	"encoding/hex"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/repository"
	"github.com/damekr/backer/common/dataproto"
	"path"
	"strings"
)

// BUFFERSIZE determines how big is piece of data that will be send in one frame
const BUFFERSIZE = 1024

// StartTransferServer setup listener on specified port
func StartTransferServer(srvConfig *config.ServerConfig) {
	lis, err := net.Listen("tcp", net.JoinHostPort(srvConfig.DataTransferInterface, srvConfig.DataPort))
	if err != nil {
		log.Errorf("Cannot start data transfer server on interface: %s and port: %s. Error: %s", srvConfig.DataTransferInterface, srvConfig.DataPort, err.Error())
	}
	log.Info("Data server started successfully on interface: ", lis.Addr())
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Errorf("Couldn't accept connection from host: %s, error: %s", conn.RemoteAddr(), err.Error())
		}
		// Starting gourutine to handle single connection
		go dataTransferHandler(conn)
	}
}

func dataTransferHandler(conn net.Conn) {
	log.Debug("Handling data transfer from: ", conn.RemoteAddr())
	header, err := dataproto.UnmarshalTransferHeader(conn)
	if err != nil {
		log.Error("Couldn't read data transfer header")
	}
	log.Debug("Transfer from: ", header.From)
	transferType := strings.ToLower(header.TType)
	log.Debug("Transfer type: ", transferType)
	switch transferType {
	case "fullbackup":
		log.Info("Starting full backup data transfer")
		log.Info("Creating new client bucket for data")
		savesetFullPath, err := repository.CreateClientSaveset(header.From)
		if err != nil {
			log.Error("Failed during creation client saveset")
			conn.Close()
		}
		receiveFiles(conn, savesetFullPath)
		log.Debug("Closing connection with client: ", conn.RemoteAddr())
	}

	// defer conn.Close()
}

func receiveFiles(conn net.Conn, savesetFullPath string) {
	log.Debugf("Using saveset with name: %s for this fullbackup", savesetFullPath)
	for {
		fileInfo, err := dataproto.UnmarshalFileInfoHeader(conn)
		if err == io.EOF {
			log.Debug("Received last file, stopping transfer")
			break
			// TODO - handle random errors, repeate file transfer?
		}
		log.Debugf("Received header %#v about file being transferd", fileInfo)
		log.Info("Starting receiving file: ", fileInfo.Name)
		err = receiveFile(fileInfo.Size, savesetFullPath, fileInfo.Name, fileInfo.Location, fileInfo.Checksum, conn)
		if err != nil {
			log.Error("Received error: ", err.Error())
			log.Errorf("File %s has not been properly received", fileInfo.Name)
		}

	}
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
	}
	log.Debugf("Calculation of checksum of file: %s passsed", file.Name())
	return nil
}

func receiveFile(fileSize int64, savesetFullPath, fileName, fileFullLocation, checksum string, connection net.Conn) error {
	log.Debugf("Creating file: %s in saveset: %s", fileName, savesetFullPath)
	fileDir, _ := path.Split(fileFullLocation)
	log.Debug("Creating proper path under saveset: ", fileDir)
	err := os.MkdirAll(path.Join(savesetFullPath, fileDir), 0700)
	if err != nil {
		log.Error("Couldn't create proper file path under saveset")
		return err
	}
	fileUnderSavesetPath := path.Join(savesetFullPath, fileFullLocation)
	log.Debug("File under saveset: ", fileUnderSavesetPath)
	newFile, err := os.Create(fileUnderSavesetPath)
	if err != nil {
		log.Errorf("Cannot create file: %s in repository, error: %s", fileName, err.Error())
	}
	defer newFile.Close()
	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			if fileSize == 0 {
				// Fast fix for empty files
				break
			}
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	if checkFileChecksum(fileUnderSavesetPath, checksum) != nil {
		log.Error("File does not match checksum")
		return errors.New("Checksum does not match")
	}
	log.Debugf("File %s has been correctly received", fileName)
	return nil
}
