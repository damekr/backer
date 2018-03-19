package db

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/damekr/backer/cmd/bacsrv/config"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/pkg/errors"
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
	FilesMetadata []bftp.FileMetadata
}

type DB struct {
	AssetsLocation string // Means clients metadata location, where all backups related information exist
	FileSchema     BackupMetadata
}

func Get() DB {
	return DB{
		AssetsLocation: filepath.Join(config.MainConfig.Storage.Location, ".meta/db"),
		FileSchema:     BackupMetadata{},
	}
}

func (d DB) createClientMetaCatalogue(clientName string) (string, error) {
	dbLocation := filepath.Join(config.MainConfig.Storage.Location, "/.meta/db")
	log.Debugln("DB Assets Location: ", dbLocation)
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
	files, err := ioutil.ReadDir(d.AssetsLocation)
	if err != nil {
		log.Errorln("Cannot read files from db, err: ", err)
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
	files, err := ioutil.ReadDir(filepath.Join(d.AssetsLocation, clientName))
	if err != nil {
		log.Warningln("Could not find any client assets")
	}
	for _, v := range files {
		clientAssets = append(clientAssets, d.readClientAssets(filepath.Join(d.AssetsLocation, clientName, v.Name())))
	}
	return clientAssets
}

func (d DB) GetBackupMetadata(backupID int) (BackupMetadata, error) {
	var seekingBackupMetadata BackupMetadata
	backupsMetadata := d.GetBackupsMetadata()
	for _, v := range backupsMetadata {
		log.Println("Backup id: ", v.BackupID)
		if v.BackupID == backupID {
			seekingBackupMetadata = v
		}
	}
	if seekingBackupMetadata.BackupID == 0 {
		return seekingBackupMetadata, errors.New("Backup metadata not found")
	}
	return seekingBackupMetadata, nil
}

func (d DB) readClientAssets(clientAssetPath string) BackupMetadata {
	var asset BackupMetadata
	raw, err := ioutil.ReadFile(clientAssetPath)
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
