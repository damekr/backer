package api


import (
	"net"
	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	pb "github.com/backer/bacsrv/api/proto"
	"github.com/backer/baclnt/transfer"
	"io"
	"os"
)



const (
	port = ":9000"
)

func init(){
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

type server struct{}

func (s *server)TriggerBackup(stream pb.Baclnt_TriggerBackupServer) error{
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

func ServeServer() {
	lis, err := net.Listen("tcp", port)
	log.Info("Starting baclnt api on addr: ", lis.Addr())
	if err != nil {
		log.Errorf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBaclntServer(s, &server{})
	s.Serve(lis)
}