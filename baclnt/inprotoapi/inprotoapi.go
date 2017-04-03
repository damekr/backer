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
			// TODO dispatcher shall be triggered in the same goroutine, the work inside should be in goroutines
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
			paths = append(paths, path.Path)
			//  TODO The same case as in backup, maybe consider sending to messages --> hello with authentication and then paths
			serverAddress = path.Name
		}
		if err == io.EOF {
			// TODO here always is sending "OK" message should be executed some checks and then send a proper message
			log.Debugf("Received all paths from server: %s, sending ok message to server", serverAddress)
			dispatcher.DispatchRestoreStart(paths, serverAddress)
			return pathsStream.SendAndClose(&pb.HelloReply{
				Name: HostName,
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
