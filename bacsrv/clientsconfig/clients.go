package clientsconfig

import (
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/spf13/viper"
)

type Client struct {
	Name     string `json:"clientName"`
	Address  string `json:"clientAddress"`
	BackupID string `json:"backupId"`
	Platform string `json:"platform"`
	CID      string `json:"cid"`
}

// TODO - GENERAL - It needs to be refactored, remove repeated functions checking executing. Error creating instead of returning nil in case of error.

// ClientsConfigInstance represents a new instance of viper config library for configuration
var ClientsConfigInstance = viper.New()

// InitClientsConfig creating a new instance of Viper configuration with read
func InitClientsConfig(srvConfig *config.ServerConfig) {
	log.Info("Client config path: ", srvConfig.ClientsConfig)
	// TODO Add checking file
	ClientsConfigInstance.SetConfigName("config")
	ClientsConfigInstance.AddConfigPath(srvConfig.ClientsConfig)
	err := ClientsConfigInstance.ReadInConfig()
	if err != nil {

		log.Errorf("Cannot read clients config file, an error: %s", err)
	}
}

func GetAllClients() []Client {
	var Clients []Client
	clients := ClientsConfigInstance.AllSettings()
	for k, v := range clients {
		clientProperty := ClientsConfigInstance.GetStringMapString(k)
		Clients = append(Clients, Client{
			Name:     k,
			Address:  clientProperty["ip"],
			BackupID: clientProperty["backupid"],
		})

		log.Debugf("Key: %s, value: %s", k, v)
	}
	return Clients
}

func GetClientInformation(name string) *Client {
	if !DoesClientExist(name) {
		return &Client{}
	}
	client := ClientsConfigInstance.GetStringMapString(name)
	return &Client{
		Name:     name,
		Address:  client["ip"],
		BackupID: client["backupid"],
	}

}

func DoesClientExistWithIP(address string) bool {
	clients := GetAllClients()
	for _, v := range clients {
		if v.Address == address {
			log.Debugf("Client: %s with address: %s exists", v.Name, address)
			return true
		}
	}
	return false
}

// GetClientIP returns an ip address from name
func GetClientIP(name string) (string, error) {
	client := GetClientInformation(name)
	if client.Address == "" {
		return "", errors.New("Client does not exist")
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
