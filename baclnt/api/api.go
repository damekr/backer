package api

import (
	"context"
	"net"

	"github.com/damekr/backer/baclnt/config"
	"github.com/damekr/backer/baclnt/fs"
	"github.com/damekr/backer/baclnt/network"
	"github.com/damekr/backer/common/protoclnt"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

var log = logrus.WithFields(logrus.Fields{"prefix": "api"})

// Ping returns hostname of client
func (s *server) Ping(ctx context.Context, in *protoclnt.PingRequest) (*protoclnt.PingResponse, error) {
	log.Printf("Got request to ping client: %s", in.Ip)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)

	return &protoclnt.PingResponse{Message: "FROM CLIENT"}, nil
}

func (s *server) Backup(ctx context.Context, backupRequest *protoclnt.BackupRequest) (*protoclnt.BackupResponse, error) {
	log.Printf("Got request to fs client: %s", backupRequest.Ip)
	log.Println("Paths to be validated: ", backupRequest.Paths)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
	fileSystem := fs.FS{}
	validatedPaths := fileSystem.GetAbsolutePaths(backupRequest.Paths)
	log.Printf("Validated paths: ", validatedPaths)

	//TODO backupRequest.IP is probably client ip, this message should contains server external ip
	err := runBackup(validatedPaths, backupRequest.Ip)
	if err != nil {
		log.Errorf("Backup Failed, err: ", err.Error())
	}
	return &protoclnt.BackupResponse{Validpaths: validatedPaths}, nil
}

func runBackup(paths []string, serverIp string) error {
	client := network.CreateTransferClient()
	session, err := client.Connect(serverIp, config.MainConfig.ServerDataPort)
	if err != nil {
		log.Errorln("Cannot initialize connection")
	}
	log.Println("Session ID: ", session.Id)
	err = session.StartBackup(paths)
	if err != nil {
		log.Error("Backup failed, err: ", err.Error())
	}
	//TODO defer?
	log.Debug("Finished backup, closing session...")
	err = session.CloseSession()
	if err != nil {
		log.Errorln("Could not close session, err: ", err.Error())
	}
	return nil
}

func (s *server) Restore(ctx context.Context, restoreRequest *protoclnt.RestoreRequest) (*protoclnt.RestoreResponse, error) {
	log.Printf("Got request to fs client: %s", restoreRequest.Ip)
	log.Println("Paths to be validated: ", restoreRequest.Paths)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)

	err := runRestore(restoreRequest.Paths, restoreRequest.Ip)
	if err != nil {
		log.Errorf("Backup Failed, err: ", err.Error())
	}
	return &protoclnt.RestoreResponse{Status: "OK"}, nil
}

func runRestore(paths []string, serverIp string) error {
	client := network.CreateTransferClient()
	session, err := client.Connect(serverIp, config.MainConfig.ServerDataPort)
	if err != nil {
		log.Errorln("Could not connect to server for restore, err: ", err.Error())
		return err
	}
	log.Println("Session ID: ", session.Id)
	err = session.StartRestore(paths)
	if err != nil {
		log.Errorln("Restore failed, err: ", err.Error())
	}
	//TODO defer?
	err = session.CloseSession()
	if err != nil {
		log.Errorln("Could not close session, err: ", err.Error())
		return err
	}
	return nil
}

// Start method starts a grpc server on specific port
func Start() error {
	lis, err := net.Listen("tcp", ":"+config.MainConfig.MgmtPort)
	if err != nil {
		log.Errorf("Failed to listen: %v", err)
		return err
	}
	log.Info("Starting bacsrv protoapi on addr: ", lis.Addr())
	s := grpc.NewServer()
	protoclnt.RegisterBaclntServer(s, &server{})
	s.Serve(lis)
	return nil
}
