package transfer

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"net"
	"strings"
)

// InitTransferServerDispatcher starts main part of transfer server and dispatch which type of transfer will be executed, possibilities:
// - backup (client sends to server data)
// - restore (client sends requests for restore and gets data)
func InitTransferServerDispatcher(srvConfig *config.ServerConfig) {
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
			err := dispatchTransferConnection(connection)
			if err != nil {
				log.Error("An error occured during dispatching transfer connection, error: ", err.Error())
			}
		}
	}
}

func dispatchTransferConnection(connection net.Conn) error {
	log.Debug("Dispatching transfer connection from: ", connection.RemoteAddr())
	header, err := readTransferConnectionHeader(connection)
	if err != nil {
		log.Error("Cannot get header, error: ", err.Error())
	}
	log.Debug("Received header: ", header)
	switch header {
	case "backup":
		log.Info("Dispatching backup transfer operation")
		fileSize := GetFileSize(connection)
		fileName := GetFileName(connection)
		// Part of receiving file
		ReceiveFile(fileSize, fileName, connection)
	case "restore":
		log.Info("Dispatching restore transfer operation")
		requestedArchiveName := GetFileName(connection)
		log.Debug("Client requesting archive: ", requestedArchiveName)
		err := SendArchive(connection, requestedArchiveName)
		if err != nil {
			log.Errorf("Cannot send archive: %s to client: %s", requestedArchiveName, connection.RemoteAddr())
		}
	default:
		log.Errorf("The header: %s does not mean anything", header)
		return errors.New("Did not recognize header as an operation")
	}

	return nil
}

func readTransferConnectionHeader(connection net.Conn) (string, error) {
	headerBuff := make([]byte, 20)
	_, err := connection.Read(headerBuff)
	if err != nil {
		log.Error("Cannot read header from host: ", connection.RemoteAddr())
		return "", err
	}
	log.Debug("Got not trimmed header: ", string(headerBuff))
	header := strings.Trim(string(headerBuff), ":")
	return header, nil
}
