package repository

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/clientsconfig"
	"os"
)

type ClientBucket struct {
	Name string
}

func checkIfClientBucketExists(name string) bool {
	repolocation := MainRepository.Location
	log.Debugf("Checking if %s  bucket exists, under mainrepository: %s", name, repolocation)
	repo, err := os.Stat(repolocation + name)
	if err == nil && repo.IsDir() {
		// TODO make more validations
		log.Infof("Client %s bucket exists", name)
		return true
	}
	return false
}

func InitClientsBuckets() error {
	allClients := clientsconfig.GetAllClients()
	log.Debug("Number of integrated clients: ", len(allClients))
	for _, v := range allClients {
		log.Printf("Client info: %s", v.Name)
		if !checkIfClientBucketExists(v.Name) {
			log.Infof("Client %s bucket does not exist, creating...", v.Name)
			_ = CreateClient(v.Name)
		}
	}
	return nil
}

func CreateClient(name string) *ClientBucket {
	repo := GetRepository()
	repo.CreateClientBucket(name)
	return &ClientBucket{Name: name}
}
