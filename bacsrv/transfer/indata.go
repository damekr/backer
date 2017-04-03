package transfer

import (
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
)

// BUFFERSIZE determines how big is piece of data that will be send in one frame
const BUFFERSIZE = 1024

// InitTransferServer starts main part of transfering data, consider later to running this on demand
func InitTransferServer(srvConfig *config.ServerConfig) {
	listener, err := net.Listen("tcp", "localhost:"+srvConfig.DataPort)
	log.Info("Starting transfer server on addr: ", listener.Addr())
	if err != nil {
		log.Error("Cannot start transfer server, error: ", err.Error())
	}
	for {
		connection, err := listener.Accept()
		if err != nil {
			log.Error("Cannot accept conection, error: ", err.Error())
		} else {
			log.Debug("A new transfer connection estabilished, from: ", connection.RemoteAddr())
			fileSize := GetFileSize(connection)
			fileName := GetFileName(connection)
			// Part of receiving file
			ReceiveFile(fileSize, fileName, connection)
		}
	}
}

// ReceiveFile is able to read data from buffer and save them in created file.
// It also checks if retrived file is equeal to sent earlier in first chunks of data.
func ReceiveFile(fileSize int64, fileName string, connection net.Conn) {
	newFile, err := os.Create(config.GetMainRepositoryLocation() + fileName)
	if err != nil {
		log.Panic("Cannot create uniq local file in repository")
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
}

// GetFileSize gets file size from given buffer --> remember to send and receive data in proper order
func GetFileSize(connection net.Conn) int64 {
	bufferFileSize := make([]byte, 10)
	_, err := connection.Read(bufferFileSize)
	if err != nil {
		log.Error("Cannot get file size through transfer, error: ", err.Error())
	}
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	log.Debug("Received file with size: ", fileSize)
	return fileSize
}

// GetFileName returns filename from given connection --> remember to send data and read
func GetFileName(connection net.Conn) string {
	bufferFileName := make([]byte, 64)
	_, err := connection.Read(bufferFileName)
	if err != nil {
		log.Error("Cannot read file name, error: ", err.Error())
	}
	fileName := strings.Trim(string(bufferFileName), ":")
	log.Debug("Receiving file name: ", fileName)
	return fileName
}
