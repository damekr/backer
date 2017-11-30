package storage

import (
	"os"

	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/storage/local"
	"github.com/sirupsen/logrus"
)

/*
REPOSITORY SCHEMA (TEMPORARY)

/.meta
/.meta/db
/.meta/init
/data
/data/<client_name>
/locks

*/

var log = logrus.WithFields(logrus.Fields{"prefix": "storage"})

var DefaultStorage Storage

type Storage interface {
	CreateBucket(clientName string) (string, error)
	CreateSaveset(bucketLocation string) (string, error)
	CreateFile(savesetLocation, fileName string) (*os.File, error)
	OpenFile(fileLocation string) (*os.File, error)
}

func setDefaultStorage(storage Storage) {
	DefaultStorage = storage
}

func Create(storageType string) error {
	switch storageType {
	case "local":
		localStorage, err := local.Create(config.MainConfig.Storage.Location)
		if err != nil {
			log.Error("Cannot create local storage")
			return err
		}
		setDefaultStorage(localStorage)
		return nil

	}
	return nil
}
