package db

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"

	"github.com/damekr/backer/pkg/bftp"
)

const (
	dbFilesSuffix      = ".json"
	dbFilesPermissions = 0700
)

type Json struct {
	DBLocation string
}

func GetJsonBackupDB(dbLocation string) Json {
	return Json{
		// TODO It should be applied at the beginning of application start
		DBLocation: dbLocation,
	}
}

func (j Json) CreateAssetMetadata(assetMetadata bftp.AssetMetadata) error {
	clientMetadataLocation, err := j.createClientMetaCatalogue(assetMetadata.ClientName)
	if err != nil {
		log.Errorln("Cannot create client metadata location, err: ", err.Error())
		return err
	}

	log.Debugln("Creating backup metadata file")
	s := strconv.Itoa(assetMetadata.ID)
	file, err := os.Create(filepath.Join(clientMetadataLocation, s+dbFilesSuffix))
	defer file.Close()
	if err != nil {
		return err
	}
	jsonData, err := json.Marshal(assetMetadata)
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
	log.Debugln("DB Assets DBLocation: ", j.DBLocation)
	clientDbLocation := filepath.Join(j.DBLocation, clientName)
	if err := os.MkdirAll(clientDbLocation, dbFilesPermissions); err != nil {
		return "", err
	}
	return clientDbLocation, nil

}

func (j Json) ReadAssetsMetadata() ([]bftp.AssetMetadata, error) {
	var backupsMetadata []bftp.AssetMetadata
	clientsNames, err := j.ReadClientsNames()
	if err != nil {
		return backupsMetadata, err
	}
	for _, v := range clientsNames {
		clientAsset, err := j.ReadAssetsMetadataOfClient(v)
		if err != nil {
			log.Errorln("Could not read backup metadata of client: ", v)
		} else {
			backupsMetadata = append(backupsMetadata, clientAsset...)
		}
	}
	return backupsMetadata, nil

}
func (j Json) ReadAssetsMetadataOfClient(clientName string) ([]bftp.AssetMetadata, error) {
	var backupsMetadata []bftp.AssetMetadata
	clientMetadataLocation := filepath.Join(j.DBLocation, clientName)
	log.Debugln("DBLOCATION: ", j.DBLocation)
	log.Debugln("Reading client metadata from location: ", clientMetadataLocation)
	files, err := ioutil.ReadDir(clientMetadataLocation)
	if err != nil {
		return nil, clientMetadataNotFound
	}
	for _, v := range files {
		log.Debugln("Checking json files for metadata: ", v)
		backupMetadata, err := j.readClientBackupMetadata(filepath.Join(j.DBLocation, clientName, v.Name()))
		if err != nil {
			log.Errorln("Cannot read client metadata file, err: ", err)
		} else {
			backupsMetadata = append(backupsMetadata, *backupMetadata)
		}
	}
	return backupsMetadata, nil

}

func (j Json) readClientBackupMetadata(clientAssetPath string) (*bftp.AssetMetadata, error) {
	rawData, err := ioutil.ReadFile(clientAssetPath)
	if err != nil {
		return nil, err
	}
	backupMetadata := new(bftp.AssetMetadata)
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

func (j Json) ReadAssetMetadata(backupID int) (*bftp.AssetMetadata, error) {
	seekingBackupMetadata := new(bftp.AssetMetadata)
	backupsMetadata, err := j.ReadAssetsMetadata()
	if err != nil {
		return seekingBackupMetadata, err
	}
	for _, v := range backupsMetadata {
		log.Println("Backup id: ", v.ID)
		if v.ID == backupID {
			seekingBackupMetadata = &v
		}
	}
	if seekingBackupMetadata.ID == 0 {
		return seekingBackupMetadata, errors.New("backup metadata not found")
	}
	return seekingBackupMetadata, nil
}
