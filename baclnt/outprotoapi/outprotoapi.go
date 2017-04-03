package outprotoapi

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/config"
	"github.com/damekr/backer/common/protosrv"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"time"
)

// GENERAL COMMENT: This file needs to use types of message from bacsrv proto file
// METADATA: Using metadata of grpc can be done authentication, below example how to proceed
const (
	srvMgmtPort     = ":8090"
	timestampFormat = time.StampNano
)

func SayHelloToServer(address string) (string, error) {
	md := metadata.Pairs("timestamp", time.Now().Format(timestampFormat))
	ctx := metadata.NewContext(context.Background(), md)
	log.Printf("Sending message to: %s%s, hostName: %s", address, srvMgmtPort, config.ClntConfig.ExternalName)
	conn, err := grpc.Dial(address+srvMgmtPort, grpc.WithInsecure())
	log.Print("ERROR: ", err)
	if err != nil {
		log.Warningf("Cannot connect to address %s", address)
		return "", err
	}
	defer conn.Close()
	c := protosrv.NewBacsrvClient(conn)
	//Contact the server and print out its response.
	r, err := c.SayHello(ctx, &protosrv.HelloRequest{Name: config.ClntConfig.ExternalName})
	if err != nil {
		log.Warningf("Could not get client name: %v", err)
		return "", err
	}
	log.Debugf("Received client name: %s", r.Name)
	return r.Name, nil
}
