package client

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/common/protosrv"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"net"
	"os"
	"time"
)

func init() {
	log.Debug("Initializing bacli client")
}

type Client interface {
	ListAllInSecure() ([]string, error)
	PingInSecure() (string, error)
	RunBackup([]string) error
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

func (c ClientGRPC) PingInSecure() (string, error) {
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx := metadata.NewContext(context.Background(), md)
	log.Printf("Sending message to: %s:%s", c.Server, c.Port)
	conn, err := c.ConnectInSecure(c.Server, c.Port)
	if err != nil {
		log.Warningf("Cannot connect to address %s", c.Server)
		return "", err
	}
	defer conn.Close()
	cn := protosrv.NewBacsrvClient(conn)
	//Contact the server and print out its response.
	hostname, err := os.Hostname()
	if err != nil {
		log.Error("Cannot get hostname setting default")
		hostname = "client"
	}
	r, err := cn.SayHello(ctx, &protosrv.HelloRequest{Name: hostname})
	if err != nil {
		log.Warningf("Could not get client name: %v", err)
		return "", err
	}
	log.Debugf("Received client name: %s", r.Name)
	return r.Name, nil
}

func (c ClientREST) PingInSecure() (string, error) {
	log.Info("Pinging from REST.....")
	return "", nil
}

func (c ClientGRPC) ListAllInSecure() ([]string, error) {
	md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
	ctx := metadata.NewContext(context.Background(), md)
	log.Printf("Sending message to: %s:%s", c.Server, c.Port)
	conn, err := c.ConnectInSecure(c.Server, c.Port)
	if err != nil {
		log.Warningf("Cannot connect to address %s", c.Server)
		return nil, err
	}
	defer conn.Close()

	cn := protosrv.NewBacsrvClient(conn)
	//Contact the server and print out its response.
	hostname, err := os.Hostname()
	if err != nil {
		log.Error("Cannot get hostname setting default")
		hostname = "client"
	}

	r, err := cn.ListClients(ctx, &protosrv.HelloRequest{Name: hostname})
	if err != nil {
		log.Warningf("Could not get client name: %v", err)
		return nil, err
	}

	log.Debugf("Received client name: %s", r.Clients)
	return r.Clients, nil
}

func (c ClientREST) ListAllInSecure() ([]string, error) {
	log.Info("Listing clients rest....")
	return nil, nil
}

func (c ClientGRPC) RunBackupInSecure(clients []string) error {
	log.Infof("Using GRPC protocol to run backup")
	conn, err := c.ConnectInSecure(c.Server, c.Port)
	if err != nil {
		log.Warningf("Cannot connect to address %s", c.Server)
		return err
	}
	defer conn.Close()
	for k, v := range clients {
		log.Infof("Starting backup of client: %s number: %d", v, k)
		md := metadata.Pairs("timestamp", time.Now().Format(time.StampNano))
		ctx := metadata.NewContext(context.Background(), md)
		log.Printf("Sending message to: %s:%s", c.Server, c.Port)
		cn := protosrv.NewBacsrvClient(conn)
		//Contact the server and print out its response.
		hostname, err := os.Hostname()
		if err != nil {
			log.Error("Cannot get hostname setting default")
			hostname = "client"
		}
		r, err := cn.RunBackup(ctx, &protosrv.Client{
			Name:  hostname,
			Cname: v,
		})
		if err != nil {
			log.Warningf("Could not get client name: %v", err)
		}
		log.Debug("Received status of backup: ", r.Backup)
	}
	return nil
}

func (c ClientREST) RunBackupInSecure([]string, error) error {
	log.Info("Running backup of clients rest....")
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
