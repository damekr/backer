package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

// Backup specifies a backup
type Backup struct {
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

func showNumberOfAddedBackupsDefinitions() {
	keys := BackupConfigInstance.AllSettings()
	log.Debug(keys)
	log.Info("Available backup definition(s): ", len(keys))
}
