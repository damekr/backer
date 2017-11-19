package ping

import (
	"fmt"
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
	fmt.Println("Starting backup of client: ", p.ClientIP)
	//conn, err := establishConnection(clientIP)
	//if err != nil {
	//	log.Warningf("Cannot connect to address %s", clientIP)
	//	return "", err
	//}
	//defer conn.Close()
	//c := pb.NewBaclntClient(conn)
	//r, err := c.Ping(context.Background(), &pb.PingRequest{Message: "Message from server"})
	//if err != nil {
	//	log.Warningf("Could not get client name: %v", err)
	//	return "", err
	//}
	//log.Debugf("Received client message: %s", r.Message)
	//return r.Message, nil
	p.Message = "SADA"

}

func (p *Ping) Stop() {
	fmt.Println("Stopping")
}
