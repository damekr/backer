package dispatcher

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/archiver"
	"github.com/damekr/backer/baclnt/transfer"
)

// DataPort TODO Must be excluded to config file, or should be received from server during an integration
const DataPort = 8000

func DispatchBackupStart(paths []string, serverAddress string) {
	archive := archiver.NewArchive(paths)
	tarlocation := archive.MakeArchive()
	log.Debugf("An archive has been created at location %s", tarlocation)
	transferConnection := transfer.TransferConnection{
		Host: serverAddress,
		Port: DataPort,
	}
	log.Debugf("Establishing connection with: %s, on port %s", serverAddress, DataPort)
	backupConfig := &transfer.BackupConfig{
		Paths:  paths,
		TRConn: transferConnection.InitConnection(),
	}
	backupConfig.SendArchive(tarlocation)
}
