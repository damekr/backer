package transfer

import (
	"io"
	"net"
	"os"
	// "strconv"
	// "strings"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/repository"
	"github.com/damekr/backer/common/dataproto"
	"path"
	"strings"
)

// BUFFERSIZE determines how big is piece of data that will be send in one frame
const BUFFERSIZE = 1024

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
		log.Debugf("Using saveset with name: %s for this fullbackup", savesetFullPath)
		fileInfo, err := dataproto.UnmarshalFileInfoHeader(conn)
		if err != nil {
			log.Error("Failed during read fileinfo header")
			conn.Close()
		}
		log.Debugf("Received header %#v about file being transferd", fileInfo)
		log.Info("Starting receiving file: ", fileInfo.Name)
		err = receiveFile(fileInfo.Size, savesetFullPath, fileInfo.Name, conn)
		if err != nil {
			log.Errorf("File %s has not been properly received", fileInfo.Name)
		}
		log.Debug("Closing connection with client: ", conn.RemoteAddr())
	}

	defer conn.Close()
}

func receiveFile(fileSize int64, savesetFullPath, fileName string, connection net.Conn) error {
	log.Debugf("Creating file: %s in saveset: %s", fileName, savesetFullPath)
	fileUnderSavesetPath := path.Join(savesetFullPath, fileName)
	log.Debug("File under saveset: ", fileUnderSavesetPath)
	newFile, err := os.Create(fileUnderSavesetPath)
	if err != nil {
		log.Errorf("Cannot create file: %s in repository, error: %s", fileName, err.Error())
	}
	defer newFile.Close()
	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	log.Debugf("File %s has been correctly received", fileName)
	return nil
}
