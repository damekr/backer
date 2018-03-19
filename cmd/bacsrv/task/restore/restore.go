package restore

import (
	"context"

	"github.com/damekr/backer/api/protoclnt"
	"github.com/damekr/backer/cmd/bacsrv/db"
	"github.com/damekr/backer/cmd/bacsrv/network"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "task:restore"})

type Restore struct {
	ClientIP       string `json:"clientIP"`
	BackupID       int    `json:"backupID"`
	FilesMetadata  []bftp.FileMetadata
	Progress       int    `json:"-"`
	Status         bool   `json:"status"`
	BucketLocation string `json:"bucketLocation"`
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
	var restoreFilesInfo []*protoclnt.RestoreFileInfo
	log.Debugln("FIles metadata: ", r.FilesMetadata)
	for _, v := range r.FilesMetadata {
		fileMeta := protoclnt.RestoreFileInfo{
			LocationOnServer: v.LocationOnServer,
			OriginalLocation: v.OriginalFileLocation,
		}
		restoreFilesInfo = append(restoreFilesInfo, &fileMeta)
	}
	log.Debugln("Restore Files info: ", restoreFilesInfo)
	response, err := c.Restore(context.Background(),
		&protoclnt.RestoreRequest{Ip: r.ClientIP,
			RestoreFileInfo: restoreFilesInfo})

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

//Setup configures restore job, should be splited into different kind of setups(singleDir, wholeBackup etc.).
func (r *Restore) Setup(remotePath string, singleDirPath string) error {
	backupMetadata, err := db.Get().GetBackupMetadata(r.BackupID)
	if err != nil {
		return err
	}

	if singleDirPath != "" {
		for _, v := range backupMetadata.FilesMetadata {
			if v.OriginalFileLocation == singleDirPath {
				log.Infoln("Adding to restore single dir: ", v.OriginalFileLocation)
				r.FilesMetadata = append(r.FilesMetadata, v)
			}
		}
	} else {
		for _, v := range backupMetadata.FilesMetadata {
			r.FilesMetadata = append(r.FilesMetadata, v)
		}
	}

	log.Debugln("Files to be restored metadata: ", r.FilesMetadata)
	return nil
}
