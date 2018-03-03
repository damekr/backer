package restore

import (
	"context"
	"path/filepath"

	"github.com/damekr/backer/bacsrv/db"
	"github.com/damekr/backer/bacsrv/network"
	"github.com/damekr/backer/common/protoclnt"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "task:restore"})

type Restore struct {
	ClientIP               string   `json:"clientIP"`
	BackupID               int      `json:"backupID"`
	OriginalFilesLocations []string `json:"originalFilesLocations"` // Paths to be restored on client
	FilesLocationOnServer  []string `json:"filesLocationOnServer"`  // Paths on server
	Progress               int      `json:"-"`
	Status                 bool     `json:"status"`
	BucketLocation         string   `json:"bucketLocation"`
}

func Create(clientIP string, backupID int) *Restore {
	return &Restore{
		ClientIP: clientIP,
		BackupID: backupID,
	}
}

func (r *Restore) Run() {
	log.Println("Running restore of client: ", r.ClientIP)
	conn, err := network.EstablishGRPCConnection(r.ClientIP)
	if err != nil {
		log.Errorf("Cannot connect to address %s", r.ClientIP)
		return
	}
	// TODO Consider close grpc connection before restore gets done
	defer conn.Close()
	c := protoclnt.NewBaclntClient(conn)
	response, err := c.Restore(context.Background(),
		&protoclnt.RestoreRequest{Ip: r.ClientIP,
			PathsOnServer: r.FilesLocationOnServer, OriginalPaths: r.OriginalFilesLocations})
	if err != nil {
		log.Warningf("Could not get response from restore request, err: ", err)
		r.Status = false
		return
	}
	r.Status = true
	log.Println("Response: ", response)
}

func (r *Restore) Stop() {
	log.Println("Stopping")
}

func (r *Restore) Setup() error {
	backupMetadata, err := db.Get().GetBackupMetadata(r.BackupID)
	if err != nil {
		return err
	}
	for _, v := range backupMetadata.FilesMetadata {
		r.FilesLocationOnServer = append(r.FilesLocationOnServer, filepath.Join(backupMetadata.SavesetPath, v.OriginalFileLocation))
		r.OriginalFilesLocations = append(r.OriginalFilesLocations, v.OriginalFileLocation)
	}
	log.Debugln("Local on server backup files: ", r.FilesLocationOnServer)
	log.Debugln("Original files paths: ", r.OriginalFilesLocations)
	return nil
}
