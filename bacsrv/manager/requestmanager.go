package manager

import (
	log "github.com/Sirupsen/logrus"
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
func SendBackupTriggerMessage(paths []string, address string) error {
	// TODO It sould have logic like if client is integrated
	// TODO Should I send checking paths message before?
	err := protoapi.SendBackupRequest(paths, address)
	if err != nil {
		return err
	}
	return nil
}
