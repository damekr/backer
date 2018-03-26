package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/damekr/backer/cmd/bacsrv/config"
	"github.com/damekr/backer/pkg/bftp"
)

const (
	dbFilesSuffix      = ".json"
	dbFilesPermissions = 0700
)

type Json struct {
	DBLocation string
}

func GetJsonsBackupDB(dbLocation string) Json {
	return Json{
		// TODO It should be applied at the beginning of application start
		DBLocation: dbLocation,
	}
}

func (j Json) CreateBackupMetadata(backupMetadata bftp.BackupMetaData) error {
	clientMetadataLocation, err := j.createClientMetaCatalogue(backupMetadata.ClientName)
	if err != nil {
		log.Errorln("Cannot create client metadata location, err: ", err.Error())
		return err
	}

	log.Debugln("Creating backup metadata file")
	file, err := os.Create(filepath.Join(clientMetadataLocation, filepath.Base(backupMetadata.SavesetPath)) + dbFilesSuffix)
	defer file.Close()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(backupMetadata)
	if err != nil {
		return err
	}
	wrote, err := file.Write(jsonData)
	if err != nil {
		return err
	}
	log.Debugln("Wrote backup metadata: ", wrote)
	return nil
}

func (j Json) createClientMetaCatalogue(clientName string) (string, error) {
	dbLocation := filepath.Join(config.MainConfig.Storage.Location, "/.meta/db")
	log.Debugln("DB Assets DBLocation: ", dbLocation)
	clientDbLocation := filepath.Join(dbLocation, clientName)
	if err := os.MkdirAll(clientDbLocation, dbFilesPermissions); err != nil {
		return "", err
	}
	return clientDbLocation, nil

}

func (j Json) ReadBackupsMetadata() ([]BackupMetadata, error) {
	var backupsMetadata []BackupMetadata
	clientsNames, err := j.ReadClientsNames()
	if err != nil {
		return backupsMetadata, err
	}
	for _, v := range clientsNames {
		clientAsset, err := j.ReadBackupsMetadataOfClient(v)
		if err != nil {
			log.Errorln("Could not read backup metadata of client: ", v)
		} else {
			backupsMetadata = append(backupsMetadata, clientAsset...)
		}
	}
	return backupsMetadata, nil

}
func (j Json) ReadBackupsMetadataOfClient(clientName string) ([]BackupMetadata, error) {
	files, err := ioutil.ReadDir(filepath.Join(j.DBLocation, clientName))
	if err != nil {
		return nil, clientMetadataNotFound
	}

	backupsMetadata := make([]BackupMetadata, len(files))

	for _, v := range files {
		backupMetadata, err := j.readClientBackupMetadata(filepath.Join(j.DBLocation, clientName, v.Name()))
		if err != nil {
			log.Errorln(err)
		} else {
			backupsMetadata = append(backupsMetadata, *backupMetadata)
		}
	}
	return backupsMetadata, nil

}

func (j Json) readClientBackupMetadata(clientAssetPath string) (*BackupMetadata, error) {
	rawData, err := ioutil.ReadFile(clientAssetPath)
	if err != nil {
		return nil, err
	}
	backupMetadata := new(BackupMetadata)
	err = json.Unmarshal(rawData, backupMetadata)
	if err != nil {
		log.Error("Could not unmarshal json file, err: ", err.Error())
	}
	return backupMetadata, nil
}

func (j Json) ReadClientsNames() ([]string, error) {
	var clientNames []string
	files, err := ioutil.ReadDir(j.DBLocation)
	if err != nil {
		log.Errorln("Cannot read files from db, err: ", err)
		return clientNames, err
	}
	for _, v := range files {
		if v.IsDir() {
			clientNames = append(clientNames, v.Name())
		}
	}
	return clientNames, nil

}

func (j Json) ReadBackupMetadata(backupID int) (*BackupMetadata, error) {

	backupsMetadata, err := j.ReadBackupsMetadata()
	if err != nil {
		return nil, err
	}
	seekingBackupMetadata := new(BackupMetadata)
	for _, v := range backupsMetadata {
		log.Println("Backup id: ", v.BackupID)
		if v.BackupID == backupID {
			seekingBackupMetadata = &v
		}
	}
	if seekingBackupMetadata.BackupID == 0 {
		return seekingBackupMetadata, errors.New("Backup metadata not found")
	}
	return seekingBackupMetadata, nil
}
