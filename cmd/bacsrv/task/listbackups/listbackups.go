package listbackups

import (
	"fmt"

	"github.com/damekr/backer/cmd/bacsrv/db"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "task:listbackups"})

type ListBackups struct {
	ClientName string  `json:"clientName"`
	BackupIDs  []int64 `json:"backupIDs"`
}

//TODO Move it, its internal operation hence should not be here
func Create(clientName string) *ListBackups {
	return &ListBackups{
		ClientName: clientName,
	}
}

func (l *ListBackups) Run() {
	log.Println("Getting backups")
	var backupIds []int64
	database := db.DB()
	clientAssets, err := database.ReadBackupsMetadataOfClient(l.ClientName)
	if err != nil {
		log.Errorln("Could not read backups metadata of client, err: ", err)
	} else {
		for _, v := range clientAssets {
			backupIds = append(backupIds, int64(v.BackupID))
		}
		l.BackupIDs = backupIds
		log.Debugln("Found client backups ids: ", backupIds)
	}
}

func (l *ListBackups) Stop() {
	fmt.Println("Stopping")
}
