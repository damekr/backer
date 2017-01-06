package config

import (
	"fmt"
	"github.com/spf13/viper"
)

const configName = "bacsrv"

type ServerConfig struct {
	MgmtPort  string
	DataPort  string
	LogOutput string // STDOUT, FILE, SYSLOG
	Debug     bool
}

func fillMainConfigStruct() *ServerConfig {
	return &ServerConfig{
		MgmtPort:  viper.GetString("server.MgmtPort"),
		DataPort:  viper.GetString("server.DataPort"),
		LogOutput: viper.GetString("server.LogOutput"),
		Debug:     viper.GetBool("server.Debug"),
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
	config := fillMainConfigStruct()
	return config
}

func ReadConfigFile() {
	viper.ReadInConfig()
}

func GetMgmtPort() string {
	return GetServerConfig().MgmtPort
}

func GetTransferPort() string {
	return GetServerConfig().DataPort
}
