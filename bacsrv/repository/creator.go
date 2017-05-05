package repository

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"os"
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

type Repository struct {
	Location string
	// Size uint64
}

var (
	MainRepository *Repository
)

func checkIfRepoExists(repolocation string) bool {
	log.Debugf("Checking if %s repository exists...", repolocation)
	repo, err := os.Stat(repolocation)
	if err == nil && repo.IsDir() {
		// TODO make more repository validations
		return true
	}
	return false
}

func CreateRepository() (*Repository, error) {
	repolocation := config.GetMainRepositoryLocation()
	if checkIfRepoExists(repolocation) {
		log.Infof("Repository %s exists, skipping creating", repolocation)
		MainRepository = &Repository{Location: repolocation}
		return MainRepository, nil
	}
	err := os.MkdirAll(repolocation+"/.meta/init", 0700)
	if err != nil {
		log.Errorf("Cannot create repository %s...", repolocation)
		return nil, err
	}
	errd := os.MkdirAll(repolocation+"/data", 0700)
	if errd != nil {
		log.Error("Cannot create data directory inside repository")
		return nil, errd
	}

	erri := os.MkdirAll(repolocation+"/locks", 0700)
	if erri != nil {
		log.Error("Cannot create locks directory inside repository")
		return nil, erri
	}
	errdb := os.MkdirAll(repolocation+"/.meta/db", 0700)
	if erri != nil {
		log.Error("Cannot create dbs directory inside repository")
		return nil, errdb
	}
	log.Infof("Repository %s has been created successfully", repolocation)
	MainRepository = &Repository{Location: repolocation}
	return MainRepository, nil
}
