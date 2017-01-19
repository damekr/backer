package api

import (
	"os"

	log "github.com/Sirupsen/logrus"
	pb "github.com/damekr/backer/bacsrv/api/proto"
	"github.com/damekr/backer/bacsrv/config"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func makePbPaths(path string) *pb.Paths {
	return &pb.Paths{Path: path}
}

func preparePaths(paths []string) []*pb.Paths {
	var pbpaths []*pb.Paths
	for _, l := range paths {
		pbpaths = append(pbpaths, makePbPaths(l))
	}
	return pbpaths
}

func triggerCheckingPaths(client pb.BaclntClient, paths []*pb.Paths) {
	stream, err := client.GetStatusPaths(context.Background())
	if err != nil {
		log.Errorf("Cannot estabilish stream for path checking connection")
	}
	for _, path := range paths {
		if err := stream.Send(path); err != nil {
			log.Errorf("Cannot send path  %s over stream", path)
		}
	}
	// TODO Bellow probably does not work, was changed bacause of compiling errors.
	reply, err := stream.Recv()
	if err != nil {
		log.Errorf("Cannot get stream replay during checking paths, err: %v", err)
	}
	log.Debugf("Received paths from client informations: %v", reply)

}

func CheckIfPathsExists(paths []string, clientaddr string) {
	address := clientaddr + ":" + config.GetMgmtPort()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Warnf("Cannot connect to client: %s", err)
	}
	defer conn.Close()
	c := pb.NewBaclntClient(conn)
	pbpaths := preparePaths(paths)
	log.Debug(c, pbpaths)
}

func triggerBackup(client pb.BaclntClient, paths []*pb.Paths) {
	stream, err := client.TriggerBackup(context.Background())
	if err != nil {
		log.Errorf("Could not greet: %v", err)
	}
	for _, path := range paths {
		if err := stream.Send(path); err != nil {
			log.Errorf("Error %v", err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Errorf("An error occured: %v", err)
	}
	log.Debugf("Route summary: %v", reply)
}

func SendBackupRequest(paths []string) {
	address := "localhost:" + config.GetMgmtPort()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewBaclntClient(conn)
	pbpaths := preparePaths(paths)
	triggerBackup(c, pbpaths)
}
