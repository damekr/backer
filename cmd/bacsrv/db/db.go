package db

import (
	"path/filepath"

	"github.com/damekr/backer/cmd/bacsrv/config"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const (
	dbHiddenPath = ".meta/db"
)

var (
	log                    = logrus.WithFields(logrus.Fields{"prefix": "db"})
	clientMetadataNotFound = errors.New("not found client metadata")
)

type BackupsDB interface {
	CreateAssetMetadata(assetMetadata bftp.AssetMetadata) error
	ReadClientsNames() ([]string, error)
	ReadAssetsMetadata() ([]bftp.AssetMetadata, error)
	ReadAssetsMetadataOfClient(clientName string) ([]bftp.AssetMetadata, error)
	ReadAssetMetadata(assetID int) (*bftp.AssetMetadata, error)
}

func DB() BackupsDB {
	// IN case of use different DBs
	dbLocation := filepath.Join(config.MainConfig.DBLocation, dbHiddenPath)
	return GetJsonBackupDB(dbLocation)
}
