package manager

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/backupconfig"
	"github.com/damekr/backer/bacsrv/clientsconfig"

	"github.com/damekr/backer/bacsrv/outprotoapi"
)

// SendHelloMessageToClient is responsible for proxing restapi reqests to clients
func SendHelloMessageToClient(clntAddress string) (string, error) {
	clntHostname, err := outprotoapi.SayHelloToClient(clntAddress)
	if err != nil {
		log.Errorf("Given client on address %s is not available", clntAddress)
		return "", err
	}
	return clntHostname, nil

}

func preBackupChecks(paths []string, clntAddr string) ([]string, error) {
	log.Debug("Starting executing prebackup checks...")
	log.Debug("Checking if client repsponds")
	hostname, err := SendHelloMessageToClient(clntAddr)
	if err != nil {
		log.Errorf("Client %s does not responds", clntAddr)
		return nil, err
	}
	log.Debug("Client sent it's own hostname, and it is: ", hostname)
	checkedPaths, err := outprotoapi.CheckPaths(clntAddr, paths)
	if err != nil {
		log.Error("An error ocurred during checking paths, error: ", err.Error())
		return nil, err
	}
	return checkedPaths, nil
}

func StartBackup(backupConfig *backupconfig.Backup, clntAddr string) error {
	log.Info("Starting backup of client: ", clntAddr)
	validatedPaths, err := preBackupChecks(backupConfig.Paths, clntAddr)
	if err != nil {
		log.Error("Cannot validate paths on client side")
		return err
	}
	log.Debugf("Got validated paths from client: %s starting backup...", validatedPaths)
	err = outprotoapi.SendBackupRequest(validatedPaths, clntAddr)
	if err != nil {
		log.Error("Triggering backup failed!")
		return err
	}
	log.Info("Backup has been triggered properly!")
	return nil
}

// IntegrateClient creates client config and should add it to configuration file
func IntegrateClient(name string, address string, backupID string) error {
	//	client := &clientconfig.Client{
	//		Name:     name,
	//		Address:  address,
	//		BackupID: backupID,
	//	}
	_, err := outprotoapi.SayHelloToClient(address)
	if err != nil {
		log.Errorf("Client %s is not available", name)
		return err
	}
	//TODO Database create trigger  to add a new client shall be here
	//TODO Creating client backet should be here
	return nil

}

// GetAllIntegratedClients simply fetching clients from clients configuration file, at least now and shows them
func GetAllIntegratedClients() []clientsconfig.Client {
	return clientsconfig.GetAllClients()
}
