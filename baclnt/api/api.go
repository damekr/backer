package api


import (
	"log"
	"net"

	"google.golang.org/grpc"
	pb "github.com/backer/bacsrv/api/proto"
	"io"
)

const (
	port = ":9000"
)
type server struct{}

func (s *server)TriggerBackup(stream pb.Baclnt_TriggerBackupServer) error{
	for {
		path, err := stream.Recv()
		if path != nil {
			log.Println("Path: ", path)
		}
		
		if err == io.EOF {
			return stream.SendAndClose(&pb.Status{
				Message: "OK",
			})
		}
		if err != nil {
			return nil
		}
	}

}

func ServeServer() {
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBaclntServer(s, &server{})
	s.Serve(lis)
}