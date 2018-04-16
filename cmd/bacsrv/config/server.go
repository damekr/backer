package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

var (
	EmptyMainConfig = errors.New("MainConfig: Configuration file is empty")
	MainConfig      = ServerConfig{}
)

type Storage struct {
	Type     string
	Location string
}

type ServerConfig struct {
	MgmtPort              string
	DataPort              string
	RestAPIPort           string
	ClntMgmtPort          string
	LogOutput             string // STDOUT, FILE, SYSLOG
	Debug                 bool
	Storage               Storage
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
	storage := Storage{
		Type:     viper.GetString("storage.type"),
		Location: viper.GetString("storage.location"),
	}
	MainConfig = ServerConfig{
		MgmtPort:              viper.GetString("server.MgmtPort"),
		DataPort:              viper.GetString("server.DataPort"),
		RestAPIPort:           viper.GetString("server.RestApiPort"),
		ClntMgmtPort:          viper.GetString("server.ClntMgmtPort"),
		ExternalName:          viper.GetString("server.ExternalName"),
		DataTransferInterface: viper.GetString("server.DataTransferInterface"),
		LogOutput:             viper.GetString("server.LogOutput"),
		Debug:                 viper.GetBool("server.Debug"),
		Storage:               storage,
		DBLocation:            viper.GetString("server.DBLocation"),
	}
	return nil
}
