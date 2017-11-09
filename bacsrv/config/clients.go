package config

import (
	"errors"
	logger "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

var (
	EmptyClientsConfig = errors.New("ClientsConfigFilePath: Empty clients config file")
	AllClients         = []clientDefinition{}
	log                = logger.New()
)

func init() {
	log.WithFields(logger.Fields{"MODULE": "CLIENTS-CONFIG"})
}

type clientDefinition struct {
	Name                 string   `json:"clientName"`
	IPAddress            string   `json:"clientAddress"`
	BackupsIDs           []string `json:"backupsIDs"`
	BackupDefinitions    []backupDefinition
	SchedulesIDs         []string `json:"schedulesIds"`
	SchedulesDefinitions []scheduleDefinition
	Platform             string `json:"platform"`
	CID                  string `json:"cid"`
}

// backupDefinition specifies a backup
type backupDefinition struct {
	ID          string   `json:"id"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Paths       []string `json:"paths"`
	Excluded    []string `json:"excludedPaths"`
	Retention   string   `json:"retentionTime"`
}

type scheduleDefinition struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Hour        time.Time `json:"startHour"`
	Day         string    `json:"startDay"`
}

// TODO Must be places in communication protocol
type BackupTriggerMessage struct {
	ClientName   string `json:"clientName"`
	BackupConfig backupDefinition
}

func ReadInClientsConfig(clientsConfigPath, backupsConfigPath, scheduleConfigPath string) error {
	backupsViper, err := readInBackupsConfig(backupsConfigPath)
	if err != nil {
		return err
	}
	backupsDefinitions := setBackupsDefinitions(*backupsViper)

	schedulesViper, err := readInSchedulesConfig(scheduleConfigPath)
	if err != nil {
		return err
	}
	schedulesDefinitions := setSchedulesDefinitions(*schedulesViper)

	clientsViper, err := readInClientsConfig(clientsConfigPath)
	if err != nil {
		return err
	}
	clientsDefinitions := setClientsDefinitions(*clientsViper)
	for _, client := range clientsDefinitions {
		matchedClient, err := matchClientDefinition(&client, backupsDefinitions, schedulesDefinitions)
		if err != nil {
			AllClients = append(AllClients, *matchedClient)
		}
	}

	return nil
}

func matchClientDefinition(clntDefinition *clientDefinition, backupsDefinitions []backupDefinition, schedulesDefinitions []scheduleDefinition) (*clientDefinition, error) {
	for bk, bv := range backupsDefinitions {
		for _, bid := range clntDefinition.BackupsIDs {
			if bid == bv.ID {
				clntDefinition.BackupDefinitions = append(clntDefinition.BackupDefinitions, backupsDefinitions[bk])
			}
		}
	}
	for sk, sv := range schedulesDefinitions {
		for _, sid := range clntDefinition.SchedulesIDs {
			if sid == sv.ID {
				clntDefinition.SchedulesDefinitions = append(clntDefinition.SchedulesDefinitions, schedulesDefinitions[sk])
			}
		}
	}
	return clntDefinition, nil
}

func readInClientsConfig(clientsConfigPath string) (*viper.Viper, error) {
	log.Debugln("Clients config file path: ", clientsConfigPath)
	clientsViper := viper.New()
	clientsViper.SetConfigFile(clientsConfigPath)
	err := clientsViper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return clientsViper, nil
}

func readInBackupsConfig(backupsConfigPath string) (*viper.Viper, error) {
	log.Debugln("Backups config file backupsConfigPath: ", backupsConfigPath)
	backupsViper := viper.New()
	backupsViper.SetConfigFile(backupsConfigPath)
	err := backupsViper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return backupsViper, nil
}

func readInSchedulesConfig(scheduleConfigPath string) (*viper.Viper, error) {
	log.Debugln("Backups config file scheduleConfigPath: ", scheduleConfigPath)
	schedulesViper := viper.New()
	schedulesViper.SetConfigFile(scheduleConfigPath)
	err := schedulesViper.ReadInConfig()
	if err != nil {
		return nil, err
	}
	return schedulesViper, nil
}

func setBackupsDefinitions(backupsViper viper.Viper) []backupDefinition {
	if len(backupsViper.AllKeys()) == 0 {
		return []backupDefinition{}
	}
	var backupsDefinitions []backupDefinition
	for k, v := range backupsViper.AllSettings() {
		log.Debugf("Key: %s, value: %s", k, v)
		var bDefinition backupDefinition
		validDefinition := true
		backupPropertySlices := backupsViper.GetStringMapStringSlice(k)
		backup := backupsViper.GetStringMapString(k)
		if len(backup) == 0 {
			validDefinition = false
		}
		if backup["id"] != "" {
			bDefinition.ID = backup["id"]
		} else {
			validDefinition = false
		}
		if backup["name"] != "" {
			bDefinition.Name = backup["name"]
		} else {
			validDefinition = false
		}
		if backup["retention"] != "" {
			bDefinition.Retention = backup["retention"]
		} else {
			validDefinition = false
		}
		if backup["description"] != "" {
			bDefinition.Description = backup["description"]
		} else {
			validDefinition = false
		}
		if len(backupPropertySlices["paths"]) > 0 {
			bDefinition.Paths = backupPropertySlices["paths"]
		} else {
			validDefinition = false
		}
		if len(backupPropertySlices["excludes"]) > 0 {
			bDefinition.Excluded = backupPropertySlices["excludes"]
		} else {
			validDefinition = false
		}
		if validDefinition {
			backupsDefinitions = append(backupsDefinitions, bDefinition)
		}
	}
	return backupsDefinitions
}

func setSchedulesDefinitions(schedulesViper viper.Viper) []scheduleDefinition {
	if len(schedulesViper.AllKeys()) == 0 {
		return []scheduleDefinition{}
	}
	var scheduleDefinitions []scheduleDefinition
	for k, v := range schedulesViper.AllSettings() {
		log.Debugf("Key: %s, value: %s", k, v)
		var sDefinition scheduleDefinition
		validDefinition := true
		schedule := schedulesViper.GetStringMapString(k)
		scheduleHour := schedulesViper.GetStringMap(k)
		if schedule["id"] != "" {
			sDefinition.ID = schedule["id"]
		} else {
			validDefinition = false
		}
		if schedule["name"] != "" {
			sDefinition.Name = schedule["name"]
		} else {
			validDefinition = false
		}
		if schedule["day"] != "" {
			sDefinition.Day = schedule["day"]
		} else {
			validDefinition = false
		}
		if schedule["description"] != "" {
			sDefinition.Description = schedule["description"]
		} else {
			validDefinition = false
		}
		if scheduleHour["hour"] != "" {
			sDefinition.Hour = scheduleHour["hour"].(time.Time)
		} else {
			validDefinition = false
		}
		if validDefinition {
			scheduleDefinitions = append(scheduleDefinitions, sDefinition)
		}
	}

	return scheduleDefinitions
}

func setClientsDefinitions(clientsViper viper.Viper) []clientDefinition {
	if len(clientsViper.AllKeys()) == 0 {
		return []clientDefinition{}
	}
	var clientsDefinitions []clientDefinition
	for k, v := range clientsViper.AllSettings() {
		log.Debugf("Key: %s, value: %s", k, v)
		var cDefinition clientDefinition
		validDefinition := true
		clientDef := clientsViper.GetStringMapString(k)
		clientsPropertySlices := clientsViper.GetStringMapStringSlice(k)
		if clientDef["name"] != "" {
			cDefinition.Name = clientDef["name"]
		} else {
			validDefinition = false
		}
		if clientDef["ipaddress"] != "" {
			cDefinition.IPAddress = clientDef["ipaddress"]
		} else {
			validDefinition = false
		}
		if len(clientsPropertySlices["backupsids"]) > 0 {
			cDefinition.BackupsIDs = clientsPropertySlices["backupsids"]
		} else {
			validDefinition = false
		}
		if len(clientsPropertySlices["schedulesids"]) > 0 {
			cDefinition.SchedulesIDs = clientsPropertySlices["schedulesids"]
		} else {
			validDefinition = false
		}
		if validDefinition {
			clientsDefinitions = append(clientsDefinitions, cDefinition)
		}
	}

	return clientsDefinitions
}

//
//	backups := BackupsViper.AllSettings()
//	for k, v := range backups {
//		backupPropertySlices := BackupsViper.GetStringMapStringSlice(k)
//		backupPropertyString := BackupsViper.GetStringMapString(k)
//		MainBackupsConfig.AllBackups = append(MainBackupsConfig.AllBackups, backupDefinition{
//			ID:        backupPropertyString["id"],
//			Paths:     backupPropertySlices["paths"],
//			Excluded:  backupPropertySlices["excluded"],
//			Retention: backupPropertyString["retention"],
//		})
//
//		log.Debugf("Key: %s, value: %s", k, v)
//	}
//	return nil
//}
