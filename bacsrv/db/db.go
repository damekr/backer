package db

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/damekr/backer/bacsrv/config"
	"os"
	"path/filepath"
)

const (
	clientsBucket = "clients"
)

var (
	db *bolt.DB
)

func createClientsBucket() error {
	var err error
	db.Update(func(tx *bolt.Tx) error {
		_, err = tx.CreateBucketIfNotExists([]byte(clientsBucket))
		if err != nil {
			return fmt.Errorf("create bucket: %s", err)
		}
		return nil
	})
	return nil
}

func InitDB() {
	var err error
	db, err = bolt.Open(filepath.Join(config.GetDBLocation(), "server.db"), 0600, nil)
	if err != nil {
		log.Fatal(err)
		os.Exit(4)
	}
	err = createClientsBucket()
	if err != nil {
		log.Fatal("Fatal error during creating clients bucket in DB")
		os.Exit(5)
	}
}
