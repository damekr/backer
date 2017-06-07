package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const configName = "bacsrv"

var (
	server *ServerConfig
)

type ServerConfig struct {
	MgmtPort              string
	DataPort              string
	RestAPIPort           string
	LogOutput             string // STDOUT, FILE, SYSLOG
	Debug                 bool
	RepositoryConfig      string
	ClientsConfig         string
	ExternalName          string
	DataTransferInterface string
	DBLocation            string
}

func fillMainConfigStruct() *ServerConfig {
	return &ServerConfig{
		MgmtPort:              viper.GetString("server.MgmtPort"),
		DataPort:              viper.GetString("server.DataPort"),
		RestAPIPort:           viper.GetString("server.RestApiPort"),
		ExternalName:          viper.GetString("server.ExternalName"),
		DataTransferInterface: viper.GetString("server.DataTransferInterface"),
		LogOutput:             viper.GetString("server.LogOutput"),
		Debug:                 viper.GetBool("server.Debug"),
		ClientsConfig:         viper.GetString("clients.ConfigFile"),
		DBLocation:            viper.GetString("server.DBLocation"),
	}
}

func (c *ServerConfig) ShowConfig() {
	fmt.Printf("Config Struct: %#v\n", c)
}

func SetConfigPath(path string) {
	// Viper can cooperate with Cobra arg parser consider reading config file path from
	viper.SetConfigFile(path)
	//viper.SetConfigName(configName)
	//viper.AddConfigPath(path)
}

func GetServerConfig() *ServerConfig {
	ReadConfigFile()
	server = fillMainConfigStruct()
	return server
}

func ReadConfigFile() {
	viper.ReadInConfig()
}

func GetClientconfigPath() string {
	return GetServerConfig().ClientsConfig
}

func GetMgmtPort() string {
	return GetServerConfig().MgmtPort
}

func GetTransferPort() string {
	return GetServerConfig().DataPort
}

func GetExternalName() string {
	return GetServerConfig().ExternalName
}

func GetDataTransferInterface() string {
	return GetServerConfig().DataTransferInterface
}

func GetDBLocation() string {
	return server.DBLocation
}
