package api

import (
	"io"
	"net"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/config"
	"github.com/damekr/backer/baclnt/transfer"
	pb "github.com/damekr/backer/bacsrv/protoapi/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type server struct{}

// SayHello returns hostname of client
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	name, err := os.Hostname()
	if err != nil {
		log.Error("Cannot get hostname, sending default value")
		return &pb.HelloReply{Name: "0"}, nil
	}
	return &pb.HelloReply{Name: name}, nil
}

func (s *server) TriggerBackup(stream pb.Baclnt_TriggerBackupServer) error {
	log.Debug("Backup has been triggered")
	var paths []string
	for {
		path, err := stream.Recv()
		if path != nil {
			log.Debug("Received path to backup: ", path.Path)
			paths = append(paths, path.Path)
		}
		if err == io.EOF {
			log.Debug("Recivied all paths, sending OK message to server...")
			return stream.SendAndClose(&pb.Status{
				Message: "OK",
			})
		}
		absPaths := transfer.GetAbsolutePaths(paths)
		log.Debugf("Absolutive paths: %v", absPaths)

		if err != nil {
			return nil
		}
	}

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

func ServeServer(config *config.ClientConfig) {
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
