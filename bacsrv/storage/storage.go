package storage

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"os"
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

type Backend interface {
	SaveFile()
	RemoveFile()
}


type Storage struct {
	Type Backend
	// Size uint64
}


func Create(storageType string) (Backend, error){
	switch storageType {
	case "local":
		//TODO Get location from config
		localStorage, err := local.Create("/tmp")
		if err != nil {
			log.Error("Cannot create local storage")
		}
		return localStorage, nil

	}
}


func checkIfRepoExists(repolocation string) bool {
	log.Debugf("Checking if %s storage exists...", repolocation)
	repo, err := os.Stat(repolocation)
	if err == nil && repo.IsDir() {
		// TODO make more storage validations
		return true
	}
	return false
}


