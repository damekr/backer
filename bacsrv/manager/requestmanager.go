package manager

import (
	"errors"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/backupconfig"
	"github.com/damekr/backer/bacsrv/clientsconfig"
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
