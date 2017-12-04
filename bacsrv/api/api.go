package api

import (
	"context"

	"github.com/damekr/backer/bacsrv/job"
	"github.com/damekr/backer/bacsrv/network"
	"github.com/damekr/backer/bacsrv/task/backup"
	"github.com/damekr/backer/bacsrv/task/listbackups"
	"github.com/damekr/backer/bacsrv/task/ping"
	"github.com/damekr/backer/common/protosrv"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

var log = logrus.WithFields(logrus.Fields{"prefix": "api"})

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

func pingClient(clientIP string) (string, error) {
	log.Println("PINGING CLIENT: ", clientIP)
	pingTask := ping.CreatePing(clientIP)
	pingJob := job.Create("ping")
	pingJob.AddTask(pingTask)
	pingJob.Start()

	return pingTask.Message, nil
}

func (s *server) Backup(ctx context.Context, backupRequest *protosrv.BackupRequest) (*protosrv.BackupResponse, error) {
	log.Printf("Got request to backup client: %s", backupRequest.Ip)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)

	//Sending gRPC request to start backup (client initialize)
	status, err := backupClient(backupRequest.Ip, backupRequest.Paths)
	if err != nil {
		log.Errorln("Cannot backup client, err: ", err)
	}
	return &protosrv.BackupResponse{Backupstatus: status}, nil
}

func backupClient(clientIP string, paths []string) (bool, error) {
	log.Println("Creating backup job of: ", clientIP)
	backupTask := backup.CreateBackup(clientIP, paths)
	backupJob := job.Create("fs")
	backupTask.Setup(paths)
	backupJob.AddTask(backupTask)
	backupJob.Start()
	return backupTask.Status, nil
}

func (s *server) Restore(ctx context.Context, restoreRequest *protosrv.RestoreRequest) (*protosrv.RestoreResponse, error) {
	log.Printf("Got request to restore client: %s", restoreRequest.Ip)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)

	//Sending gRPC request to start restore (client initialize)
	err := restoreClient(restoreRequest.Ip, restoreRequest.Paths)
	if err != nil {
		log.Errorln("Cannot restore client, err: ", err)
	}
	return &protosrv.RestoreResponse{Status: "OK"}, nil
}

func restoreClient(clientIP string, paths []string) error {
	log.Debugln("Got request to restore client paths: ", paths)
	return nil
}

func (s *server) ListBackups(ctx context.Context, listBackupsRequest *protosrv.ListBackupsRequest) (*protosrv.ListBackupsResponse, error) {
	log.Debugln("Got request to list backups of client: ", listBackupsRequest.ClientName)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
	clientName, backupIds := listBackups(listBackupsRequest.ClientName)
	log.Debugf("Client: %s backups: %x", clientName, backupIds)
	return &protosrv.ListBackupsResponse{
		ClientName: clientName,
		BackupID:   backupIds,
	}, nil
}

func listBackups(clientName string) (string, []int64) {
	listBackups := listbackups.Create(clientName)
	listBackups.Run()
	return listBackups.ClientName, listBackups.BackupIDs
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
