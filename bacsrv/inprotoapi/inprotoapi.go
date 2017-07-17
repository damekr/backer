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
	// go operations.SendHelloMessageToClient(in.Name)
	return &protosrv.HelloReply{Name: config.GetExternalName()}, nil
}

func (s *server) ListClients(ctx context.Context, in *protosrv.HelloRequest) (*protosrv.ClientsList, error) {
	log.Debug("Received a request to check paths")
	client := in.Name
	log.Debug("Got request from: ", client)
	log.Debug("Starting checking integrated clients")
	clientsL := operations.GetAllIntegratedClients()
	clients := []string{}
	for _, v := range clientsL {
		clients = append(clients, v.Name)
	}
	return &protosrv.ClientsList{
		Clients: clients,
	}, nil
}

func (s *server) RunBackup(ctx context.Context, in *protosrv.Client) (*protosrv.Status, error) {
	log.Debug("Received a request to run backup of client: ", in.Cname)
	client := in.Name
	log.Debug("Got request from: ", client)
	log.Debug("Starting backup of: ", in.Cname)
	clientInfo := config.GetClientInformation(in.Cname)
	log.Debug("Running backup of client address: ", client)

	backup := &confi
	err := operations.StartBackup(backup, clientHost)
	if err != nil {
		log.Error("Cannot start backup, error: ", err)
	}
	return &protosrv.Status{
		Backup: true,
	}, nil
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
