package config

import (
	"github.com/spf13/viper"
	"log"
)

// TODO GENERAL Comment configuration reading now works without any validation. Needs to be done before real work of the application.
// At least required parameters need to be chacked

// ClientConfig specify client main configuration
type ClientConfig struct {
	MgmtPort              string
	DataPort              string
	LogOutput             string // STDOUT, FILE, SYSLOG
	ExternalName          string
	DataTransferInterface string
	Debug                 bool
	TempDir               string // Path to store temporary data
}

// ServerConfig specifies configuration of server to which client is integrated
type ServerConfig struct {
	MgmtPort string
	DataPort string
}

var (
	clntConfViper = viper.New()
	// ClntConfig is the application configuration
	ClntConfig ClientConfig
	// SrvConfig is a struct with data read from configuration file
	SrvConfig ServerConfig
)

// GetServerConfig return fullfiled server config struct with data from config file
func GetServerConfig() *ServerConfig {
	return &ServerConfig{
		MgmtPort: clntConfViper.GetString("server.MgmtPort"),
		DataPort: clntConfViper.GetString("server.DataPort"),
	}
}

// GetServerMgmtPort returns server management port on which client has to initiate connections
func GetServerMgmtPort() string {
	return clntConfViper.GetString("server.MgmtPort")
}

// GetServerDataPort returns data port of server on which client has to initiate connections
func GetServerDataPort() string {
	return clntConfViper.GetString("server.DataPort")
}

func GetTempDir() string {
	return clntConfViper.GetString("main.TempDir")
}

func GetClientMgmtPort() string {
	return clntConfViper.GetString("main.MgmtPort")
}

func GetExternalName() string {
	return ClntConfig.ExternalName
}

func GetDataTransferInterface() string {
	return ClntConfig.DataTransferInterface
}

func fillConfigData() ClientConfig {
	return ClientConfig{
		MgmtPort:              clntConfViper.GetString("main.MgmtPort"),
		DataPort:              clntConfViper.GetString("main.DataPort"),
		LogOutput:             clntConfViper.GetString("main.LogOutput"),
		Debug:                 clntConfViper.GetBool("main.Debug"),
		TempDir:               clntConfViper.GetString("main.TempDir"),
		ExternalName:          clntConfViper.GetString("main.ExternalName"),
		DataTransferInterface: clntConfViper.GetString("main.DataTransferInterface"),
	}
}

func (c *ClientConfig) ShowConfig() {
	log.Printf("Config Struct: %#v\n", c)
}

func setConfigPath(path string) {
	// Viper can cooperate with Cobra arg parser consider reading config file path from arg
	clntConfViper.SetConfigFile(path)
}

// ReadInConfigFile reads file into memory at the beginning of application start
func ReadInConfigFile(path string) {
	setConfigPath(path)
	clntConfViper.ReadInConfig()
	ClntConfig = fillConfigData()

}
