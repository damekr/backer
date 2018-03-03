package api

import (
	"context"
	"net"
	"strings"

	"github.com/damekr/backer/baclnt/config"
	"github.com/damekr/backer/baclnt/fs"
	"github.com/damekr/backer/baclnt/network"
	"github.com/damekr/backer/baclnt/transfer"
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
	log.Printf("Got request to backup client: %s", backupRequest.Ip)
	log.Println("Paths to be validated: ", backupRequest.Paths)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
	fileSystem := fs.FS{}
	validatedPaths := fileSystem.GetAbsolutePaths(backupRequest.Paths)
	log.Printf("Validated paths: ", validatedPaths)

	//TODO backupRequest.IP is probably client ip, this message should contains server external ip
	err := runBackup(validatedPaths)
	if err != nil {
		log.Errorf("Backup Failed, err: ", err.Error())
	}
	return &protoclnt.BackupResponse{Validpaths: validatedPaths}, nil
}

func runBackup(paths []string) error {
	// TODO Consider extend gRPC API to send client and server IP or external NAME it allows trigger backup from many servers to one client
	client := network.CreateTransferClient()
	session, err := client.Connect(config.MainConfig.ServerExternalName, config.MainConfig.ServerDataPort)
	if err != nil {
		log.Errorln("Cannot initialize connection, err: ", err.Error())
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
	log.Printf("Got request to run restore of client: %s", restoreRequest.Ip)
	log.Debugln("Need to restore files on server: ", restoreRequest.PathsOnServer)
	log.Println("Original files locations: ", restoreRequest.OriginalPaths)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
	var restoreFilesMetadata []transfer.RestoreFileMetadata
	// TODO Fast fix to match paths - does not work! FIXME
	for _, v := range restoreRequest.PathsOnServer {
		for _, k := range restoreRequest.OriginalPaths {
			if strings.ContainsAny(v, k) {
				fileMetadata := transfer.RestoreFileMetadata{
					PathOnServer: v,
					PathOnClient: k,
				}
				restoreFilesMetadata = append(restoreFilesMetadata, fileMetadata)
			}
		}
	}
	err := runRestore(restoreFilesMetadata)
	if err != nil {
		log.Errorf("Restore Failed, err: ", err.Error())
	}
	return &protoclnt.RestoreResponse{Status: "OK"}, nil
}

func runRestore(restoreFilesMetadatas []transfer.RestoreFileMetadata) error {
	//TODO Consider extend gRPC API to send client and server IP or external NAME it allows trigger backup from many servers to one client
	client := network.CreateTransferClient()
	log.Debugln("Got restore files meta: ", restoreFilesMetadatas)
	session, err := client.Connect(config.MainConfig.ServerExternalName, config.MainConfig.ServerDataPort)
	if err != nil {
		log.Errorln("Could not connect to server for restore, err: ", err.Error())
		return err
	}
	log.Println("Session ID: ", session.Id)
	err = session.StartRestore(restoreFilesMetadatas)
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
