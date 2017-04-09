package inprotoapi

import (
	"io"
	"net"
	// "os"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/config"
	"github.com/damekr/backer/baclnt/dispatcher"
	pb "github.com/damekr/backer/common/protoclnt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

func (s *server) CheckPaths(ctx context.Context, in *pb.Paths) (*pb.Paths, error) {
	log.Debug("Received a request to check paths")
	requestedPaths := in.Path
	log.Debug("Got paths to be checked: ", requestedPaths)
	log.Debug("Starting checking paths it can take a while...")
	validateFilesPaths := dispatcher.ValidatePaths(requestedPaths)
	return &pb.Paths{
		Name: config.GetExternalName(),
		Path: validateFilesPaths,
	}, nil
}

// SayHello returns hostname of client
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Got request from server: %s", in.Name)
	return &pb.HelloReply{Name: config.GetExternalName()}, nil
}

func (s *server) TriggerBackup(ctx context.Context, in *pb.Paths) (*pb.Status, error) {
	md, ok := metadata.FromContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
	log.Debugf("Backup from %s has been triggered", in.Name)
	log.Debugf("Got paths to be send: ", in.Path)
	err := dispatcher.DispatchBackupStart(in.Path, in.Name)
	if err != nil {
		log.Error("An error occured during dispatching data transfer, error: ", err.Error())
	}
	return &pb.Status{
		Name:    config.GetExternalName(),
		Message: "ok",
	}, nil
}

func (s *server) TriggerRestore(ctx context.Context, request *pb.TriggerRestoreMessage) (*pb.TriggerRestoreResponse, error) {
	log.Debug("Got restore trigger with requested capacity", request.Reqcapacity)
	return &pb.TriggerRestoreResponse{Name: config.GetExternalName(), Ok: true, Listenerok: true}, nil
}

func (s *server) GetStatusPaths(stream pb.Baclnt_GetStatusPathsServer) error {
	log.Debug("Starting checking paths")
	// var paths []string
	for {
		path, err := stream.Recv()
		if path != nil {
			log.Debug("Received path to check: ", path.Path)
			// paths = append(paths, path.Path)
		}
		if err == io.EOF {
			log.Debug("Received all paths, checking locally paths..")
		}

	}
}

func (s *server) SendRestorePaths(pathsStream pb.Baclnt_SendRestorePathsServer) error {
	log.Debugf("Getting restore paths in stream...")
	var (
		paths         []string
		serverAddress string
	)
	for {
		path, err := pathsStream.Recv()
		if path != nil {
			log.Debugf("Received path to be restored: %s", path.Path)
			// paths = append(paths, path.Path)
			//  TODO The same case as in backup, maybe consider sending to messages --> hello with authentication and then paths
			serverAddress = path.Name
		}
		if err == io.EOF {
			// TODO here always is sending "OK" message should be executed some checks and then send a proper message
			log.Debugf("Received all paths from server: %s, sending ok message to server", serverAddress)
			dispatcher.DispatchRestoreStart(paths, serverAddress)
			return pathsStream.SendAndClose(&pb.HelloReply{
				Name: config.GetExternalName(),
			})
		}
	}
}

// ServeServer method starts a grpc server on specific port
func ServeServer(config *config.ClientConfig) {
	lis, err := net.Listen("tcp", ":"+config.MgmtPort)
	log.Info("Starting baclnt protoapi on addr: ", lis.Addr())
	if err != nil {
		log.Errorf("Failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBaclntServer(s, &server{})
	s.Serve(lis)
}
