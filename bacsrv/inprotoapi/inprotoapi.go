package inprotoapi

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/manager"
	"github.com/damekr/backer/common/protosrv"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"net"
	// "os"
)

var HostName string

type server struct{}

// TODO to make it visible from outside network must be able to read from config file.
func init() {
	// name, err := os.Hostname()
	// if err != nil {
	// 	log.Warning("Cannot get hostname, setting default: baclnt")
	// 	name = "baclnt"
	// }
	HostName = "127.0.0.1"
}

// SayHello returns hostname of client
func (s *server) SayHello(ctx context.Context, in *protosrv.HelloRequest) (*protosrv.HelloReply, error) {
	log.Printf("Got request from client: %s", in.Name)
	go manager.SendHelloMessageToClient(in.Name)
	return &protosrv.HelloReply{Name: HostName}, nil
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
