package transfer

import (
	"io"
	"net"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/backup"
	"github.com/damekr/backer/common/dataproto"
	"path"
	"github.com/damekr/backer/baclnt/config"
)

// BUFFERSIZE specifies how big is a chunk of data being sent
const (
	BUFFERSIZE = 1024
	SRVDATAPORT = "8000"
)

type BackupConfig struct {
	Paths       []string
	Exclude     []string
	ArchiveName string
	ArchiveSize string
}

type RestoreConfig struct {
	Override    bool
	ArchiveName string
}

func initDataConnectionWithServer(srvAddr, dataPort string) (net.Conn, error) {
	log.Debug("Starting connection with server: ", srvAddr)
	conn, err := net.Dial("tcp", net.JoinHostPort(srvAddr, dataPort))
	if err != nil {
		log.Error("Couldn't establish connection with: ", srvAddr)
		return nil, err
	}
	log.Debugf("Established connection with server: %s on port: %s", srvAddr, dataPort)
	return conn, nil
}

func initFullBackupDataConnection(srvAddr, dataPort string) (net.Conn, error){
	conn, err := initDataConnectionWithServer(srvAddr, dataPort)
	if err != nil {
		log.Error("Cannot init full backup connection")
		return nil, err
	}
	clientName := config.GetExternalName()
	log.Debugf("Using client name: %s to communicate with server", clientName)
	log.Debug("Sending transfer type header")
	transferHeader := dataproto.CreateFullBackupTypeHeader(clientName)
	err = transferHeader.SendTypeHeader(conn, srvAddr)
	if err != nil {
		log.Error("Sending transfer header failed with error: ", err)
		return nil, err
	}
	return conn, nil

}

func StartFullBackup(paths []string, srvAddr string) error {
	delimiter := make([]byte, 1)

	// Initializing full backup connection
	conn, err := initFullBackupDataConnection(srvAddr, SRVDATAPORT)
	if err != nil {
		log.Error(err)
	}
	defer conn.Close()

	for _, v := range paths {
		log.Debug("Sending file: ", v)
		err = sendFileHeader(conn, v)
		if err != nil {
			log.Errorf("An error occured during sending file: %s header, error: %s", v, err.Error())
		}
		d, err := conn.Read(delimiter)
		if err != nil {
			log.Debug("Correct read delimiter, starting data transfer")
		}
		log.Debug("Read delimiter size: ", d)
		err = sendFile(conn, v)
		if err != nil {
			log.Errorf("An error occured during sending file: %s, error: %s", v, err.Error())
		}
	}
	log.Info("Closing file transfer connection")
	return nil

}

func sendTransferTypeHeader(transfer *dataproto.Transfer, conn net.Conn) error {
	log.Debug("Sending transfer type to server")
	err := dataproto.SendDataTypeHeader(transfer, conn)
	if err != nil {
		log.Error("Encoding transfer header failed")
		return err
	}

	return nil
}

func sendFileHeader(conn net.Conn, fileLocation string) error {
	log.Debug("Sending file info header")
	fileHeader, err := backup.ReadFileHeader(fileLocation)
	if err != nil {
		log.Error("File does not exist")
		return err
	}
	err = dataproto.SendFileInfoHeader(fileHeader, conn)
	if err != nil {
		log.Error("Encoding and sending file type header failed")
		return err
	}
	return nil
}

func sendFile(conn net.Conn, fileLocation string) error {
	log.Debug("Sending file ", path.Base(fileLocation))
	sendBuffer := make([]byte, BUFFERSIZE)
	file, err := openFile(fileLocation)
	if err != nil {
		log.Errorf("Couldn't read file %s, skipping", fileLocation)
		return nil
	}
	defer file.Close()
	for {
		_, err := file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		conn.Write(sendBuffer)
	}
	log.Debugf("File: %s has been sent", path.Base(fileLocation))
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
