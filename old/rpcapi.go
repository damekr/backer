package api

import (
	"net"
	"os"

	pb "github.com/damekr/backer/bacsrv/api/proto"
	"google.golang.org/grpc"

	log "github.com/Sirupsen/logrus"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

// StartApi is main handler for external connection to server, it also will handle requests from
// client application
func StartApiListener(port string) {
	lis, err := net.Listen("tcp", port)
	log.Info("Starting listen server api on addr: ", lis.Addr())
	if err != nil {
		log.Errorf("Failed to listen, error: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterBaclntServer(s, &server{})
	s.Serve(lis)

