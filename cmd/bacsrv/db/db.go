package db

import (
	"path/filepath"

	"github.com/damekr/backer/cmd/bacsrv/config"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

var (
	log                    = logrus.WithFields(logrus.Fields{"prefix": "db"})
	clientMetadataNotFound = errors.New("not found client metadata")
	dbLocation             = filepath.Join(config.MainConfig.Storage.Location, ".meta/db")
)

type BackupsDB interface {
	CreateBackupMetadata(backupMetadata bftp.BackupMetaData) error
	ReadClientsNames() ([]string, error)
	ReadBackupsMetadata() ([]BackupMetadata, error)
	ReadBackupsMetadataOfClient(clientName string) ([]BackupMetadata, error)
	ReadBackupMetadata(backupID int) (*BackupMetadata, error)
}

type BackupMetadata struct {
	ClientName    string `json:"clientName"`
	BackupID      int    `json:"backupID"`
	BucketPath    string `json:"bucketLocation"`
	SavesetPath   string `json:"savesetLocation"`
	FilesMetadata []bftp.FileMetadata
}

func Get() BackupsDB {
	// IN case of use different DBs
	return GetJsonsBackupDB(dbLocation)
}
