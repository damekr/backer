package grpc

import (
	"context"

	"github.com/damekr/backer/api/protosrv"
	"github.com/damekr/backer/cmd/bacsrv/job"
	"github.com/damekr/backer/cmd/bacsrv/network"
	"github.com/damekr/backer/cmd/bacsrv/task/backup"
	"github.com/damekr/backer/cmd/bacsrv/task/listbackups"
	"github.com/damekr/backer/cmd/bacsrv/task/listclients"
	"github.com/damekr/backer/cmd/bacsrv/task/ping"
	"github.com/damekr/backer/cmd/bacsrv/task/restore"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type server struct{}

var log = logrus.WithFields(logrus.Fields{"prefix": "api"})

// Ping returns hostname of client
func (s *server) Ping(ctx context.Context, in *protosrv.PingRequest) (*protosrv.PingResponse, error) {
	log.WithField("task", "ping").Printf("Got request to ping client: %s", in.Ip)

	s.metadataHandler(ctx)

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

	s.metadataHandler(ctx)

	// Sending gRPC request to start backup (client initialize)
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
	// TODO: Setup here is not needed - task creating handles it
	backupTask.Setup(paths)
	backupJob.AddTask(backupTask)
	backupJob.Start()
	return backupTask.Status, nil
}

// RestoreWholeBackup restores whole backup chosen by backupID to the same location as it was on client
func (s *server) RestoreWholeBackup(ctx context.Context, restoreRequest *protosrv.RestoreRequest) (*protosrv.RestoreResponse, error) {
	log.Printf("Got request to restore client: %s", restoreRequest.Ip)

	s.metadataHandler(ctx)

	// Sending gRPC request to start restore (client initialize)
	err := restoreWholeBackup(restoreRequest.Ip, int(restoreRequest.Backupid))
	if err != nil {
		log.Errorln("Cannot restore client, err: ", err)
	}
	return &protosrv.RestoreResponse{Status: "OK"}, nil
}

func restoreWholeBackup(clientIP string, backupID int) error {
	log.Debugln("Creating restore job of client: ", clientIP)
	log.Debugln("Restore job on backupID: ", backupID)
	restoreOptions := bftp.RestoreOptions{
		WholeBackup: true,
		BasePath:    "",
	}
	restoreTask := restore.Create(clientIP, backupID, restoreOptions)
	err := restoreTask.Setup()
	if err != nil {
		log.Errorln("Error", err)
		return err
	}
	log.Debugln("Restore paths on the server: ", restoreTask.AssetsMetadata)
	restoreJob := job.Create("restore")
	restoreJob.AddTask(restoreTask)
	restoreJob.Start()
	return nil
}

// RestoreWholeBackupDifferentPlace restores whole backup chosen by backupID to different remote location
func (s *server) RestoreWholeBackupDifferentPlace(ctx context.Context, request *protosrv.RestoreWholeBackupDifferentPlaceRequest) (*protosrv.RestoreResponse, error) {
	log.Infoln("Got request to restore client with different remote path: ", request.Restorerequest.Ip)
	log.Debugln("Remote path to restore data: ", request.Remotedir)

	s.metadataHandler(ctx)

	// Sending gRPC request to start restore (client initialize)
	err := restoreWholeBackupDifferentPlace(request.Restorerequest.Ip, request.Remotedir, int(request.Restorerequest.Backupid))
	if err != nil {
		log.Errorln("Cannot restore client, err: ", err)
	}
	return &protosrv.RestoreResponse{Status: "OK"}, nil
}

func restoreWholeBackupDifferentPlace(clientIP, basePath string, backupID int) error {
	log.Debugln("Creating restore job of client: ", clientIP)
	log.Debugln("Restore job on backupID: ", backupID)
	restoreOptions := bftp.RestoreOptions{
		WholeBackup: true,
		BasePath:    basePath,
	}
	restoreTask := restore.Create(clientIP, backupID, restoreOptions)
	err := restoreTask.Setup()
	if err != nil {
		log.Errorln("Error", err)
		return err
	}
	log.Debugln("Restore paths on the server with different location: ", basePath)
	restoreJob := job.Create("restore")
	restoreJob.AddTask(restoreTask)
	restoreJob.Start()

	return nil
}

// RestoreDir restores single directory or file to the same location on client
func (s *server) RestoreDir(ctx context.Context, request *protosrv.RestoreDirRequest) (*protosrv.RestoreResponse, error) {
	log.Infoln("Got request to restore client with different remote path. Client IP:  ", request.Restorerequest.Ip)
	log.Debugln("AbsolutePath to restore data: ", request.ObjectPaths)

	s.metadataHandler(ctx)

	// Sending gRPC request to client to start restore (client initialize restore)
	err := restoreDir(request.Restorerequest.Ip, request.ObjectPaths, int(request.Restorerequest.Backupid))
	if err != nil {
		log.Errorln("Cannot restore client, err: ", err)
	}

	return &protosrv.RestoreResponse{Status: "OK"}, nil
}

func restoreDir(clientIP string, objectsPaths []string, backupID int) error {
	log.Debugln("Creating restore job of client: ", clientIP)
	log.Debugln("Restore job on backupID: ", backupID)
	restoreOptions := bftp.RestoreOptions{
		WholeBackup:  false,
		BasePath:     "",
		ObjectsPaths: objectsPaths,
	}
	restoreTask := restore.Create(clientIP, backupID, restoreOptions)
	err := restoreTask.Setup()
	if err != nil {
		log.Errorln("Error", err)
		return err
	}
	log.Debugf("Restore path: %s to the same location", objectsPaths)
	restoreJob := job.Create("restore")
	restoreJob.AddTask(restoreTask)
	restoreJob.Start()

	return nil
}

func (s *server) RestoreDirRemoteDifferentPlace(ctx context.Context, request *protosrv.RestoreDirRemoteDifferentPlaceRequest) (*protosrv.RestoreResponse, error) {
	log.Infoln("Got request to restore client's dir in different remote path. Client IP:  ", request.Restorerequest.Ip)
	log.Debugln("AbsolutePath  to be restored: ", request.ObjectsPaths)

	s.metadataHandler(ctx)

	err := restoreDirRemoteDifferentPlace(request.Restorerequest.Ip, request.Remotedir, request.ObjectsPaths, int(request.Restorerequest.Backupid))
	if err != nil {
		log.Errorln("Cannot restore client, err: ", err)
	}

	return &protosrv.RestoreResponse{Status: "OK"}, nil
}

func restoreDirRemoteDifferentPlace(clientIP, basePath string, objectsPaths []string, backupID int) error {
	log.Debugln("Creating restore job of client: ", clientIP)
	log.Debugln("Restore job on backupID: ", backupID)
	restoreOptions := bftp.RestoreOptions{
		WholeBackup:  false,
		BasePath:     basePath,
		ObjectsPaths: objectsPaths,
	}

	restoreTask := restore.Create(clientIP, backupID, restoreOptions)
	err := restoreTask.Setup()
	if err != nil {
		log.Errorln("Error", err)
		return err
	}
	log.Debugf("Restore path: %s to different location: %s\n", objectsPaths, basePath)
	restoreJob := job.Create("restore")
	restoreJob.AddTask(restoreTask)
	restoreJob.Start()
	return nil
}

func (s *server) ListBackups(ctx context.Context, listBackupsRequest *protosrv.ListBackupsRequest) (*protosrv.ListBackupsResponse, error) {
	log.Debugln("Got request to list backups of client: ", listBackupsRequest.ClientName)
	s.metadataHandler(ctx)
	if listBackupsRequest.ClientName == "" {
		log.Println("No client given")
		// TODO Return an error? Needs to be tested
		return &protosrv.ListBackupsResponse{
			ClientName: "No clients given",
			BackupID:   []int64{},
		}, nil
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

func (s *server) ListClients(ctx context.Context, listClientsRequest *protosrv.ListClientsRequest) (*protosrv.ListClientsResponse, error) {
	log.Debugln("Got request to list clients")
	s.metadataHandler(ctx)
	clients := listClients()
	return &protosrv.ListClientsResponse{
		Clients: clients,
	}, nil
}

func listClients() []string {
	listClients := listclients.Create()
	listClients.Run()
	return listClients.Names
}

// TODO Just for now, for having it in one place
func (s *server) metadataHandler(ctx context.Context) {
	md, ok := metadata.FromIncomingContext(ctx)
	log.Print("OK: ", ok)
	log.Print("METADATA: ", md)
}

// Start method starts a grpc server on specific port
func Start() error {
	list, err := network.StartTCPManagementServer()
	if err != nil {
		log.Errorln("Cannot start Mgmt server, err: ", err)
	}
	s := grpc.NewServer()
	protosrv.RegisterBacsrvServer(s, &server{})
	s.Serve(list)
	return nil
}
