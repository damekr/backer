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

type TransferConnection struct {
	Port       int
	Host       string
	Timeout    int
	BUFFERSIZE int
}

type BackupConfig struct {
	Paths       []string
	Exclude     []string
	ArchiveName string
	ArchiveSize string
	TRConn      net.Conn
}

func (c *TransferConnection) InitConnection() net.Conn {
	log.Debugf("Trying to intialize transfer connection with %s, on port: %s", c.Host, c.Port)
	connection, err := net.Dial("tcp", c.Host+":"+strconv.Itoa(c.Port))
	if err != nil {
		log.Fatal("Cannot initialize transfer connection")
	}
	return connection
}

// func (b *BackupConfig) CreateArchive(paths []string) string {
// 	b.Paths = paths
// 	archive := archiver.NewArchive(b.Paths, "Archiwum")
// 	log.Print("Creating archive in: %v with paths: %v", archivePath, paths)
// 	archive.MakeArchive(Config.TempPath, "Archive")
// 	return archivePath
// }

func (b *BackupConfig) SendArchive(archiveLocation string) {
	b.ArchiveName, b.ArchiveSize = readArchiveMetadata(archiveLocation)
	log.Debugf("Read arch metadata, name: %s, size: %s", b.ArchiveName, b.ArchiveSize)
	// Sending archive size to compare that all has been sent
	outSize, err := b.TRConn.Write([]byte(b.ArchiveSize))
	if err != nil {
		log.Println("An error occured: " + err.Error())
	}
	log.Println(outSize, "bytes sent Name")
	// Sending archive name to use on backend side
	sentDataSize, err := b.TRConn.Write([]byte(b.ArchiveName))
	if err != nil {
		log.Println("An error occured: " + err.Error())
	}
	log.Println(sentDataSize, "bytes sent size")
	// TODO I am not sure that this is proper to close in this place
	// connection, make it maybe in seperate method?
	defer b.TRConn.Close()
	// Sending archive
	sendBuffer := make([]byte, BUFFERSIZE)
	file := readArchive(archiveLocation)
	defer file.Close()
	for {
		_, err := file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		b.TRConn.Write(sendBuffer)
	}
	log.Debug("File has been sent, closing connection!")
	return
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

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}
