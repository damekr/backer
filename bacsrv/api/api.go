package api

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/backer/bacsrv/api/proto"
	"github.com/backer/bacsrv/config"
	log "github.com/Sirupsen/logrus"
	"os"
)

func init(){
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func makePbPaths(path string) *pb.Paths{
	return &pb.Paths{Path: path}
}

func preparePaths(paths []string) []*pb.Paths{
	var pbpaths []*pb.Paths
	for _, l := range paths{
		pbpaths = append(pbpaths, makePbPaths(l))
	}
	return pbpaths
}

func triggerBackup(client pb.BaclntClient, paths []*pb.Paths) {
	stream, err := client.TriggerBackup(context.Background())
	if err != nil {
		log.Errorf("Could not greet: %v", err)
	}
	for _, path := range paths{
		if err := stream.Send(path); err != nil{
			log.Errorf("Error %v", err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Errorf("An error occured: %v", err)
	}
	log.Debugf("Route summary: %v", reply)
}


func SendBackupRequest(paths []string){
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

