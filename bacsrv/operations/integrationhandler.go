package operations

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
)

func init() {
	log.Debugln("Initializing integration module")
}

func getClientCID(address string) (string, error) {
	return "", nil
}

func saveClientInformation(client *config.Client) error {
	return nil
}
