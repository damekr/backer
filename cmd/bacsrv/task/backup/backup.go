package backup

import (
	"context"
	"fmt"

	"github.com/damekr/backer/api/protoclnt"
	"github.com/damekr/backer/cmd/bacsrv/network"
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
	c := protoclnt.NewBaclntClient(conn)
	response, err := c.Backup(context.Background(), &protoclnt.BackupRequest{Paths: b.RequestedPaths})
	if err != nil {
		log.Warningf("Could not get paths of client: %v", err)
		b.Status = false
		return
	}
	b.ValidPaths = response.Validpaths
	b.Status = true
	log.Println("Response: ", response)
}

func (b *Backup) Stop() {
	fmt.Println("Stopping")
}

func (b *Backup) Setup(paths []string) {
	b.RequestedPaths = paths
}
