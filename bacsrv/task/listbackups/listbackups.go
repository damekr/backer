package listbackups

import (
	"fmt"

	"github.com/damekr/backer/bacsrv/db"
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
	database := db.Get()
	clientAssets := database.GetClientAssets(l.ClientName)
	for _, v := range clientAssets {
		backupIds = append(backupIds, int64(v.BackupID))
	}
	l.BackupIDs = backupIds
	log.Debugln("Found client backups ids: ", backupIds)

}

func (l *ListBackups) Stop() {
	fmt.Println("Stopping")
}
