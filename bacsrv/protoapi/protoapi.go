package protoapi

import (
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	pb "github.com/damekr/backer/bacsrv/protoapi/proto"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// TODO - GENERAL - Should be considered if main messages functions should get specific STRUCTs Like BackupMessageConfig

const (
	clntMgmtPort = ":9090"
)

var Name string

func init() {
	name, err := os.Hostname()
	if err != nil {
		log.Error("Cannot get server hostname, setting default")
		name = "bacsrv"
	}
	Name = name
}

func SayHelloToClient(address string) (string, error) {
	conn, err := grpc.Dial(address+clntMgmtPort, grpc.WithInsecure())
	if err != nil {
		log.Warningf("Cannot connect to address %s", address)
		return "", err
	}
	defer conn.Close()
	c := pb.NewBaclntClient(conn)
	//Contact the server and print out its response.
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: Name})
	if err != nil {
		log.Warningf("Could not get client name: %v", err)
		return "", err
	}
	log.Debugf("Received client name: %s", r.Name)
	return r.Name, nil
}

func makePbPaths(path string) *pb.Paths {
	return &pb.Paths{
		Name: Name,
		Path: path,
	}
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

// SendBackupRequest creates connection to client with specified address and trigger a backup
func SendBackupRequest(paths []string, clntAddress string) error {
	address := clntAddress + config.GetMgmtPort()
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Errorf("Could not connect to client: %s", address)
		return err
	}
	defer conn.Close()
	c := pb.NewBaclntClient(conn)
	pbpaths := preparePaths(paths)
	triggerBackup(c, pbpaths)
	return nil
}
