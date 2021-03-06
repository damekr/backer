package local

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	"github.com/damekr/backer/cmd/bacsrv/config"
	"github.com/damekr/backer/pkg/bftp"
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
	log.Infof("Storage %s has been created successfully or already exists", location)
	local := &Local{Location: location}
	return local, nil
}

func (l Local) CreateSaveset(bucketLocation string, assetID int) (string, error) {
	savesetName := "fullbackup" + "_" + strconv.Itoa(assetID)
	log.Debug("Creating saveset: ", savesetName)
	savesetLocation := filepath.Join(bucketLocation, savesetName)
	err := os.MkdirAll(savesetLocation, 0700)
	if err != nil {
		log.Error("Error occured during creation saveset, error: ", err.Error())
		return "", err
	}

	return savesetLocation, nil
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

func (l Local) ReadFile(filePath string) (io.ReadCloser, error) {
	log.Debugln("Reading file: ", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorln("Cannot open file for reading, err: ", err)
		return nil, err
	}
	return io.ReadCloser(file), nil
}

func (l Local) CreateFile(savesetLocation, fileOriginalPath string) (*os.File, error) {
	fileName := filepath.Base(fileOriginalPath)
	log.Infof("Creating file: %s, in saveset: %s", fileName, savesetLocation)

	file, err := os.Create(filepath.Join(savesetLocation, fileOriginalPath))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (l Local) CreateDir(savesetLocation string, dirMetadata bftp.DirMetadata) error {
	log.Debugf("Creating dir: %s under saveset: %s\n", dirMetadata.Path, savesetLocation)
	dirFullPath := filepath.Join(savesetLocation, dirMetadata.Path)
	err := os.MkdirAll(dirFullPath, dirMetadata.Mode)
	if err != nil {
		log.Errorf("Cannot create dir: %s under saveset: %s\n", dirMetadata, savesetLocation)
		return err
	}
	return nil
}

func (l Local) RemoveBucket(clientName string) {

}
