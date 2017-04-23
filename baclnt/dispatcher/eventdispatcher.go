package dispatcher

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/backup"
	"github.com/damekr/backer/baclnt/outprotoapi"
)

func SendHelloMessageToServer(srvAddress string) error {
	out, err := outprotoapi.SayHelloToServer(srvAddress)
	if err != nil {
		log.Error("Could not send hello message to server: ", srvAddress)
		return err
	}
	log.Debugf("Client %s is saying hello", out)
	return nil
}

// ValidatePaths got paths and checks if they exist in the system, returns only available.
func ValidatePaths(paths []string) []string {
	log.Debug("Starting paths validation, received paths: ", paths)
	validatedPaths := backup.GetAbsolutePaths(paths)
	log.Debug("Checked and resolved all files: ", validatedPaths)
	return validatedPaths
}
