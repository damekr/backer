package task

import (
	log "github.com/Sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
)

type Task interface {
	Run()
	Stop()
}

const (
	clntMgmtPort = "9090"
	//timestampFormat = time.StampNano
)

func establishConnection(clientIP string) (*grpc.ClientConn, error) {
	log.Debug("Establishing grpc connection with: ", clientIP)
	conn, err := grpc.Dial(net.JoinHostPort(clientIP, clntMgmtPort), grpc.WithInsecure())
	if err != nil {
		log.Error("Cannot establish connection with client: ", clientIP)
		return nil, err
	}
	log.Debug("Successfully established grpc connection with client: ", clientIP)

	return conn, nil
}
