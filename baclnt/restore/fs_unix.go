package restore

import (
	log "github.com/Sirupsen/logrus"
	"io"
)

func init() {
	log.Debugln("Initializes restore for unix fs")
}

// SaveFileOnUnixFS saves file on unix fs
func SaveFileOnUnixFS(w *io.Writer, override bool) error {
	return nil
}
