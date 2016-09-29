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

func checkIfRepoExists(){
    repolocation := config.GetRepositoryLocalization()
    log.Debugf("Checking if %s repository exists...", repolocation)
    repo, err := os.Stat(repolocation)
        if err == nil && repo.IsDir() {
            // make more repository validations
            log.Info("Repository exists, skipping creating.")
        } else {
            log.Info("Repository does not exist, creating...")
            out, err := createRepository(repolocation)
            if err != nil {
                log.Error("An error during creating repository, error: ", err.Error())
            }
            log.Info(out)
        }
}

func createRepository(repolocation string)(string, error){

    return "nil", nil
}