package api



import (
	"log"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "github.com/backer/bacsrv/api/proto"
	"github.com/backer/bacsrv/config"
)


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
		log.Fatalf("could not greet: %v", err)
	}
	for _, path := range paths{
		if err := stream.Send(path); err != nil{
			log.Fatalf("Error %v", err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("An error occured: %v", err)
	}
	log.Printf("Route summary: %v", reply)
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

