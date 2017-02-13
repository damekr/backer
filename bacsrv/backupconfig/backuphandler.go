package backupconfig

import (
	log "github.com/Sirupsen/logrus"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)

}

type Backup struct {
	ID        string
	Paths     []string
	Excluded  []string
	Retention string
}
