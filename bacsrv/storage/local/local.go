package local

import (
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/damekr/backer/bacsrv/config"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var logger = logrus.New()
var log = &logrus.Entry{}

func init() {
	logger.Formatter = new(prefixed.TextFormatter)
	logger.Level = logrus.DebugLevel
	log = logger.WithFields(logrus.Fields{"prefix": "storage:local"})

}

type Local struct {
	Location string
}

func Create(location string) (*Local, error) {

	err := os.MkdirAll(location+"/.meta/init", 0700)
	if err != nil {
		log.Errorf("Cannot create storage %s...", location)
		return nil, err
	}
	errd := os.MkdirAll(location+"/data", 0700)
	if errd != nil {
		log.Error("Cannot create data directory inside storage")
		return nil, errd
	}

	erri := os.MkdirAll(location+"/locks", 0700)
	if erri != nil {
		log.Error("Cannot create locks directory inside storage")
		return nil, erri
	}
	errdb := os.MkdirAll(location+"/.meta/db", 0700)
	if erri != nil {
		log.Error("Cannot create dbs directory inside storage")
		return nil, errdb
	}
	log.Infof("Storage %s has been created successfully", location)
	local := &Local{Location: location}
	return local, nil
}

func (l Local) CreateBucket(clientName string) (string, error) {
	log.Debugln("Creating client bucket: ", clientName)
	//TODO If bucket exist, just return it. Should be handled already by MkdirAll
	bucketLocation := filepath.Join(config.MainConfig.Storage.Location, "data", clientName)
	err := os.MkdirAll(bucketLocation, 0700)
	if err != nil {
		log.Errorln("Could not create client bucket")
		return "", err
	}
	return bucketLocation, nil
}

func (l Local) OpenFile(fileLocation string) (*os.File, error) {
	log.Println("Opening file: ", fileLocation)
	file, err := os.Open(fileLocation)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (l Local) RemoveBucket(clientName string) {

}

func (l Local) CreateSaveset(bucketLocation string) (string, error) {
	savesetName := "fullbackup" + "_" + strconv.Itoa(time.Now().Nanosecond())
	log.Debug("Creating saveset: ", savesetName)
	savesetLocation := filepath.Join(bucketLocation, savesetName)
	err := os.MkdirAll(savesetLocation, 0700)
	if err != nil {
		log.Error("Error occured during creation saveset, error: ", err.Error())
		return "", err
	}

	return savesetLocation, nil
}

func createFileOriginalPath(savesetLocation, filePath string) (string, error) {
	log.Infof("Creating path: %s under saveset: %s\n", filePath, savesetLocation)
	fullFilePath := filepath.Join(savesetLocation, filePath)
	err := os.MkdirAll(fullFilePath, 0700)
	if err != nil {
		log.Errorf("Cannot create path: %s under saveset: %s\n", filePath, savesetLocation)
		return "", err
	}
	return fullFilePath, nil
}

func (l Local) CreateFile(savesetLocation, fileOriginalPath string) (*os.File, error) {
	filePath := filepath.Dir(fileOriginalPath)
	fileName := filepath.Base(fileOriginalPath)
	log.Infof("Creating file: %s, in saveset: %s", fileName, savesetLocation)
	fullFilePath, err := createFileOriginalPath(savesetLocation, filePath)
	if err != nil {
		log.Errorln("Cannot create path for file backup")
	}
	file, err := os.Create(filepath.Join(fullFilePath, fileName))
	if err != nil {
		return nil, err
	}
	return file, nil
}
