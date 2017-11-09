package config

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type Client struct {
	Name     string `json:"clientName"`
	Address  string `json:"clientAddress"`
	BackupsIDs []string `json:"backupsIDs"`
	BackupDefinitions []Backup
	SchedulesIDs []string `json:"schedulesIds"`
	SchedulesDefinitions []Schedule
	Platform string `json:"platform"`
	CID      string `json:"cid"`
}

type ClientsConfigs struct {
	AllClients []Client
}


// ClientsViper represents a new instance of viper config library for configuration
var (
	ClientsViper       = viper.New()
	MainClientsConfig = ClientsConfigs{}
	EmptyClientsConfig = errors.New("ClientsConfigFilePath: Empty clients config file")
)

func ReadInClientsConfig(clientsConfigPath, backupsConfigPath, scheduleConfigPath string) error {
	log.Debugln("Clients config file path: ", clientsConfigPath)
	ClientsViper.AddConfigPath(clientsConfigPath)
	ClientsViper.SetConfigName("clients")
	err := ClientsViper.ReadInConfig()
	if err != nil {
		return err
	}
	err = readInBackupsConfig(backupsConfigPath)
	if err != nil {
		return err
	}
	err = readInSchedulesConfig(scheduleConfigPath)
	if err != nil {
		return err
	}
	clients := ClientsViper.AllSettings()
	for k, v := range clients {
		clientProperty := ClientsViper.GetStringMapString(k)
		clientIDsSlices := ClientsViper.GetStringMapStringSlice(k)
		MainClientsConfig.AllClients = append(MainClientsConfig.AllClients, Client{
			Name:     k,
			Address:  clientProperty["ip"],
			BackupsIDs: clientIDsSlices["backupsids"],
			SchedulesIDs: clientIDsSlices["schedulesids"],
		})

		log.Debugf("Key: %s, value: %s", k, v)
	}
	if len(MainClientsConfig.AllClients) == 0 {
		return EmptyClientsConfig
	}
	return nil
}

func matchBackupsConfigs(configPath, backupId string)([]Backup, error){

	return nil, nil

}

func matchSchedulesConfigs(configPath, scheduleId string)([]Schedule, error){

	return nil, nil
}


//
//func GetClientInformation(name string) *Client {
//	if !DoesClientExist(name) {
//		return &Client{}
//	}
//	client := ClientsViper.GetStringMapString(name)
//	return &Client{
//		Name:     name,
//		Address:  client["ip"],
//		BackupID: client["backupid"],
//	}
//
//}
//
//func DoesClientExistWithIP(address string) bool {
//	clients := GetAllClients()
//	for _, v := range clients {
//		if v.Address == address {
//			log.Debugf("Client: %s with address: %s exists", v.Name, address)
//			return true
//		}
//	}
//	return false
//}
//
//// GetClientIP returns an ip address from name
//func GetClientIP(name string) (string, error) {
//	client := GetClientInformation(name)
//	if client.Address == "" {
//		return "", errors.New("Client does not exist")
//	}
//	return client.Address, nil
//}
//
//// DoesClientExist checks if client has been added to clients configuration file
//func DoesClientExist(name string) bool {
//	client := ClientsViper.GetStringMapString(name)
//	if client == nil {
//		log.Warningf("Requested client %s does not exist", name)
//		return false
//	}
//	log.Debugf("Requested client %s exists", name)
//	return true
//}
//
//func PrintValues() {
//	fmt.Println(viper.AllSettings())
//}
