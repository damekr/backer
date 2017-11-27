package backup

import (
	"context"
	"fmt"

	"github.com/d8x/bftp"
	"github.com/damekr/backer/bacsrv/network"
	"github.com/damekr/backer/common/proto"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "task:backup"})

type Backup struct {
	ClientIP       string
	RequestedPaths []string
	Progress       int
	ValidPaths     []string
	Status         bool
}

// TODO - maybe make tasks like: backupDefinition, BackupTask, PreBackupTask, PostBackupTask

func CreateBackup(clientIP string, paths []string) *Backup {
	return &Backup{
		ClientIP:       clientIP,
		RequestedPaths: paths,
	}
}

func (b *Backup) Run() {
	log.Println("Running fs of client client: ", b.ClientIP)
	conn, err := network.EstablishGRPCConnection(b.ClientIP)
	if err != nil {
		log.Errorf("Cannot connect to address %s", b.ClientIP)
		return
	}
	defer conn.Close()
	c := proto.NewBacsrvClient(conn)
	r, err := c.Backup(context.Background(), &proto.BackupRequest{Paths: b.RequestedPaths})
	if err != nil {
		log.Warningf("Could not get paths of client: %v", err)
		b.ValidPaths = []string{}
		b.Status = false
		return
	}
	b.Status = true
	log.Println("Response: ", r)
	b.ValidPaths = r.BaclntBackupResponse.Validpaths
	b.StartBackup()

}

func (b *Backup) StartBackup() error {
	bftpClient := bftp.CreateBFTPClient()
	session, err := bftpClient.Connect("127.0.0.1", 8000)
	if err != nil {
		log.Errorln("Cannot connect to client: ", b.ClientIP)
	}
	log.Println("Session ID: ", session.Id)
	return nil
}

func (b *Backup) Stop() {
	fmt.Println("Stopping")
}

func (b *Backup) Setup(paths []string) {
	b.RequestedPaths = paths
}
