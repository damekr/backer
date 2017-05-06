package db

import (
	"github.com/HouzuoGuo/tiedot/db"
	"github.com/HouzuoGuo/tiedot/dberr"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"path/filepath"
)

const (
	clientsDocName = "clients"
	dbLocation     = ".meta"
)

func InitDBs() (*db.DB, error) {
	repoLocation := config.GetMainRepositoryLocation()
	DB, err := db.OpenDB(filepath.Join(repoLocation, dbLocation))
	if dberr.Type(err) == dberr.ErrorIO {
		log.Error("Cannot create DB")
		return nil, err
	}
	if err := createClientsDBDoc(DB); err != nil {
		log.Error("Cannot fully initialize DB because of clients DOC")
	}
	return DB, nil
}

func createClientsDBDoc(db *db.DB) error {
	log.Debug("Creating clients doc in DB")
	err := db.Create(clientsDocName)
	if dberr.Type(err) == dberr.ErrorIO {
		log.Warning("IO Error during creating Clients DOC")
		return nil
	}
	return nil
}

func OpenDB() (*db.DB, error) {
	repoLocation := config.GetMainRepositoryLocation()
	db, err := db.OpenDB(filepath.Join(repoLocation, dbLocation))
	if dberr.Type(err) == dberr.ErrorIO {
		log.Error("Cannot create DB")
		return nil, err
	}
	return db, nil
}
