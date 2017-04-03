package outprotoapi

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/common/protosrv"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	// "os"
)

// GENERAL COMMENT: This file needs to use types of message from bacsrv proto file

const (
	srvMgmtPort = ":8090"
)

var HostName string

func init() {
	// name, err := os.Hostname()
	// if err != nil {
	// 	log.Warning("Cannot get hostname, setting default: baclnt")
	// 	name = "baclnt"
	// }
	HostName = "127.0.0.1"
}

func SayHelloToServer(address string) (string, error) {
	log.Printf("Sending message to: %s%s", address, srvMgmtPort)
	conn, err := grpc.Dial(address+srvMgmtPort, grpc.WithInsecure())
	log.Print("ERROR: ", err)
	if err != nil {
		log.Warningf("Cannot connect to address %s", address)
		return "", err
	}
	defer conn.Close()
	c := protosrv.NewBacsrvClient(conn)
	//Contact the server and print out its response.
	r, err := c.SayHello(context.Background(), &protosrv.HelloRequest{Name: HostName})
	if err != nil {
		log.Warningf("Could not get client name: %v", err)
		return "", err
	}
	log.Debugf("Received client name: %s", r.Name)
	return r.Name, nil
}
