package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	MainConfig = ClientConfig{}
	log        = logrus.WithFields(logrus.Fields{"prefix": "config"})
)

// ClientConfig specify client main configuration
type ClientConfig struct {
	ManagementPort        string
	DataPort              string
	LogOutput             string // STDOUT, FILE, SYSLOG
	ExternalName          string
	DataTransferInterface string
	Debug                 bool
	TempDir               string // AbsolutePath to store temporary data
	CID                   string
	ServerDataPort        string
	ServerExternalName    string
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
		ManagementPort:        viper.GetString("main.ManagementPort"),
		DataPort:              viper.GetString("main.DataPort"),
		LogOutput:             viper.GetString("main.LogOutput"),
		Debug:                 viper.GetBool("main.Debug"),
		TempDir:               viper.GetString("main.TempDir"),
		ExternalName:          viper.GetString("main.ExternalName"),
		DataTransferInterface: viper.GetString("main.DataTransferInterface"),
		ServerDataPort:        viper.GetString("server.DataPort"),
		ServerExternalName:    viper.GetString("server.ExternalName"),
	}
	return nil
}

func (c *ClientConfig) ShowConfig() {
	log.Printf("Config Struct: %#v\n", c)
}
