package ping

import (
	"context"

	"github.com/damekr/backer/bacsrv/network"
	"github.com/damekr/backer/common/proto"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "task:ping"})

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

// TODO Handle baclnt client not available.

func (p *Ping) Run() {
	log.Println("Pinging client: ", p.ClientIP)
	conn, err := network.EstablishGRPCConnection(p.ClientIP)
	if err != nil {
		log.Warningf("Cannot connect to address %s", p.ClientIP)

	}
	defer conn.Close()
	c := proto.NewBacsrvClient(conn)
	r, err := c.Ping(context.Background(), &proto.PingRequest{Ip: "Message from server"})
	if err != nil {
		log.Errorf("Could not connect to client err: %v", err)
		p.Message = err.Error()
		return
	}
	if r.Message != "" {
		log.Debugf("Received client message: %s", r.Message)
		p.Message = r.Message
	} else {
		p.Message = ""
	}

}

func (p *Ping) Stop() {
	log.Println("Stopping")
}
