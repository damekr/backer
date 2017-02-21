package config

import (
	"github.com/spf13/viper"
)

func init() {
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
