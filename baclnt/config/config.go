package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type ClientConfig struct {
	MgmtPort  string
	LogOutput string // STDOUT, FILE, SYSLOG
	Debug     bool
}

var clntConfig = viper.New()

const (
	MgmtPort = "9090"
)

func fillConfigData() *ClientConfig {
	return &ClientConfig{
		MgmtPort:  MgmtPort,
		LogOutput: clntConfig.GetString("main.LogOutput"),
		Debug:     clntConfig.GetBool("main.Debug"),
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
	return config
}
