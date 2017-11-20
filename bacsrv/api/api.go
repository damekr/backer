package api

import (
	"context"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/job"
	"github.com/damekr/backer/bacsrv/network"
	"github.com/damekr/backer/bacsrv/task/backup"
	"github.com/damekr/backer/bacsrv/task/ping"
	"github.com/damekr/backer/common/protosrv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

// Ping returns hostname of client
func (s *server) Ping(ctx context.Context, in *protosrv.PingRequest) (*protosrv.PingResponse, error) {
	log.Printf("Got request to ping client: %s", in.Ip)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
	clientMessage, err := pingClient(in.Ip)
	if err != nil {
		log.Errorln("Cannot ping client, err: ", err)
	}
	return &protosrv.PingResponse{Message: clientMessage}, nil
}

func (s *server) Backup(ctx context.Context, backupRequest *protosrv.BackupRequest) (*protosrv.BackupResponse, error) {
	log.Printf("Got reqest to backup client: %s", backupRequest.Ip)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
	validatedPaths, err := backupClient(backupRequest.Ip, backupRequest.Paths)
	if err != nil {
		log.Errorln("Cannot backup client, err: ", err)
	}
	log.Printf("Validated paths from client: ", validatedPaths)
	bacsrvBackupResponse := &protosrv.BacsrvBackupResponse{
		Backupstatus: true,
	}
	return &protosrv.BackupResponse{BacsrvBackupResponse: bacsrvBackupResponse}, nil
}

func pingClient(clientIP string) (string, error) {
	log.Println("PINGING CLIENT: ", clientIP)
	pingTask := ping.CreatePing(clientIP)
	pingJob := job.New("ping")
	pingJob.AddTask(pingTask)
	pingJob.Start()

	return pingTask.Message, nil
}

func backupClient(clientIP string, paths []string) ([]string, error) {
	log.Println("Creating backup job of: ", clientIP)
	backupTask := backup.CreateBackup(clientIP, paths)
	backupJob := job.New("backup")
	backupTask.Setup(paths)
	backupJob.AddTask(backupTask)
	backupJob.Start()

	return backupTask.Paths, nil
}

// Start method starts a grpc server on specific port
func Start() error {
	list, err := network.StartTCPMgmtServer()
	if err != nil {
		log.Errorln("Cannot start Mgmt server, err: ", err)
	}
	s := grpc.NewServer()
	protosrv.RegisterBacsrvServer(s, &server{})
	s.Serve(list)
	return nil
}
