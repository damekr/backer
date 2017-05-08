package db

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/tiedot/db"
	"github.com/damekr/backer/tiedot/dberr"
	"path/filepath"
)

const (
	clientsDocName = "clients"
	backupsDocName = "backups"
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
	if err := createBackupsDBDoc(DB); err != nil {
		log.Error("Cannot fully initialize DB because of clients DOC")
	}
	return DB, nil
}

func createClientsDBDoc(db *db.DB) error {
	log.Debug("Creating clients doc in DB")
	err := db.Create(clientsDocName)
	if dberr.Type(err) == dberr.ErrorIO {
		log.Warning("IO Error during creating Clients DOC")
		return err
	}
	return nil
}

func createBackupsDBDoc(db *db.DB) error {
	log.Debug("Creating backups doc in DB")
	err := db.Create(backupsDocName)
	if dberr.Type(err) == dberr.ErrorIO {
		log.Warning("IO Error during creating Backups DOC")
		return err
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
