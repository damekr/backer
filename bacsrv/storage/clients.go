package storage

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

type ClientBucket struct {
	Location string
}

func checkIfClientBucketExists(name string) bool {
	repolocation := Storage.Location
	log.Debugf("Checking if %s  bucket exists, under mainrepository: %s", name, repolocation)
	bucketFolder := filepath.Join(repolocation, bucketsLocation, name)
	log.Debug("Checking file bucket as foleder: ", bucketFolder)
	repo, err := os.Stat(bucketFolder)
	if err == nil && repo.IsDir() {
		// TODO make more validations
		log.Infof("Client %s bucket exists", name)
		return true
	}
	return false
}

func InitClientsBuckets() error {
	repo := GetRepository()
	allClients := config.GetAllClients()
	log.Debug("Number of integrated clients: ", len(allClients))
	for _, v := range allClients {
		log.Printf("Client info: %s", v.Name)
		if !checkIfClientBucketExists(v.Name) {
			log.Infof("Client %s bucket does not exist, creating...", v.Name)
			err := repo.CreateClientBucket(v.Name)
			if err != nil {
				log.Errorf("Could not create client %s bucket", v.Name)
			}
		}
	}
	return nil
}

func CreateClientSaveset(name string) (string, error) {
	repo := GetRepository()
	log.Debugf("Getting client %s bucket", name)
	bucket, err := repo.GetClientBucket(name)
	if err != nil {
		log.Errorf("Cannot create saveset because client %s bucket does not exist", name)
		return "", errors.New("Clients bucket does not exist")
	}
	savesetName := "fullbackup" + "_" + strconv.Itoa(time.Now().Nanosecond()) + "_" + name
	log.Debug("Creating saveset: ", savesetName)
	err = os.MkdirAll(filepath.Join(bucket.Location, savesetName), 0700)
	if err != nil {
		log.Error("Error occured during creation saveset, error: ", err.Error())
	}
	return filepath.Join(bucket.Location, savesetName), nil
}
