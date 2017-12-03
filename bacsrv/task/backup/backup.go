package backup

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/damekr/backer/bacsrv/network"
	"github.com/damekr/backer/bacsrv/storage"
	"github.com/damekr/backer/common/proto"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "task:backup"})

type Backup struct {
	ClientIP       string   `json:"clientIP"`
	RequestedPaths []string `json:"requestedPaths"`
	Progress       int      `json:"-"`
	Status         bool     `json:"status"`
	BucketLocation string   `json:"bucketLocation"`
	ValidPaths     []string `json:"validPaths"`
}

// TODO - maybe make tasks like: backupDefinition, BackupTask, PreBackupTask, PostBackupTask

func CreateBackup(clientIP string, paths []string) *Backup {
	return &Backup{
		ClientIP:       clientIP,
		RequestedPaths: paths,
	}
}

func (b *Backup) Run() {
	log.Println("Running backup of client client: ", b.ClientIP)
	conn, err := network.EstablishGRPCConnection(b.ClientIP)
	if err != nil {
		log.Errorf("Cannot connect to address %s", b.ClientIP)
		return
	}
	defer conn.Close()
	c := proto.NewBacsrvClient(conn)
	response, err := c.Backup(context.Background(), &proto.BackupRequest{Paths: b.RequestedPaths})
	if err != nil {
		log.Warningf("Could not get paths of client: %v", err)
		b.Status = false
		return
	}
	b.ValidPaths = response.BaclntBackupResponse.Validpaths
	b.Status = true
	if err := b.createMetadata(); err != nil {
		log.Errorln("Could not write metadata of backup, err: ", err.Error())
	}
	log.Println("Response: ", response)
}

func (b *Backup) Stop() {
	fmt.Println("Stopping")
}

func (b *Backup) Setup(paths []string) {
	b.RequestedPaths = paths
}

func (b *Backup) createMetadata() error {
	log.Infoln("Creating backup metadata...")
	jsonData, err := json.Marshal(b)
	if err != nil {
		return err
	}
	err = storage.WriteBackupMetadata(jsonData)
	if err != nil {
		return err
	}
	return nil
}
