package api

import (
	"io"
	"net"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/config"
	"github.com/damekr/backer/baclnt/dispatcher"
	pb "github.com/damekr/backer/bacsrv/protoapi/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

var HostName string

func init() {
	name, err := os.Hostname()
	if err != nil {
		log.Warning("Cannot get hostname, setting default: baclnt")
		name = "baclnt"
	}
	HostName = name
}

type server struct{}

// SayHello returns hostname of client
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Printf("Got request from server: %s", in.Name)
	return &pb.HelloReply{Name: HostName}, nil
}

func (s *server) TriggerBackup(stream pb.Baclnt_TriggerBackupServer) error {
	log.Debug("Backup has been triggered")
	var paths []string
	var serverName string
	for {
		path, err := stream.Recv()
		if path != nil {
			log.Debugf("Received path to backup: %s, from server: %s", path.Path, path.Name)
			paths = append(paths, path.Path)
			// TODO It will be overrided on each iteration, should be improved, or each path will have name specification
			serverName = path.Name
		}
		if err == io.EOF {
			go dispatcher.DispatchBackupStart(paths, serverName)
			log.Debug("Recivied all paths, sending OK message to server...")
			return stream.SendAndClose(&pb.Status{
				Name:    HostName,
				Message: "OK",
			})
		}

	}

}

func (s *server) TriggerRestore(ctx context.Context, request *pb.TriggerRestoreMessage) (*pb.TriggerRestoreResponse, error) {
	log.Debug("Got restore trigger with requested capacity", request.Reqcapacity)
	return &pb.TriggerRestoreResponse{Name: HostName, Ok: true, Listenerok: true}, nil
}

func (s *server) GetStatusPaths(stream pb.Baclnt_GetStatusPathsServer) error {
	log.Debug("Starting checking paths")
	var paths []string
	for {
		path, err := stream.Recv()
		if path != nil {
			log.Debug("Received path to check: ", path.Path)
			paths = append(paths, path.Path)
		}
		if err == io.EOF {
			log.Debug("Received all paths, checking locally paths..")
		}

	}
}

// ServeServer method starts a grpc server on specific port
func ServeServer(config *config.ClientConfig) {
	// TODO Find a way to get client informations.
	//listOfFiles := []string{"/home/damekr/d8x.github.io"}
	//files := transfer.GetFilesInformations(listOfFiles)
	//log.Debug("Files ", files)
	lis, err := net.Listen("tcp", ":"+config.MgmtPort)
	log.Info("Starting baclnt api on addr: ", lis.Addr())
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBaclntServer(s, &server{})
	s.Serve(lis)
}
