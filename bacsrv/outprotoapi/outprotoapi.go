package outprotoapi

import (
	"time"
)

// TODO - GENERAL - Should be considered if main messages functions should get specific STRUCTs Like BackupMessageConfig

const (
	clntMgmtPort    = "9090"
	timestampFormat = time.StampNano
)

//func establishConnection(clntAddr string) (*grpc.ClientConn, error) {
//	log.Debug("Establishing grpc connection with: ", clntAddr)
//	conn, err := grpc.Dial(net.JoinHostPort(clntAddr, clntMgmtPort), grpc.WithInsecure())
//	if err != nil {
//		log.Error("Cannot establish connection with client: ", clntAddr)
//		return nil, err
//	}
//	log.Debug("Successfully established grpc connection with client: ", clntAddr)
//
//	return conn, nil
//}
//
//func SayHelloToClient(clntAddress string) (string, error) {
//	conn, err := establishConnection(clntAddress)
//	if err != nil {
//		log.Warningf("Cannot connect to address %s", clntAddress)
//		return "", err
//	}
//	defer conn.Close()
//	c := pb.NewBaclntClient(conn)
//	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: config.GetExternalName()})
//	if err != nil {
//		log.Warningf("Could not get client name: %v", err)
//		return "", err
//	}
//	log.Debugf("Received client name: %s", r.Name)
//	return r.Name, nil
//}
//
//func CheckPaths(clntAddr string, paths []string) ([]string, error) {
//	log.Debug("Starting checking paths on client: ", clntAddr)
//	md := metadata.Pairs("ServerExternalName", config.GetExternalName())
//	ctx := metadata.NewContext(context.Background(), md)
//	conn, err := establishConnection(clntAddr)
//	if err != nil {
//		log.Error("Cannot connect to client: ", clntAddr)
//	}
//	defer conn.Close()
//	con := pb.NewBaclntClient(conn)
//	availableFiles, err := con.CheckPaths(ctx, makePbPaths(paths))
//	if err != nil {
//		log.Error("Cannot check paths, got error: ", err.Error())
//		return nil, err
//	}
//	log.Debug("ClientConfig respond with his name: ", availableFiles.Name)
//	log.Debug("Got validated, resolved paths of requested files: ", availableFiles.Path)
//	return availableFiles.Path, nil
//}
//
//// SendBackupRequest creates connection to client with specified address and trigger a backup
//func SendBackupRequest(paths []string, clntAddress string) error {
//	log.Debug("Starting backup on client: ", clntAddress)
//	md := metadata.Pairs("ServerExternalName", config.GetExternalName())
//	ctx := metadata.NewContext(context.Background(), md)
//	conn, err := establishConnection(clntAddress)
//	if err != nil {
//		log.Warningf("Cannot connect to address %s", clntAddress)
//		return err
//	}
//	defer conn.Close()
//	c := pb.NewBaclntClient(conn)
//	status, err := c.TriggerBackup(ctx, &pb.Paths{
//		Name: config.GetExternalName(),
//		Path: paths,
//	})
//	if err != nil {
//		log.Error("Couldn't send trigger backup message")
//		return err
//	}
//	log.Debugf("Got status message: %#v from client: %s", status, clntAddress)
//	return nil
//}
//
//func triggerCheckingPaths(client pb.BaclntClient, paths []*pb.Paths) {
//	stream, err := client.GetStatusPaths(context.Background())
//	if err != nil {
//		log.Errorf("Cannot establish stream for path checking connection")
//	}
//	for _, path := range paths {
//		if err := stream.Send(path); err != nil {
//			log.Errorf("Cannot send path  %s over stream", path)
//		}
//	}
//	// TODO Bellow probably does not work, was changed bacause of compiling errors.
//	reply, err := stream.Recv()
//	if err != nil {
//		log.Errorf("Cannot get stream replay during checking paths, err: %v", err)
//	}
//	log.Debugf("Received paths from client informations: %v", reply)
//
//}
//
//func makePbPaths(path []string) *pb.Paths {
//	return &pb.Paths{
//		Name: config.GetExternalName(),
//		Path: path,
//	}
//}
//
//func prepareRestoreTriggerMessage(reqcapacity int64, startlistener bool) *pb.TriggerRestoreMessage {
//	return &pb.TriggerRestoreMessage{
//		Name:          config.GetExternalName(),
//		Reqcapacity:   reqcapacity,
//		Startlistener: startlistener,
//	}
//}
//
//func triggerRestore(client pb.BaclntClient, pbMessage *pb.TriggerRestoreMessage) error {
//	response, err := client.TriggerRestore(context.Background(), pbMessage)
//	if err != nil {
//		log.Error("Cannot send restore message, error: ", err)
//		return err
//	}
//	if response.Ok {
//		log.Debug("ClientConfig %s has enough space", response.Name)
//	} else if response.Listenerok {
//		log.Debug("Sterted data listener on client %s side", response.Name)
//	} else {
//		log.Error("There is problem with space or starting data listener on client: ", response.Name)
//		return errors.New("Cannot start data listener or client has not enough space")
//	}
//	return nil
//}
//
//func SendRestoreRequest(reqcapacity int64, startlistener bool, clntAddress string) error {
//	address := clntAddress + clntMgmtPort
//	log.Debug("Sending restore rpc request to client: ", address)
//	conn, err := grpc.Dial(address, grpc.WithInsecure())
//	defer conn.Close()
//
//	if err != nil {
//		log.Error("Could not estabilsh connection with client: ", address)
//		return err
//	}
//
//	clnt := pb.NewBaclntClient(conn)
//	pbMessage := prepareRestoreTriggerMessage(reqcapacity, startlistener)
//
//	err = triggerRestore(clnt, pbMessage)
//	if err != nil {
//		log.Error("Cannot trigger restore on client: ", address)
//		return err
//	}
//
//	return nil
//}
//
//func sendGrpcPathsToClient(clnt pb.BaclntClient, paths []*pb.Paths) error {
//	stream, err := clnt.SendRestorePaths(context.Background())
//	if err != nil {
//		log.Error("Cannot send restore paths, error: ", err.Error())
//		return err
//	}
//	for _, path := range paths {
//		if err := stream.Send(path); err != nil {
//			log.Errorf("Error %v", err)
//		}
//	}
//	reply, err := stream.CloseAndRecv()
//	if err != nil {
//		log.Errorf("An error occured: %v", err)
//	}
//	log.Debugf("Replay summary: %v", reply)
//	return nil
//}
//
//// SendIntegrationRequest sends a request to client to get needed vales from the client
//func SendIntegrationRequest(client *config.Client) (*config.Client, error) {
//	log.Debug("Sending rpc integration messagee to client with address: ", client.Address)
//	md := metadata.Pairs("ServerExternalName", config.GetExternalName())
//	ctx := metadata.NewContext(context.Background(), md)
//	conn, err := establishConnection(client.Address)
//	if err != nil {
//		log.Error("Couldn't establish rpc connection with client: ", client.Address)
//		return client, err
//	}
//	defer conn.Close()
//	c := pb.NewBaclntClient(conn)
//	clientInfo, err := c.TriggerIntegration(ctx, &pb.HelloRequest{Name: config.GetExternalName()})
//	if err != nil {
//		log.Errorf("Couldn't read client CID from remote client")
//		return client, err
//	}
//	log.Debug("Got information about client: %#v", clientInfo)
//	client.Name = clientInfo.Name
//	client.CID = clientInfo.Cid
//	client.Platform = clientInfo.Platform
//	return client, nil
//}
