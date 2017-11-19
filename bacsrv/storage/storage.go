package storage

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/storage/local"

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

var DefaultStorage Backend

type Backend interface {
	CreateBucket(clientName string) *local.ClientBucket
	RemoveBucket(clientName string)
}

type Storage struct {
	Type Backend
	// Size uint64
}

func setDefaultStorage(storage Backend) {
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
