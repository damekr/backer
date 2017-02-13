package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	// TODO It cannot be in init because of overwriting
	ReadConfigFile()
}

type RepositoryConfig struct {
	Localization string
}

func GetRepositoryConfig() *RepositoryConfig {
	return &RepositoryConfig{
		Localization: viper.GetString("Repository.Localization"),
	}
}

func GetRepositoryLocalization() string {
	return GetRepositoryConfig().Localization
}
