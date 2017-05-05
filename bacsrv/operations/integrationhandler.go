package operations

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/clientsconfig"
)

func init() {
	log.Debugln("Initializing integration module")
}

func getClientCID(address string) (string, error) {
	return "", nil
}

func saveClientInformation(client *clientsconfig.Client) error {
	return nil
}
