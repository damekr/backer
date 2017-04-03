package dispatcher

import (
	log "github.com/Sirupsen/logrus"
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
