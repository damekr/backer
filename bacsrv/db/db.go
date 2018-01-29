package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/damekr/backer/bacsrv/config"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"prefix": "db"})
)

type BackupMetadata struct {
	ClientName    string `json:"clientName"`
	BackupID      int    `json:"backupID"`
	BucketPath    string `json:"bucketLocation"`
	SavesetPath   string `json:"savesetLocation"`
	FilesMetadata []FileMetaData
}

type FileMetaData struct {
	FileWithPath string `json:"fileWithPath"`
	BackupTime   string `json:"backupTime"`
}

type DB struct {
	Location   string
	FileSchema BackupMetadata
}

func Get() DB {
	return DB{
		Location:   filepath.Join(config.MainConfig.Storage.Location, ".meta/db"),
		FileSchema: BackupMetadata{},
	}
}

func (d DB) createClientMetaCatalogue(clientName string) (string, error) {
	dbLocation := filepath.Join(config.MainConfig.Storage.Location, "/.meta/db")
	clientDbLocation := filepath.Join(dbLocation, clientName)
	if err := os.MkdirAll(clientDbLocation, 0700); err != nil {
		return "", err
	}
	return clientDbLocation, nil

}

func (d DB) WriteBackupMetadata(data []byte, fileName, clientName string) error {
	clientMetadataDbLocation, err := d.createClientMetaCatalogue(clientName)
	if err != nil {
		log.Errorln("Cannot create client DB location, err: ", err.Error())
		return err
	}
	log.Debugln("Creating backup metadata file")
	file, err := os.Create(filepath.Join(clientMetadataDbLocation, fileName) + ".json")
	defer file.Close()
	if err != nil {
		return err
	}
	wrote, err := file.Write(data)
	log.Debugln("Wrote backup metadata: ", wrote)
	return nil
}

func (d DB) GetClientsNames() []string {
	var clientNames []string
	files, err := ioutil.ReadDir(d.Location)
	if err != nil {
		log.Errorln("Cannot read files from db")
	}
	for _, v := range files {
		if v.IsDir() {
			clientNames = append(clientNames, v.Name())
		}
	}
	return clientNames
}

func (d DB) GetBackupsMetadata() []BackupMetadata {
	var clientsAssets []BackupMetadata
	clientsNames := d.GetClientsNames()
	for _, v := range clientsNames {
		clientAsset := d.GetClientBackupsMetadata(v)
		clientsAssets = append(clientsAssets, clientAsset...)
	}
	return clientsAssets
}

func (d DB) GetClientBackupsMetadata(clientName string) []BackupMetadata {
	var clientAssets []BackupMetadata
	files, err := ioutil.ReadDir(filepath.Join(d.Location, clientName))
	if err != nil {
		log.Warningln("Could not find any client assets")
	}
	for _, v := range files {
		clientAssets = append(clientAssets, d.readAsset(filepath.Join(d.Location, clientName, v.Name())))
	}
	return clientAssets
}

func (d DB) readAsset(filePath string) BackupMetadata {
	var asset BackupMetadata
	raw, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Error("Cannot read client asset, err: ", err.Error())
		return asset
	}
	err = json.Unmarshal(raw, &asset)
	if err != nil {
		log.Error("Could not unmarshal json file, err: ", err.Error())
	}
	return asset
}
