package db

import (
	"github.com/HouzuoGuo/tiedot/db"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/clientsconfig"
	"github.com/damekr/backer/bacsrv/config"
	"path/filepath"
)

const (
	clientsDBName     = "clientsDB"
	clientsDBLocation = ".meta/db"
)

func init() {
	log.Debug("Initializing db module")
}

func GetClientDBConnection() *db.DB {
	repoLocation := config.GetMainRepositoryLocation()
	clientsDB, err := db.OpenDB(filepath.Join(repoLocation, clientsDBLocation, clientsDBName))
	if err != nil {
		log.Error("Cannot create or open clients DB")
	}
	return clientsDB
}

func AddClient(client *clientsconfig.Client) error {
	dbConnection := GetClientDBConnection()
	log.Debugf("Adding client: %s into database ", client.Name)
	if err := dbConnection.Create("Feeds"); err != nil {
		panic(err)
	}
	return nil
}
