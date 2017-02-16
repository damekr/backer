package dispatcher

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/archiver"
	"os"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func DispatchBackupStart(paths []string, serverAddress string) {
	archive := archiver.NewArchive(paths)
	tarlocation := archive.MakeArchive()
	log.Debugf("An archive has been created at location %s", tarlocation)

}
