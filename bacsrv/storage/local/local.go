package local

import (
	"io"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"

	"github.com/damekr/backer/bacsrv/config"
	"github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
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

type ClientBucket struct {
	Location string
}

type ClientSaveset struct {
	Bucket   ClientBucket
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

func (l Local) CreateBucket(clientName string) *ClientBucket {
	log.Debugln("Creating client bucket: ", clientName)
	bucketLocation := filepath.Join(config.MainConfig.Storage.Location, "data", clientName)
	err := os.MkdirAll(bucketLocation, 0700)
	if err != nil {
		log.Errorln("Could not create client bucket")
		return nil
	}
	return &ClientBucket{
		Location: bucketLocation,
	}
	return nil
}

func (l Local) RemoveBucket(clientName string) {

}

func (b *ClientBucket) CreateSaveset() *ClientSaveset {
	savesetName := "fullbackup" + "_" + strconv.Itoa(time.Now().Nanosecond())
	log.Debug("Creating saveset: ", savesetName)
	savesetLocation := filepath.Join(b.Location, savesetName)
	err := os.MkdirAll(savesetLocation, 0700)
	if err != nil {
		log.Error("Error occured during creation saveset, error: ", err.Error())
	}
	saveset := ClientSaveset{
		Bucket:   *b,
		Location: savesetLocation,
	}
	return &saveset
}

func (b *ClientBucket) RemoveSaveset() {

}

func (s *ClientSaveset) CreateFile(fileName string) (io.Writer, error) {
	log.Infoln("Creating file: ", fileName)
	file, err := os.Create(path.Join(s.Location, fileName))
	if err != nil {
		return nil, err
	}

	return file, nil
}

func (s *ClientSaveset) RemoveFile(saveset *ClientSaveset) {

}
