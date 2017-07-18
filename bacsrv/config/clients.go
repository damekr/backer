package config

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type ClientConfig struct {
	Name     string `json:"clientName"`
	Address  string `json:"clientAddress"`
	BackupID string `json:"backupId"`
	// TODO BackupID must be an integer
	Platform string `json:"platform"`
	CID      string `json:"cid"`
}

// TODO - GENERAL - It needs to be refactored, remove repeated functions checking executing. Error creating instead of returning nil in case of error.

// ClientsConfigInstance represents a new instance of viper config library for configuration
var ClientsConfigInstance = viper.New()

// InitClientsConfig creating a new instance of Viper configuration with read
func InitClientsConfig(srvConfig *ServerConfig) {
	log.Info("ClientConfig config path: ", srvConfig.ClientsConfig)
	// TODO Add checking file
	ClientsConfigInstance.SetConfigName("clients")
	ClientsConfigInstance.AddConfigPath(srvConfig.ClientsConfig)
	err := ClientsConfigInstance.ReadInConfig()
	if err != nil {

		log.Errorf("Cannot read clients config file, an error: %s", err)
	}
}

func GetAllClients() []ClientConfig {
	var Clients []ClientConfig
	clients := ClientsConfigInstance.AllSettings()
	for k, v := range clients {
		clientProperty := ClientsConfigInstance.GetStringMapString(k)
		Clients = append(Clients, ClientConfig{
			Name:     k,
			Address:  clientProperty["ip"],
			BackupID: clientProperty["backupid"],
		})

		log.Debugf("Key: %s, value: %s", k, v)
	}
	return Clients
}

func GetClientInformation(name string) *ClientConfig {
	if !DoesClientExist(name) {
		return &ClientConfig{}
	}
	client := ClientsConfigInstance.GetStringMapString(name)
	return &ClientConfig{
		Name:     name,
		Address:  client["ip"],
		BackupID: client["backupid"],
	}

}

func DoesClientExistWithIP(address string) bool {
	clients := GetAllClients()
	for _, v := range clients {
		if v.Address == address {
			log.Debugf("ClientConfig: %s with address: %s exists", v.Name, address)
			return true
		}
	}
	return false
}

// GetClientIP returns an ip address from name
func GetClientIP(name string) (string, error) {
	client := GetClientInformation(name)
	if client.Address == "" {
		return "", errors.New("ClientConfig does not exist")
	}
	return client.Address, nil
}

// DoesClientExist checks if client has been added to clients configuration file
func DoesClientExist(name string) bool {
	client := ClientsConfigInstance.GetStringMapString(name)
	if client == nil {
		log.Warningf("Requested client %s does not exist", name)
		return false
	}
	log.Debugf("Requested client %s exists", name)
	return true
}

func PrintValues() {
	fmt.Println(viper.AllSettings())
}
