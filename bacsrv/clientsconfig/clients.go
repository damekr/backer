package clientsconfig

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/spf13/viper"
	"os"
)

type Client struct {
	Name     string `json:"clientName"`
	Address  string `json:"clientAddress"`
	BackupID string `json:"backupId"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}

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
		log.Error("Cannot read clients config file")
	}
}

// DoesClientIntegrated checks if client is in configuration files and is integrated
//func DoesClientIntegrated(name string) {
//	nil

//}

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

func GetClientInformations(name string) *Client {
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
