package outprotoapi

import (
	"errors"
	// "os"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	pb "github.com/damekr/backer/common/protoclnt"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// TODO - GENERAL - Should be considered if main messages functions should get specific STRUCTs Like BackupMessageConfig

const (
	clntMgmtPort = ":9090"
)

func SayHelloToClient(address string) (string, error) {
	conn, err := grpc.Dial(address+clntMgmtPort, grpc.WithInsecure())
	if err != nil {
		log.Warningf("Cannot connect to address %s", address)
		return "", err
	}
	defer conn.Close()
	c := pb.NewBaclntClient(conn)
	//Contact the server and print out its response.
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: config.GetExternalName()})
	if err != nil {
		log.Warningf("Could not get client name: %v", err)
		return "", err
	}
	log.Debugf("Received client name: %s", r.Name)
	return r.Name, nil
}

func triggerCheckingPaths(client pb.BaclntClient, paths []*pb.Paths) {
	stream, err := client.GetStatusPaths(context.Background())
	if err != nil {
		log.Errorf("Cannot establish stream for path checking connection")
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
	address := clientaddr + clntMgmtPort
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Warnf("Cannot connect to client: %s", err)
	}
	defer conn.Close()
	c := pb.NewBaclntClient(conn)
	pbpaths := preparePaths(paths)
	log.Debug(c, pbpaths)
}

func makePbPaths(path string) *pb.Paths {
	return &pb.Paths{
		Name: config.GetExternalName(),
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
	address := clntAddress + clntMgmtPort
	log.Debug("Sending backup request to client: ", address)
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

func prepareRestoreTriggerMessage(reqcapacity int64, startlistener bool) *pb.TriggerRestoreMessage {
	return &pb.TriggerRestoreMessage{
		Name:          config.GetExternalName(),
		Reqcapacity:   reqcapacity,
		Startlistener: startlistener,
	}
}

func triggerRestore(client pb.BaclntClient, pbMessage *pb.TriggerRestoreMessage) error {
	response, err := client.TriggerRestore(context.Background(), pbMessage)
	if err != nil {
		log.Error("Cannot send restore message, error: ", err)
		return err
	}
	if response.Ok {
		log.Debug("Client %s has enough space", response.Name)
	} else if response.Listenerok {
		log.Debug("Sterted data listener on client %s side", response.Name)
	} else {
		log.Error("There is problem with space or starting data listener on client: ", response.Name)
		return errors.New("Cannot start data listener or client has not enough space")
	}
	return nil
}

func SendRestoreRequest(reqcapacity int64, startlistener bool, clntAddress string) error {
	address := clntAddress + clntMgmtPort
	log.Debug("Sending restore request to client: ", address)
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	defer conn.Close()

	if err != nil {
		log.Error("Could not estabilsh connection with client: ", address)
		return err
	}

	clnt := pb.NewBaclntClient(conn)
	pbMessage := prepareRestoreTriggerMessage(reqcapacity, startlistener)

	err = triggerRestore(clnt, pbMessage)
	if err != nil {
		log.Error("Cannot trigger restore on client: ", address)
		return err
	}

	return nil
}

func sendGrpcPathsToClient(clnt pb.BaclntClient, paths []*pb.Paths) error {
	stream, err := clnt.SendRestorePaths(context.Background())
	if err != nil {
		log.Error("Cannot send restore paths, error: ", err.Error())
		return err
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
	log.Debugf("Replay summary: %v", reply)
	return nil
}

func SendRestorePaths(paths []string, clientAddr string) error {
	log.Debugf("Sending paths %s to be restored to client %s", paths, clientAddr)
	address := clientAddr + clntMgmtPort
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Errorf("Could not establish connection with client: ", address)
		return err
	}
	grpcPaths := preparePaths(paths)
	clnt := pb.NewBaclntClient(conn)
	err = sendGrpcPathsToClient(clnt, grpcPaths)
	if err != nil {
		log.Errorf("Error occured during sending paths to be restored, error content: ", err)
		return err
	}
	return nil
}
