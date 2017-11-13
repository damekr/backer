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
	log                = logger.New().WithFields(logger.Fields{"MODULE": "CLIENTS-CONFIG"})
)

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
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Hour        string `json:"startHour"`
	Day         string `json:"startDay"`
}

// TODO Must be places in communication protocol
type BackupTriggerMessage struct {
	ClientName   string `json:"clientName"`
	BackupConfig backupDefinition
}

func ReadInClientsConfig(clientsConfigPath, backupsConfigPath, scheduleConfigPath string) error {
	backupsViper, err := readInBackupsConfig(backupsConfigPath)
	if err != nil {
		log.Error("Error occured, err: ", err)
		return err
	}
	backupsDefinitions := setBackupsDefinitions(backupsViper)

	schedulesViper, err := readInSchedulesConfig(scheduleConfigPath)
	if err != nil {
		log.Error("Error occured, err: ", err)

		return err
	}
	schedulesDefinitions := setSchedulesDefinitions(schedulesViper)

	clientsViper, err := readInClientsConfig(clientsConfigPath)
	if err != nil {
		log.Error("Error occured, err: ", err)

		return err
	}
	clientsDefinitions := setClientsDefinitions(clientsViper)
	for _, client := range clientsDefinitions {
		matchedClient := matchClientDefinition(client, backupsDefinitions, schedulesDefinitions)

		AllClients = append(AllClients, matchedClient)

	}

	return nil
}

func matchClientDefinition(clntDefinition clientDefinition, backupsDefinitions []backupDefinition, schedulesDefinitions []scheduleDefinition) clientDefinition {
	for bk, bv := range backupsDefinitions {
		for _, bid := range clntDefinition.BackupsIDs {
			if bid == bv.ID {
				log.Debug("Found matched backup ID: ", bid)
				clntDefinition.BackupDefinitions = append(clntDefinition.BackupDefinitions, backupsDefinitions[bk])
			}
		}
	}
	for sk, sv := range schedulesDefinitions {
		for _, sid := range clntDefinition.SchedulesIDs {
			if sid == sv.ID {
				log.Debug("Found matched schedule ID: ", sid)
				clntDefinition.SchedulesDefinitions = append(clntDefinition.SchedulesDefinitions, schedulesDefinitions[sk])
			}
		}
	}
	return clntDefinition
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

func checkValidConfigKey(keyMap map[string]string, key, ctx string) (string, bool) {
	log.Debug("Checking map of: ", ctx)
	if keyMap[key] != "" {
		log.Debug("Found value in key: ", key)
		return keyMap[key], true
	} else {
		log.Debug("Did not find value in key: ", key)
		return "", false
	}
}

func checkValidConfigKeySlice(keyMap map[string][]string, key, ctx string) ([]string, bool) {
	log.Debug("Checking map with slices of: ", ctx)
	if len(keyMap[key]) > 0 {
		log.Debug("Found values in key: ", key)
		return keyMap[key], true
	} else {
		log.Debug("Did not found at least 1 value in slice in key: ", key)
		return nil, false
	}
}

func parseHour(hour string) {
	// TODO Must be clarified to const format of time
	form := "2006-Jan-02 07:04 PM"
	dateNow := time.Now().Format("2006-Jan-02")
	log.Println("DATE: ", dateNow)

	log.Println(time.Parse(form, dateNow+" "+hour))
}

func parseDayHour(day, hour string) {

}

func setBackupsDefinitions(backupsViper *viper.Viper) []backupDefinition {
	if len(backupsViper.AllKeys()) == 0 {
		log.Error("Backup config does not contain any backup definitions")
		return []backupDefinition{}
	}
	var backupsDefinitions []backupDefinition
	backupCtx := "backup"
	for k, v := range backupsViper.AllSettings() {
		log.Debugf("Key: %s, value: %s", k, v)
		var bDefinition backupDefinition
		validDefinition := true
		backupPropertySlices := backupsViper.GetStringMapStringSlice(k)
		backup := backupsViper.GetStringMapString(k)
		bDefinition.Name = k
		bDefinition.ID, validDefinition = checkValidConfigKey(backup, "id", backupCtx)
		bDefinition.Retention, validDefinition = checkValidConfigKey(backup, "retention", backupCtx)
		bDefinition.Description, validDefinition = checkValidConfigKey(backup, "description", backupCtx)

		bDefinition.Paths, validDefinition = checkValidConfigKeySlice(backupPropertySlices, "paths", backupCtx)
		bDefinition.Excluded, validDefinition = checkValidConfigKeySlice(backupPropertySlices, "excluded", backupCtx)

		if validDefinition {
			backupsDefinitions = append(backupsDefinitions, bDefinition)
		} else {
			log.Warning("Found invalid backup definition")
		}
	}
	return backupsDefinitions
}

func setSchedulesDefinitions(schedulesViper *viper.Viper) []scheduleDefinition {
	if len(schedulesViper.AllKeys()) == 0 {
		log.Error("Schedule config does not contain any valid schedule definitions")
		return []scheduleDefinition{}
	}
	var scheduleDefinitions []scheduleDefinition
	scheduleCtx := "schedule"
	parseHour("07:00pm")
	for k, v := range schedulesViper.AllSettings() {
		log.Debugf("Key: %s, value: %s", k, v)
		var sDefinition scheduleDefinition
		validDefinition := true
		schedule := schedulesViper.GetStringMapString(k)
		sDefinition.Name = k
		sDefinition.ID, validDefinition = checkValidConfigKey(schedule, "id", scheduleCtx)
		sDefinition.Day, validDefinition = checkValidConfigKey(schedule, "day", scheduleCtx)
		sDefinition.Description, validDefinition = checkValidConfigKey(schedule, "description", scheduleCtx)
		sDefinition.Hour, validDefinition = checkValidConfigKey(schedule, "hour", scheduleCtx)

		if validDefinition {
			scheduleDefinitions = append(scheduleDefinitions, sDefinition)
		} else {
			log.Error("Found invalid schedule definition")
		}
	}
	return scheduleDefinitions
}

func setClientsDefinitions(clientsViper *viper.Viper) []clientDefinition {
	if len(clientsViper.AllKeys()) == 0 {
		log.Error("Client config does not contain any valid client definitions")
		return []clientDefinition{}
	}
	var clientsDefinitions []clientDefinition
	clientCtx := "client"
	for k, v := range clientsViper.AllSettings() {
		log.Debugf("Key: %s, value: %s", k, v)
		var cDefinition clientDefinition
		validDefinition := true
		clientDef := clientsViper.GetStringMapString(k)
		clientsPropertySlices := clientsViper.GetStringMapStringSlice(k)

		cDefinition.Name = k
		cDefinition.IPAddress, validDefinition = checkValidConfigKey(clientDef, "ipaddress", clientCtx)
		cDefinition.BackupsIDs, validDefinition = checkValidConfigKeySlice(clientsPropertySlices, "backupsids", clientCtx)
		cDefinition.SchedulesIDs, validDefinition = checkValidConfigKeySlice(clientsPropertySlices, "schedulesids", clientCtx)
		if validDefinition {
			clientsDefinitions = append(clientsDefinitions, cDefinition)
		} else {
			log.Error("Found invalid client definition")
		}
	}
	return clientsDefinitions
}
