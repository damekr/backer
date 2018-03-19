package storage

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"

)

type Local struct {
	Location string
}

func Create(location string) (*Local, error) {

	err := os.MkdirAll(filepath.Join(location, "bftp"), 0700)
	if err != nil {
		log.Println("Cannot create storage: ", location)
		return nil, err
	}
	local := &Local{Location: location}
	return local, nil
}

func (l Local) CreateBucket(clientName string) (string, error) {
		return "", nil
}

func (l Local) RemoveBucket(clientName string) {

}

func (l Local) CreateSaveset(bucketLocation string) (string, error) {

	return "", nil
}

func (l Local) CreateFile(savesetLocation, fileName string) (*os.File, error) {
	log.Println("Creating file: ", fileName)
	file, err := os.Create(path.Join(savesetLocation, fileName))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l Local) OpenFile(fileLocation string) (*os.File, error) {
	log.Println("Opening file: ", fileLocation)
	file, err := os.Open(fileLocation)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func CheckIfFileExists(fullFilePath string) bool {
	if _, err := os.Stat(fullFilePath); err != nil {
		fmt.Print(err)
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetFileSize(fullPath string) int64 {
	file, err := os.Open(fullPath)
	defer file.Close()
	fstat, err := file.Stat()
	if err != nil {
		log.Println("Cannot do stat on file, returning 0")
		return 0
	}
	return fstat.Size()
}
