package storage

import (
	"os"
	"log"
)

type Storage interface {
	CreateBucket(clientName string) (string, error)
	CreateSaveset(bucketLocation string) (string, error)
	CreateFile(savesetLocation, fileName string) (*os.File, error)
	OpenFile(fileLocation string) (*os.File, error)
}



func NewStorage(storageType, location string) (Storage, error) {
	switch storageType {
	case "local":
		localStorage, err := Create(location)
		if err != nil {
			log.Println("Cannot create local storage")
			return nil, err
		}
		return localStorage, nil
	}
	//FIXME
	return nil, nil
}