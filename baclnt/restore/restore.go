package restore

import (
	"archive/tar"
	log "github.com/Sirupsen/logrus"
)

func init() {
	log.Debugln("Initializes restore module")
}

func Archive(archive *tar.Reader) error {
	return nil
}
