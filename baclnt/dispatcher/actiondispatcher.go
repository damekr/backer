package dispatcher

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/integration"
	"github.com/damekr/backer/baclnt/transfer"
)

// DataPort TODO shall be excluded to config file, or should be received from server during an integration
const DataPort = "8000"

// TODO: Think about authentication here

func DispatchBackupStart(paths []string, serverAddress string) error {
	log.Debugf("Establishing connection with: %s, on port %s", serverAddress, DataPort)
	err := transfer.StartFullBackup(paths, serverAddress)
	if err != nil {
		log.Error("An error occured during sending backup files")
		return err
	}
	return nil

}

func DispatchRestoreStart(paths []string, serverAddress string) error {
	log.Debugf("Dispatching restore for paths: ", paths)
	log.Debugf("Dispatching request from server: ", serverAddress)
	return nil

}

func DispatchIntegration() *integration.ClientInfo {
	log.Debug("Dispatching client integration")
	return integration.GetClientInfo()
}
