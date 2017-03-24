package transfer

import (
	log "github.com/Sirupsen/logrus"
	"net"
)

// SendArchive sends an archive from repository to given client, here we are sure that client requests a backup.
func SendArchive(connection net.Conn, archName string) error {
	log.Debug("Starting sending archive: %s", archName)

	return nil
}
