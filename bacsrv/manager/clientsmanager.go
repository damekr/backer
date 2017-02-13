package manager

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/clientsconfig"
	"github.com/damekr/backer/bacsrv/protoapi"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func IntegrateClient(name string, address string, backupID string) error {
	//	client := &clientconfig.Client{
	//		Name:     name,
	//		Address:  address,
	//		BackupID: backupID,
	//	}
	_, err := protoapi.SayHelloToClient(address)
	if err != nil {
		log.Errorf("Client %s is not available", name)
		return err
	}
	//TODO Database create trigger  to add a new client shall be here
	//TODO Creating client backet should be here
	return nil

}

// GetAllIntegratedClients simply fetching clients from clients configuration file, at least now
func GetAllIntegratedClients() []clientsconfig.Client {
	return clientsconfig.GetAllClients()
}
