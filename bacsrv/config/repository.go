package config

import (
	"github.com/spf13/viper"
)



type RepositoryConfig struct {
	Localization string
}

func GetRepositoryConfig() *RepositoryConfig {
	return &RepositoryConfig{
		Localization: viper.GetString("Repository.Localization"),
	}
}

func GetMainRepositoryLocation() string {
	// TODO This below does not look well, needs to be changed
	repoLocalization := GetRepositoryConfig().Localization
	return repoLocalization
}
