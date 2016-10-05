package repository

import (
	log "github.com/Sirupsen/logrus"
	"os"
    "github.com/backer/bacsrv/config"
)


/*
REPOSITORY SCHEMA (TEMPORARY)

/.meta
/.meta/init
/data
/data/<client_name>
/locks

*/

func init(){
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}

func checkIfRepoExists() bool {
    repolocation := config.GetRepositoryLocalization()
    log.Debugf("Checking if %s repository exists...", repolocation)
    repo, err := os.Stat(repolocation)
        if err == nil && repo.IsDir() {
            // make more repository validations
            return true
        } else {
            return false
           
        }
}

func CreateRepository() (*Repository, error){
    repolocation := config.GetRepositoryLocalization()
    if checkIfRepoExists(){
        log.Infof("Repository %s exists, skipping creating", repolocation)
        return &Repository{Location: repolocation}, nil
    }
    err := os.MkdirAll(repolocation + "/.meta/init", 0700)
    if err != nil {
        log.Errorf("Cannot create repository %s...", repolocation)
        return nil, err
    }
    errd := os.MkdirAll(repolocation + "/data", 0700)
    if errd != nil {
        log.Error("Cannot create data directory inside repository")
        return nil, errd
    }

    erri := os.MkdirAll(repolocation + "/locks", 0700)
    if erri != nil {
        log.Error("Cannot create locks directory inside repository")
        return nil, erri
    }
    log.Infof("Repository %s has been created successfully", repolocation)
    return &Repository{Location: repolocation}, nil
}

