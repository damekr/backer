package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ClientConfig struct {
	MgmtPort  string
	LogOutput string // STDOUT, FILE, SYSLOG
	Debug     bool
	TempDir   string // Path to store temporary data
}

var clntConfig = viper.New()

const (
	MgmtPort = "9090"
)

func GetTempDir() string {
	return clntConfig.GetString("main.TempDir")
}

func fillConfigData() *ClientConfig {
	return &ClientConfig{
		MgmtPort:  MgmtPort,
		LogOutput: clntConfig.GetString("main.LogOutput"),
		Debug:     clntConfig.GetBool("main.Debug"),
		TempDir:   clntConfig.GetString("main.TempDir"),
	}
}

func (c *ClientConfig) ShowConfig() {
	fmt.Printf("Config Struct: %#v\n", c)
}

func setConfigPath(path string) {
	// Viper can cooperate with Cobra arg parser consider reading config file path from arg
	clntConfig.SetConfigFile(path)
}

func ReadConfigFile(path string) *ClientConfig {
	setConfigPath(path)
	clntConfig.ReadInConfig()
	config := fillConfigData()
	fmt.Printf("Config temp path: %s", config.TempDir)
	return config
}
