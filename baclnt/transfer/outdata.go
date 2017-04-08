package transfer

import (
	"io"
	"net"
	"os"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/archiver"
	"github.com/damekr/backer/baclnt/config"
	"github.com/damekr/backer/common/dataproto"
	"path"
)

// BUFFERSIZE specifies how big is a chunk of data being sent
const BUFFERSIZE = 1024

var Config *config.ClientConfig

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

func InitConnectionWithServer(srvAddr, dataPort string) (net.Conn, error) {
	log.Debug("Starting connection with server: ", srvAddr)
	conn, err := net.Dial("tcp", net.JoinHostPort(srvAddr, dataPort))
	if err != nil {
		log.Error("Couldn't establish connection with: ", srvAddr)
		return nil, err
	}
	log.Debugf("Established connection with server: %s on port: %s", srvAddr, dataPort)
	return conn, nil
}

func SendTransferTypeHeader(ttype, from string, conn net.Conn) error {
	log.Debug("Sending transfer type to server")
	transfer := &dataproto.Transfer{
		TType: ttype,
		From:  from,
	}
	err := dataproto.SendDataTypeHeader(transfer, conn)
	if err != nil {
		log.Error("Encoding transfer header failes")
		return err
	}
	_ = sendFileHeader(conn, "/var/tmp/ala")
	_ = sendFile(conn, "/var/tmp/ala")
	return nil
}

func sendFileHeader(conn net.Conn, fileLocation string) error {
	log.Debug("Sending file info header")
	fileHeader, err := archiver.ReadFileHeader(fileLocation)
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
	file := openFile(fileLocation)
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

func openFile(fileLocation string) *os.File {
	// TODO Check if file exists
	file, err := os.Open(fileLocation)
	if err != nil {
		log.Error(err.Error())
	}
	return file
}

// SendArchive sends created archive to the server
func (b *BackupConfig) SendArchive(transferConn net.Conn, archiveLocation string) {
	log.Debug("Sending backup header")
	backupHeader := fillString("backup", 20)
	log.Debug("Restore header: ", backupHeader)
	sentHeaderSize, err := transferConn.Write([]byte(backupHeader))
	if err != nil {
		log.Error("An error occured during sending header message, error: ", err.Error())
	}
	log.Debug("Sent header size: ", sentHeaderSize)
	b.ArchiveName, b.ArchiveSize = readArchiveMetadataInConnectionFormat(archiveLocation)
	log.Debugf("Read arch metadata, name: %s, size: %s", b.ArchiveName, b.ArchiveSize)
	// Sending archive size to compare that all has been sent
	outSize, err := transferConn.Write([]byte(b.ArchiveSize))
	if err != nil {
		log.Println("An error occured: " + err.Error())
	}
	log.Println(outSize, "bytes sent Name")
	// Sending archive name to use on backend side
	sentDataSize, err := transferConn.Write([]byte(b.ArchiveName))
	if err != nil {
		log.Println("An error occured: " + err.Error())
	}
	log.Println(sentDataSize, "bytes sent size")
	// TODO I am not sure that this is proper to close in this place
	// connection, make it maybe in seperate method?
	defer transferConn.Close()
	// Sending archive
	sendBuffer := make([]byte, BUFFERSIZE)
	file := openFile(archiveLocation)
	defer file.Close()
	for {
		_, err := file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		transferConn.Write(sendBuffer)
	}
	log.Debug("File has been sent, closing connection!")
	return
}

func (r *RestoreConfig) ReceiveArchive(transferConn net.Conn, tempRestoreLocation string) error {
	archive, err := os.Create(tempRestoreLocation + "/" + r.ArchiveName)
	if err != nil {
		log.Error("Cannot create archive being restored")
	}
	defer archive.Close()
	var archiveSize int64
	archiveSize = GetFileSize(transferConn)
	log.Debug("Restoring file size: ", archiveSize)
	archName := GetFileName(transferConn)
	log.Debug("Restoring file name: ", archName)
	var receivedBytes int64
	for {
		if (archiveSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(archive, transferConn, (archiveSize - receivedBytes))
			transferConn.Read(make([]byte, (receivedBytes+BUFFERSIZE)-archiveSize))
			break
		}
		io.CopyN(archive, transferConn, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
	}
	log.Debugf("File %s has been correctly received", archName)
	return nil
}

// SendRestoreHeader sends "restore" command for requesting a restore
func SendRestoreHeader(transferConn net.Conn) error {
	restoreHeader := fillString("restore", 20)
	log.Print("Restore header: ", restoreHeader)
	sentHeaderSize, err := transferConn.Write([]byte(restoreHeader))
	if err != nil {
		log.Error("An error occured during sending header message, error: ", err.Error())
	}
	log.Debug("Sent header size: ", sentHeaderSize, " Sending archive Name")
	sentArchNameSize, err := transferConn.Write([]byte("dummyName"))
	if err != nil {
		log.Error("An error occured during sending header message, error: ", err.Error())
	}
	log.Debug("Sent Archive name size: ", sentArchNameSize)
	return nil
}

func readArchiveMetadataInConnectionFormat(archiveLocation string) (string, string) {
	file, err := os.Open(archiveLocation)
	if err != nil {
		log.Error(err.Error())
	}
	fileInfo, err := file.Stat()
	if err != nil {
		log.Error(err.Error())
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	log.Debugf("File Size: %s\nFile Name: %s\n", fileSize, fileName)
	return fileName, fileSize
}

func fillString(returnString string, toLength int) string {
	log.Debug("String length to be filled up: ", len(returnString))
	for {
		lengtString := len(returnString)
		if lengtString < toLength {
			returnString = returnString + ":"
			continue
		}
		break
	}
	return returnString
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
