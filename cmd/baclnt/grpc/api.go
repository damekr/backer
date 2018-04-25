package grpc

import (
	"context"
	"net"
	"time"

	"github.com/damekr/backer/api/protoclnt"
	"github.com/damekr/backer/cmd/baclnt/config"
	"github.com/damekr/backer/cmd/baclnt/fs"
	"github.com/damekr/backer/cmd/baclnt/network"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

var log = logrus.WithFields(logrus.Fields{"prefix": "api"})

// Ping returns hostname of client
func (s *server) Ping(ctx context.Context, in *protoclnt.PingRequest) (*protoclnt.PingResponse, error) {
	log.Printf("Got request to ping client: %s", in.Ip)
	s.metadataHandler(ctx)

	return &protoclnt.PingResponse{Message: "FROM CLIENT"}, nil
}

func (s *server) Backup(ctx context.Context, backupRequest *protoclnt.BackupRequest) (*protoclnt.BackupResponse, error) {
	log.Printf("Got request to backup client: %s", backupRequest.Ip)
	log.Println("Paths to be validated: ", backupRequest.Paths)
	s.metadataHandler(ctx)
	fileSystem := fs.NewLocalFileSystem()
	backupObjects, err := fileSystem.ReadBackupObjectsLocations(backupRequest.Paths)
	if err != nil {
		log.Errorln("Could not expand dirs for files, err: ", err)
	}
	log.Printf("Validated paths: ", backupObjects)

	// TODO backupRequest.IP is probably client ip, this message should contain server external ip
	err = runBackup(backupObjects)
	if err != nil {
		log.Errorf("Backup Failed, err: ", err.Error())
	}
	return &protoclnt.BackupResponse{Validpaths: backupObjects.Files}, nil
}

func runBackup(backupObjects fs.BackupObjects) error {
	// TODO Consider extend gRPC API to send client and server IP or external NAME it allows trigger backup from many servers to one client
	client := network.CreateTransferClient()
	session, err := client.Connect(config.MainConfig.ServerExternalName, config.MainConfig.ServerDataPort)
	if err != nil {
		log.Errorln("Cannot initialize connection, err: ", err.Error())
	}
	log.Println("Session ID: ", session.Id)
	err = session.StartBackup(backupObjects)
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
	log.Printf("Got request to run restore of client: %s", restoreRequest.Ip)
	log.Debugln("Need to restore from assetID: ", restoreRequest.AssetID)
	s.metadataHandler(ctx)
	startTime := time.Now()
	log.Println("Start time: ", startTime)
	var restoreMetadata bftp.AssetMetadata
	restoreMetadata.ID = int(restoreRequest.AssetID)

	var restoreOptions bftp.RestoreOptions
	restoreOptions.BasePath = restoreRequest.BasePath
	restoreOptions.WholeBackup = restoreRequest.WholeBackup
	restoreOptions.ObjectsPaths = restoreRequest.RestoreObjects
	err := runRestore(restoreMetadata, restoreOptions)
	if err != nil {
		log.Errorf("Restore Failed, err: ", err.Error())
	}
	return &protoclnt.RestoreResponse{Status: "OK"}, nil
}

func runRestore(assetMetadata bftp.AssetMetadata, options bftp.RestoreOptions) error {
	//TODO Consider extend gRPC API to send client and server IP or external NAME it allows trigger backup from many servers to one client
	client := network.CreateTransferClient()
	log.Debugln("Got restore files meta: ", assetMetadata)
	log.Debugln("Got restore options: ", options)
	session, err := client.Connect(config.MainConfig.ServerExternalName, config.MainConfig.ServerDataPort)
	if err != nil {
		log.Errorln("Could not connect to server for restore, err: ", err.Error())
		return err
	}
	log.Println("Session ID: ", session.Id)
	err = session.StartRestore(assetMetadata.ID, options)
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

// TODO Just for now, for having it in one place
func (s *server) metadataHandler(ctx context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
}

// Start method starts a grpc server on specific port
func Start() error {
	lis, err := net.Listen("tcp", ":"+config.MainConfig.ManagementPort)
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
