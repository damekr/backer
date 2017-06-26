package operations

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"

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

// IntegrateClient performs client integration with all operatinos
func IntegrateClient(client *config.Client) error {
	log.Infof("Starting %s integration...", client.Name)
	clntHostname, err := SendHelloMessageToClient(client.Address)
	if err != nil {
		log.Errorf("Client %s with address: %s does not respond", client.Name, client.Address)
		return err
	}
	log.Debugf("Got hostname: %s from client side, performing integration", clntHostname)
	remoteInformations, err := outprotoapi.SendIntegrationRequest(client)
	if err != nil {
		log.Error("Cannot get information from remote host: ", client.Address)
	}
	log.Debugf("Got information %#v about client", remoteInformations)
	return nil
}

// StartBackup start backup on client with given configuration
func StartBackup(backupConfig *config.Backup, clntAddr string) error {
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

// GetAllIntegratedClients simply fetching clients from clients configuration file, at least now and shows them
func GetAllIntegratedClients() []config.Client {
	return config.GetAllClients()
}
