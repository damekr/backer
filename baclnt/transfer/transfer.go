package transfer

import (
	"io"
	"net"
	"os"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/config"
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

func InitConnection(host string, port string) net.Conn {
	log.Debugf("Trying to intialize transfer connection with %s, on port: %s", host, port)
	connection, err := net.Dial("tcp", host+":"+port)
	if err != nil {
		log.Fatal("Cannot initialize transfer connection")
	}
	return connection
}

func (b *BackupConfig) SendArchive(transferConn net.Conn, archiveLocation string) {
	b.ArchiveName, b.ArchiveSize = readArchiveMetadata(archiveLocation)
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
	file := readArchive(archiveLocation)
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

	return nil
}

// SendRestoreHeader sends "restore" command for requesting a restore
func SendRestoreHeader(transferConn net.Conn) error {
	restoreHeader := fillString("dupa", 20)
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

func readArchive(archiveLocation string) *os.File {
	// TODO Check if file exists and if is an archive
	file, err := os.Open(archiveLocation)
	if err != nil {
		log.Error(err.Error())
		panic(err)
	}
	return file
}

func readArchiveMetadata(archiveLocation string) (string, string) {
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
