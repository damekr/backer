package config

import (
	"github.com/spf13/viper"
	"log"
)

type ClientConfig struct {
	MgmtPort     string
	LogOutput    string // STDOUT, FILE, SYSLOG
	ExternalName string
	Debug        bool
	TempDir      string // Path to store temporary data
}

var (
	clntConfViper = viper.New()
	ClntConfig    ClientConfig
)

const (
	MgmtPort = "9090"
)

func GetTempDir() string {
	return clntConfViper.GetString("main.TempDir")
}

func GetExternalName() string {
	return ClntConfig.ExternalName
}

func fillConfigData() *ClientConfig {
	return &ClientConfig{
		MgmtPort:     MgmtPort,
		LogOutput:    clntConfViper.GetString("main.LogOutput"),
		Debug:        clntConfViper.GetBool("main.Debug"),
		TempDir:      clntConfViper.GetString("main.TempDir"),
		ExternalName: clntConfViper.GetString("main.ExternalName"),
	}
}

func (c *ClientConfig) ShowConfig() {
	log.Printf("Config Struct: %#v\n", c)
}

func setConfigPath(path string) {
	// Viper can cooperate with Cobra arg parser consider reading config file path from arg
	clntConfViper.SetConfigFile(path)
}

func ReadConfigFile(path string) *ClientConfig {
	setConfigPath(path)
	clntConfViper.ReadInConfig()
	clntConfig := fillConfigData()
	ClntConfig = *clntConfig
	log.Printf("Config temp path: %s", clntConfig.TempDir)
	return clntConfig
}
