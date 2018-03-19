package grpc

import (
	"context"

	"github.com/damekr/backer/api/protosrv"
	"github.com/damekr/backer/cmd/bacsrv/job"
	"github.com/damekr/backer/cmd/bacsrv/network"
	"github.com/damekr/backer/cmd/bacsrv/task/backup"
	"github.com/damekr/backer/cmd/bacsrv/task/listbackups"
	"github.com/damekr/backer/cmd/bacsrv/task/ping"
	"github.com/damekr/backer/cmd/bacsrv/task/restore"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

var log = logrus.WithFields(logrus.Fields{"prefix": "api"})

// Ping returns hostname of client
func (s *server) Ping(ctx context.Context, in *protosrv.PingRequest) (*protosrv.PingResponse, error) {
	log.WithField("task", "ping").Printf("Got request to ping client: %s", in.Ip)
	md, ok := metadata.FromIncomingContext(ctx)
	log.WithField("task", "ping").Print("OK: ", ok)
	log.WithField("task", "ping").Print("METADATA: ", md)
	if in.Ip == "" {
		return &protosrv.PingResponse{Message: "OK FROM SERVER"}, nil
	}
	clientMessage, err := pingClient(in.Ip)
	if err != nil {
		log.WithField("task", "ping").Errorln("Cannot ping client, err: ", err)
	}
	return &protosrv.PingResponse{Message: clientMessage}, nil
}

func pingClient(clientIP string) (string, error) {
	log.WithField("task", "ping").Println("PINGING CLIENT: ", clientIP)
	pingTask := ping.CreatePing(clientIP)
	pingJob := job.Create("ping")
	pingJob.AddTask(pingTask)
	pingJob.Start()

	return pingTask.Message, nil
}

func (s *server) Backup(ctx context.Context, backupRequest *protosrv.BackupRequest) (*protosrv.BackupResponse, error) {
	log.WithField("task", "backup").Printf("Got request to backup client: %s", backupRequest.Ip)
	md, ok := metadata.FromIncomingContext(ctx)
	log.WithField("task", "backup").Print("OK: ", ok)
	log.WithField("task", "backup").Print("METADATA: ", md)

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
	backupJob := job.Create("backup")
	//TODO: Setup here is not needed - task creating handles it
	backupTask.Setup(paths)
	backupJob.AddTask(backupTask)
	backupJob.Start()
	return backupTask.Status, nil
}

//RestoreWholeBackup restores whole backup chosen by backupID to the same location as it was on client
func (s *server) RestoreWholeBackup(ctx context.Context, restoreRequest *protosrv.RestoreRequest) (*protosrv.RestoreResponse, error) {
	log.Printf("Got request to restore client: %s", restoreRequest.Ip)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)

	//Sending gRPC request to start restore (client initialize)
	err := restoreWholeBackup(restoreRequest.Ip, int(restoreRequest.Backupid))
	if err != nil {
		log.Errorln("Cannot restore client, err: ", err)
	}
	return &protosrv.RestoreResponse{Status: "OK"}, nil
}

func restoreWholeBackup(clientIP string, backupID int) error {
	log.Debugln("Creating restore job of client: ", clientIP)
	log.Debugln("Restore job on backupID: ", backupID)
	restoreTask := restore.Create(clientIP, backupID)
	err := restoreTask.Setup("", "")
	if err != nil {
		log.Errorln("Error", err)
		return err
	}
	log.Debugln("Restore paths on the server: ", restoreTask.FilesMetadata)
	restoreJob := job.Create("restore")
	restoreJob.AddTask(restoreTask)
	restoreJob.Start()
	return nil
}

//RestoreWholeBackupDifferentPlace restores whole backup chosen by backupID to different remote location
func (s *server) RestoreWholeBackupDifferentPlace(ctx context.Context, request *protosrv.RestoreWholeBackupDifferentPlaceRequest) (*protosrv.RestoreResponse, error) {
	log.Infoln("Got request to restore client with different remote path: ", request.Restorerequest.Ip)
	log.Debugln("Remote path to restore data: ", request.Remotedir)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)

	//Sending gRPC request to start restore (client initialize)
	err := restoreWholeBackupDifferentPlace(request.Restorerequest.Ip, request.Remotedir, int(request.Restorerequest.Backupid))
	if err != nil {
		log.Errorln("Cannot restore client, err: ", err)
	}
	return &protosrv.RestoreResponse{Status: "OK"}, nil
}

func restoreWholeBackupDifferentPlace(clientIP, remotePath string, backupID int) error {
	log.Debugln("Creating restore job of client: ", clientIP)
	log.Debugln("Restore job on backupID: ", backupID)
	restoreTask := restore.Create(clientIP, backupID)
	err := restoreTask.Setup(remotePath, "")
	if err != nil {
		log.Errorln("Error", err)
		return err
	}
	log.Debugln("Restore paths on the server with different location: ", remotePath)
	restoreJob := job.Create("restore")
	restoreJob.AddTask(restoreTask)
	restoreJob.Start()

	return nil
}

//RestoreDir restores single directory or file to the same location on client
func (s *server) RestoreDir(ctx context.Context, request *protosrv.RestoreDirRequest) (*protosrv.RestoreResponse, error) {
	log.Infoln("Got request to restore client with different remote path: ", request.Restorerequest.Ip)
	log.Debugln("Path to restore data: ", request.Dir)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)

	//Sending gRPC request to start restore (client initialize)
	err := restoreDir(request.Restorerequest.Ip, request.Dir, int(request.Restorerequest.Backupid))
	if err != nil {
		log.Errorln("Cannot restore client, err: ", err)
	}

	return &protosrv.RestoreResponse{Status: "OK"}, nil
}

func restoreDir(clientIP, dirPath string, backupID int) error {
	log.Debugln("Creating restore job of client: ", clientIP)
	log.Debugln("Restore job on backupID: ", backupID)
	restoreTask := restore.Create(clientIP, backupID)
	err := restoreTask.Setup("", dirPath)
	if err != nil {
		log.Errorln("Error", err)
		return err
	}
	log.Debugf("Restore path: %s to the same location", dirPath)
	restoreJob := job.Create("restore")
	restoreJob.AddTask(restoreTask)
	restoreJob.Start()

	return nil
}

func (s *server) RestoreDirRemoteDifferentPlace(ctx context.Context, request *protosrv.RestoreDirRemoteDifferentPlaceRequest) (*protosrv.RestoreResponse, error) {

	return &protosrv.RestoreResponse{Status: "OK"}, nil
}

func (s *server) ListBackups(ctx context.Context, listBackupsRequest *protosrv.ListBackupsRequest) (*protosrv.ListBackupsResponse, error) {
	log.Debugln("Got request to list backups of client: ", listBackupsRequest.ClientName)
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
	if listBackupsRequest.ClientName == "" {
		log.Println("No client given, listing all available backups")
	}
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
