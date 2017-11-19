package network

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"google.golang.org/grpc"
	"net"
)

// EstablishGRPCConnection initiate grpc connection
func EstablishGRPCConnection(clientIP string) (*grpc.ClientConn, error) {
	log.Debug("Establishing grpc connection with: ", clientIP)
	conn, err := grpc.Dial(net.JoinHostPort(clientIP, config.MainConfig.ClntMgmtPort), grpc.WithInsecure())
	if err != nil {
		log.Error("Cannot establish connection with client: ", clientIP)
		return nil, err
	}
	log.Debug("Successfully established grpc connection with client: ", clientIP)

	return conn, nil
}

func StartTCPMgmtServer() (net.Listener, error) {
	lis, err := net.Listen("tcp", ":"+config.MainConfig.MgmtPort)
	log.Info("Starting bacsrv protoapi on addr: ", lis.Addr())
	if err != nil {
		return nil, err
	}
	return lis, nil
}
