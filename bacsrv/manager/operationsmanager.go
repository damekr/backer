package manager

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/backupconfig"
	"github.com/damekr/backer/bacsrv/clientsconfig"
	"github.com/damekr/backer/bacsrv/operationshandler"
	"github.com/damekr/backer/bacsrv/protoapi"
)

// SendHelloMessage is responsible for proxing restapi reqests to clients
func SendHelloMessage(address string) (string, error) {
	clntHostname, err := protoapi.SayHelloToClient(address)
	if err != nil {
		log.Errorf("Given client on address %s is not available", address)
		return "", err
	}
	return clntHostname, nil

}

// IntegrateClient creates client config and should add it to configuration file
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

// GetAllIntegratedClients simply fetching clients from clients configuration file, at least now and shows them
func GetAllIntegratedClients() []clientsconfig.Client {
	return clientsconfig.GetAllClients()
}

// SendBackupTriggerMessage sending a message to client with specific address and does not wait for status
func SendBackupTriggerMessage(backupMessage *backupconfig.BackupTriggerMessage) (string, error) {
	// TODO It sould have logic like if client is integrated
	// TODO Should I send checking paths message before?
	client := clientsconfig.GetClientInformation(backupMessage.ClientName)
	if client.Name == "" {
		return "", errors.New("Client does not exist")
	}
	err := protoapi.SendBackupRequest(backupMessage.BackupConfig.Paths, client.Address)
	if err != nil {
		return "", err
	}
	return "Clients is not added into clients config file", nil

}

func triggerRestore(restore *operationshandler.Restore, clientName string) error {

	return nil
}

func SendRestoreTriggerMessage(restoreMessage *operationshandler.RestoreTriggerMessage) error {
	// This function will get ready structure with appriopriate saveset from repository
	// Consider make saveset as struct
	// RPC --> checking space, etc
	restoreMessage.RestoreConfig.SavesetSize = 120
	log.Printf("Restore Struct: ", restoreMessage)
	err := protoapi.SendRestoreRequest(restoreMessage.RestoreConfig.SavesetSize, true, restoreMessage.ClientName)
	if err != nil {
		log.Error("Send restore trigger message failed")
	}
	// DataTransfer
	return nil
}
