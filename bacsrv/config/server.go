package config

import (
	"fmt"
	"github.com/spf13/viper"
)

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

func (c *ServerConfig) showConfig() {
	fmt.Printf("Config Struct: %#v\n", c)
}

func setConfigPath() {
	// Viper can cooperate with Cobra arg parser consider reading config file path from arg
	viper.SetConfigName("bacsrv")
	viper.AddConfigPath("/home/damekr/dev/go/src/github.com/backer/config")
}

func GetServerSettings() *ServerConfig {
	ReadConfigFile()
	config := fillMainConfigStruct()
	return config
}

func ReadConfigFile() {
	setConfigPath()
	viper.ReadInConfig()
}

func GetMgmtPort() string {
	return GetServerSettings().MgmtPort
}

func GetTransferPort() string {
	return GetServerSettings().DataPort
}
