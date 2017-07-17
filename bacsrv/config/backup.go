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

type BackupTriggerMessage struct {
	ClientName   string `json:"clientName"`
	BackupConfig Backup
}

var BackupConfigInstance = viper.New()

// InitClientsConfig creating a new instance of Viper configuration with read
func InitBackupConfig(srvConfig *ServerConfig) {
	log.Info("Backups config path: ", srvConfig.BackupsConfig)
	// TODO Add checking file
	BackupConfigInstance.SetConfigName("backups")
	BackupConfigInstance.AddConfigPath(srvConfig.BackupsConfig)
	err := BackupConfigInstance.ReadInConfig()
	if err != nil {
		log.Errorf("Cannot read backup config file, an error: %s", err)
	}
	showNumberOfAddedBackupsDefinitions()
}

func GetAllBackupConfigs() []Backup {
	var backups []Backup
	backupsAll := BackupConfigInstance.AllSettings()
	for k, v := range backupsAll {
		backupPropertySlices := BackupConfigInstance.GetStringMapStringSlice(k)
		backupPropertyString := BackupConfigInstance.GetStringMapString(k)
		backups = append(backups, Backup{
			ID:        backupPropertyString["id"],
			Paths:     backupPropertySlices["paths"],
			Excluded:  backupPropertySlices["excluded"],
			Retention: backupPropertyString["retention"],
		})

		log.Debugf("Key: %s, value: %s", k, v)
	}
	return backups
}

func GetBackupConfigInformation(name string) *Backup {
	backupStrings := BackupConfigInstance.GetStringMapString(name)
	backupSlices := BackupConfigInstance.GetStringMapStringSlice(name)
	return &Backup{
		Paths:     backupSlices["paths"],
		Excluded:  backupSlices["excluded"],
		Retention: backupStrings["retention"],
	}
}

func GetBackupConfigByID(id string) (*Backup, error) {
	log.Debug("Getting backup config with ID: ", id)
	backupConfigs := GetAllBackupConfigs()
	for _, v := range backupConfigs {
		if v.ID == id {
			return &v, nil
		}
	}
	return nil, nil
}

func showNumberOfAddedBackupsDefinitions() {
	keys := BackupConfigInstance.AllSettings()
	log.Debug(keys)
	log.Info("Available backup definition(s): ", len(keys))
}
