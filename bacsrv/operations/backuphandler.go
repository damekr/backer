package operations

import (
	log "github.com/Sirupsen/logrus"

	"github.com/damekr/backer/bacsrv/outprotoapi"
)

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
