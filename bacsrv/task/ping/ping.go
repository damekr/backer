package ping

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/network"
	"github.com/damekr/backer/common/protosrv"
)

type Ping struct {
	ClientIP string
	Progress int
	Message  string
}

func CreatePing(clientIP string) *Ping {
	return &Ping{
		ClientIP: clientIP,
	}
}

func (p *Ping) Run() {
	log.Println("Pinging client: ", p.ClientIP)
	conn, err := network.EstablishGRPCConnection(p.ClientIP)
	if err != nil {
		log.Warningf("Cannot connect to address %s", p.ClientIP)

	}
	defer conn.Close()
	c := protosrv.NewBacsrvClient(conn)
	r, err := c.Ping(context.Background(), &protosrv.PingRequest{Ip: "Message from server"})
	if err != nil {
		log.Warningf("Could not get client name: %v", err)
	}
	log.Debugf("Received client message: %s", r.Message)
	p.Message = r.Message

}

func (p *Ping) Stop() {
	log.Println("Stopping")
}

const (
	clntMgmtPort = "9090"
	//timestampFormat = time.StampNano
)
