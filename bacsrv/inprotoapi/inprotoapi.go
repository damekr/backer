package inprotoapi

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/operations"
	"github.com/damekr/backer/common/protosrv"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	// "os"
)

type server struct{}

// SayHello returns hostname of client
func (s *server) SayHello(ctx context.Context, in *protosrv.HelloRequest) (*protosrv.HelloReply, error) {
	log.Printf("Got request from client: %s", in.Name)
	md, ok := metadata.FromContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
	go operations.SendHelloMessageToClient(in.Name)
	return &protosrv.HelloReply{Name: config.GetExternalName()}, nil
}

// ServeServer method starts a grpc server on specific port
func ServeServer(config *config.ServerConfig) {
	lis, err := net.Listen("tcp", ":"+config.MgmtPort)
	log.Info("Starting bacsrv protoapi on addr: ", lis.Addr())
	if err != nil {
		log.Errorf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protosrv.RegisterBacsrvServer(s, &server{})
	s.Serve(lis)
}
