package manager

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/protoapi"
)

// HelloMessageManager is responsible for proxing restapi reqests to clients
func HelloMessageManager(address string) (string, error) {
	clntHostname, err := protoapi.SayHelloToClient(address)
	if err != nil {
		log.Errorf("Given client on address %s is not available", address)
		return "", err
	}
	return clntHostname, nil

}
