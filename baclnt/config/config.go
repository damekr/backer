package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	MainConfig = ClientConfig{}
)

// ClientConfig specify client main configuration
type ClientConfig struct {
	MgmtPort              string
	DataPort              string
	LogOutput             string // STDOUT, FILE, SYSLOG
	ExternalName          string
	DataTransferInterface string
	Debug                 bool
	TempDir               string // Path to store temporary data
	CID                   string
}

func ReadInConfig(path string) error {
	viper.SetConfigFile(path)
	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	if len(viper.AllKeys()) == 0 {
		return nil
	}
	MainConfig = ClientConfig{
		MgmtPort:              viper.GetString("main.MgmtPort"),
		DataPort:              viper.GetString("main.DataPort"),
		LogOutput:             viper.GetString("main.LogOutput"),
		Debug:                 viper.GetBool("main.Debug"),
		TempDir:               viper.GetString("main.TempDir"),
		ExternalName:          viper.GetString("main.ExternalName"),
		DataTransferInterface: viper.GetString("main.DataTransferInterface"),
	}
	return nil
}

func (c *ClientConfig) ShowConfig() {
	log.Printf("Config Struct: %#v\n", c)
}
