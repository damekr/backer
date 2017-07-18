package job

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/outprotoapi"
)

// ResponseClientHello is sending response to client ping-pong
func ResponseClientHello(clntAddress string) {
	out, err := outprotoapi.SayHelloToClient(clntAddress)
	if err != nil {
		log.Errorf("Given client on address %s is not available", clntAddress)
	}
	log.Debug("Got name from client: ", out)

}
