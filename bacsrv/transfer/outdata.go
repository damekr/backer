package transfer

import (
	"net"

	log "github.com/sirupsen/logrus"
)

// SendArchive sends an archive from storage to given client, here we are sure that client requests a fs.
func SendArchive(connection net.Conn, archName string) error {
	log.Debug("Starting sending archive: %s", archName)

	return nil
}
