package api

import (
	"context"
	"net"

	"github.com/d8x/bftp"
	"github.com/damekr/backer/baclnt/config"
	"github.com/damekr/backer/baclnt/fs"
	"github.com/damekr/backer/common/proto"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

var log = logrus.WithFields(logrus.Fields{"prefix": "api"})

// Ping returns hostname of client
func (s *server) Ping(ctx context.Context, in *proto.PingRequest) (*proto.PingResponse, error) {
	log.Printf("Got request to ping client: %s", in.Ip)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)

	return &proto.PingResponse{Message: "FROM CLIENT"}, nil
}

func (s *server) Backup(ctx context.Context, backupRequest *proto.BackupRequest) (*proto.BackupResponse, error) {
	log.Printf("Got request to fs client: %s", backupRequest.Ip)
	log.Println("Paths to be validated: ", backupRequest.Paths)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
	fileSystem := fs.FS{}
	validatedPaths := fileSystem.GetAbsolutePaths(backupRequest.Paths)
	log.Printf("Validated paths: ", validatedPaths)
	baclntBackupResponse := &proto.BaclntBackupResponse{
		Validpaths: validatedPaths,
	}
	//TODO backupRequest.IP is probably client ip, this message should contains server external ip
	err := runBackup(validatedPaths, backupRequest.Ip)
	if err != nil {
		log.Errorf("Backup Failed, err: ", err.Error())
	}
	return &proto.BackupResponse{BaclntBackupResponse: baclntBackupResponse}, nil
}

func runBackup(paths []string, serverIp string) error {
	client := bftp.CreateBFTPClient()
	session, err := client.Connect(serverIp, 8000)
	if err != nil {
		log.Errorln("Cannot initialize connection")
	}
	log.Println("Session ID: ", session.Id)
	for _, path := range paths {
		if err := session.PutFile(path, path); err != nil {
			log.Error("Cannot send file")
		}
	}
	return nil
}

// Start method starts a grpc server on specific port
func Start() error {
	lis, err := net.Listen("tcp", ":"+config.MainConfig.MgmtPort)
	log.Info("Starting bacsrv protoapi on addr: ", lis.Addr())
	if err != nil {
		log.Errorf("Failed to listen: %v", err)
		return err
	}
	s := grpc.NewServer()
	proto.RegisterBacsrvServer(s, &server{})
	s.Serve(lis)
	return nil
}
