package client

import (
	"fmt"
	"net"
	"time"

	"github.com/damekr/backer/common/protosrv"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "client"})

type Client interface {
	ListAllInSecure() ([]string, error)
	PingInSecure() (string, error)
	RunBackup(string, []string) error
	// ConnectSecure(server string, port string, user string, password string) (*grpc.ClientConn, error)
	// ConnectInSecure(server string, port string) (*grpc.ClientConn, error) //TODO It must handle also RESTApi requests
}

type ClientGRPC struct {
	Server   string
	User     string
	Password string
	Port     string
}

type ClientREST struct {
	Server   string
	User     string
	Password string
	Port     string
}

func (c ClientGRPC) PingInSecure(clientIP string) (string, error) {
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	log.Printf("Sending message to: %s:%s", c.Server, c.Port)
	conn, err := c.ConnectInSecure(c.Server, c.Port)
	if err != nil {
		log.Warningf("Cannot connect to address %s", c.Server)
		return "", err
	}
	defer conn.Close()
	cn := protosrv.NewBacsrvClient(conn)
	//Contact the server and print out its response.

	r, err := cn.Ping(ctx, &protosrv.PingRequest{Ip: clientIP})
	if err != nil {
		log.Warningf("Could not get client name: %v", err)
		return "", err
	}
	log.Debugf("Received client name: %s", r.Message)
	return r.Message, nil
}

func (c ClientREST) PingInSecure() (string, error) {
	log.Info("Pinging from REST.....")
	return "", nil
}

//
//func (c ClientGRPC) ListAllInSecure() ([]string, error) {
//	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
//	ctx := metadata.NewContext(context.Background(), md)
//	log.Printf("Sending message to: %s:%s", c.Server, c.Port)
//	conn, err := c.ConnectInSecure(c.Server, c.Port)
//	if err != nil {
//		log.Warningf("Cannot connect to address %s", c.Server)
//		return nil, err
//	}
//	defer conn.Close()
//
//	cn := proto.NewBacsrvClient(conn)
//	//Contact the server and print out its response.
//	hostname, err := os.Hostname()
//	if err != nil {
//		log.Error("Cannot get hostname setting default")
//		hostname = "client"
//	}
//
//	r, err := cn.ListClients(ctx, &proto.HelloRequest{Name: hostname})
//	if err != nil {
//		log.Warningf("Could not get client name: %v", err)
//		return nil, err
//	}
//
//	log.Debugf("Received client name: %s", r.Clients)
//	return r.Clients, nil
//}
//
//func (c ClientREST) ListAllInSecure() ([]string, error) {
//	log.Info("Listing clients rest....")
//	return nil, nil
//}
//
func (c ClientGRPC) RunBackupInSecure(backupClientIP string, paths []string) error {
	log.Infof("Using GRPC protocol to run backup")
	conn, err := c.ConnectInSecure(c.Server, c.Port)
	if err != nil {
		log.Warningf("Cannot connect to address %s", c.Server)
		return err
	}
	defer conn.Close()

	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	log.Printf("Sending message to: %s:%s", c.Server, c.Port)
	cn := protosrv.NewBacsrvClient(conn)
	//Contact the server and print out its response.
	r, err := cn.Backup(ctx, &protosrv.BackupRequest{
		Ip:    backupClientIP,
		Paths: paths,
	})
	if err != nil {
		log.Warningf("Could not send client name: %v", err)
	}
	log.Debug("Received status of backup: ", r.Backupstatus)

	return nil
}

func (c ClientGRPC) RunRestoreInSecure(restoreClientIP string, paths []string) error {
	log.Infof("Using GRPC protocol to run restore")
	conn, err := c.ConnectInSecure(c.Server, c.Port)
	if err != nil {
		log.Warningf("Cannot connect to address %s", c.Server)
		return err
	}
	defer conn.Close()
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	log.Printf("Sending message to: %s:%s", c.Server, c.Port)
	cn := protosrv.NewBacsrvClient(conn)

	r, err := cn.Restore(ctx, &protosrv.RestoreRequest{
		Ip:    restoreClientIP,
		Paths: paths,
	})
	if err != nil {
		log.Warningf("Could not send client name: %v", err)
	}
	log.Debug("Received status of restore: ", r.Status)

	return nil
}

//func (c ClientREST) RunBackupInSecure([]string, error) error {
//	log.Info("Running fs of clients rest....")
//	return nil
//}

func (c ClientGRPC) ListBackupsInSecure(clientName string) error {
	log.Infof("Using GRPC protocol to list backups")
	conn, err := c.ConnectInSecure(c.Server, c.Port)
	if err != nil {
		log.Warningf("Cannot connect to address %s", c.Server)
		return err
	}
	defer conn.Close()
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx := metadata.NewOutgoingContext(context.Background(), md)
	log.Printf("Sending message to: %s:%s", c.Server, c.Port)
	cn := protosrv.NewBacsrvClient(conn)
	//Contact the server and print out its response.

	r, err := cn.ListBackups(ctx, &protosrv.ListBackupsRequest{
		ClientName: clientName,
	})
	if err != nil {
		log.Warningf("Could not send client name: %v", err)
	}
	fmt.Println(r)

	return nil
}

func (c ClientGRPC) ConnectInSecure(server string, port string) (*grpc.ClientConn, error) {
	log.Debug("Establishing grpc connection with: ", server)
	conn, err := grpc.Dial(net.JoinHostPort(server, port), grpc.WithInsecure())
	if err != nil {
		log.Error("Cannot establish connection with client: ", server)
		return nil, err
	}
	log.Debug("Successfully established grpc connection with client: ", server)
	return conn, nil
}
