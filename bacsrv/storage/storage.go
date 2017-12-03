package storage

import (
	"os"
	"path/filepath"

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

func WriteBackupMetadata(data []byte) error {
	dbLocation := filepath.Join(config.MainConfig.Storage.Location, "/.meta/db")
	log.Debugln("Creating backup metadata file")
	file, err := os.Create(filepath.Join(dbLocation, "backupMeta.json"))
	defer file.Close()
	if err != nil {
		return err
	}
	wrote, err := file.Write(data)
	log.Debugln("Wrote backup metadata: ", wrote)
	return nil
}

func CheckIfFileExists(fullFilePath string) bool {
	if _, err := os.Stat(fullFilePath); err != nil {
		log.Print(err)
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetFileSize(fullPath string) int64 {
	file, err := os.Open(fullPath)
	defer file.Close()
	fstat, err := file.Stat()
	if err != nil {
		log.Println("Cannot do stat on file, returning 0")
		return 0
	}
	return fstat.Size()
}
