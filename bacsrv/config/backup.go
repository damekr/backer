package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)





// Backup specifies a backup
type Backup struct {
	ID        string   `json:"id"`
	Paths     []string `json:"paths"`
	Excluded  []string `json:"excludedPaths"`
	Retention string   `json:"retentionTime"`
}

type BackupsConfigs struct {
	AllBackups []Backup
}

type BackupTriggerMessage struct {
	ClientName   string `json:"clientName"`
	BackupConfig Backup
}


var (
	BackupsViper = viper.New()
	MainBackupsConfig = BackupsConfigs{}
)


func readInBackupsConfig(path string) error {
	log.Debugln("Backups config file path: ", path )
	BackupsViper.AddConfigPath(path)
	BackupsViper.SetConfigName("backups")
	err := BackupsViper.ReadInConfig()
	if err != nil {
		return err
	}
	backups := BackupsViper.AllSettings()
	for k, v := range backups {
		backupPropertySlices := BackupsViper.GetStringMapStringSlice(k)
		backupPropertyString := BackupsViper.GetStringMapString(k)
		MainBackupsConfig.AllBackups = append(MainBackupsConfig.AllBackups, Backup{
			ID:        backupPropertyString["id"],
			Paths:     backupPropertySlices["paths"],
			Excluded:  backupPropertySlices["excluded"],
			Retention: backupPropertyString["retention"],
		})

		log.Debugf("Key: %s, value: %s", k, v)
	}
	return nil
}

//
//// InitClientsConfig creating a new instance of Viper configuration with read
//func InitBackupConfig(srvConfig *ServerConfig) {
//	log.Info("Backups config path: ", srvConfig.BackupsConfigFilePath)
//	// TODO Add checking file
//	BackupConfigInstance.SetConfigName("backups")
//	BackupConfigInstance.AddConfigPath(srvConfig.BackupsConfigFilePath)
//	err := BackupConfigInstance.ReadInConfig()
//	if err != nil {
//		log.Errorf("Cannot read backup config file, an error: %s", err)
//	}
//	showNumberOfAddedBackupsDefinitions()
//}
//
//func GetAllBackupConfigs() []Backup {
//	var backups []Backup
//
//	return backups
//}
//
//func GetBackupConfigInformation(name string) *Backup {
//	backupStrings := BackupConfigInstance.GetStringMapString(name)
//	backupSlices := BackupConfigInstance.GetStringMapStringSlice(name)
//	return &Backup{
//		Paths:     backupSlices["paths"],
//		Excluded:  backupSlices["excluded"],
//		Retention: backupStrings["retention"],
//	}
//}
//
//func GetBackupConfigByID(id string) (*Backup, error) {
//	log.Debug("Getting backup config with ID: ", id)
//	backupConfigs := GetAllBackupConfigs()
//	for _, v := range backupConfigs {
//		if v.ID == id {
//			return &v, nil
//		}
//	}
//	return nil, nil
//}
//
//func showNumberOfAddedBackupsDefinitions() {
//	keys := BackupConfigInstance.AllSettings()
//	log.Debug(keys)
//	log.Info("Available backup definition(s): ", len(keys))
//}
