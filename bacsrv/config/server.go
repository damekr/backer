package config

import (
	"github.com/spf13/viper"
	"github.com/pkg/errors"
)

var (
	EmptyMainConfig = errors.New("MainConfig: Configuration file is empty")
	MainConfig      = ServerConfig{}
)


type ServerConfig struct {
	MgmtPort              string
	DataPort              string
	RestAPIPort           string
	LogOutput             string // STDOUT, FILE, SYSLOG
	Debug                 bool
	RepositoryConfig      string
	ClientsConfigFilePath string
	BackupsConfigFilePath string
	SchedulesConfigFilePath string
	ExternalName          string
	DataTransferInterface string
	DBLocation            string
}


func ReadInServerConfig(path string) error {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	if len(viper.AllKeys()) == 0 {
		return EmptyMainConfig
	}
	MainConfig = ServerConfig{
		MgmtPort:              viper.GetString("server.MgmtPort"),
		DataPort:              viper.GetString("server.DataPort"),
		RestAPIPort:           viper.GetString("server.RestApiPort"),
		ExternalName:          viper.GetString("server.ExternalName"),
		DataTransferInterface: viper.GetString("server.DataTransferInterface"),
		LogOutput:             viper.GetString("server.LogOutput"),
		Debug:                 viper.GetBool("server.Debug"),
		ClientsConfigFilePath: viper.GetString("clients.ConfigFile"),
		BackupsConfigFilePath: viper.GetString("backups.ConfigFile"),
		SchedulesConfigFilePath: viper.GetString("schedules.ConfigFile"),
		DBLocation:            viper.GetString("server.DBLocation"),
	}
	return nil
}


