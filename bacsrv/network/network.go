package network

import (
	"net"
	"strconv"

	"github.com/d8x/bftp"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "network"})

// EstablishGRPCConnection initiate grpc connection
func EstablishGRPCConnection(clientIP string) (*grpc.ClientConn, error) {
	log.Debugf("Establishing grpc connection with: %s and port: %s ", clientIP, config.MainConfig.ClntMgmtPort)
	conn, err := grpc.Dial(net.JoinHostPort(clientIP, config.MainConfig.ClntMgmtPort), grpc.WithInsecure())
	if err != nil {
		log.Error("Cannot establish connection with client: ", clientIP)
		return nil, err
	}
	return conn, nil
}

func StartTCPDataServer(storageType storage.Storage) {
	server := bftp.CreateBFTPServer(storageType)
	server.SetIP(config.MainConfig.DataTransferInterface)
	port, err := strconv.Atoi(config.MainConfig.DataPort)
	if err != nil {
		log.Errorf("Cannot convert port to int")
	}
	server.SetPort(port)
	err = server.StartServer()
	if err != nil {
		log.Error("Cannot listen TCP Data Server, err: ", err.Error())
	}
	log.Infof("Starting bacsrv data server on addr: %s, port: %i ", config.MainConfig.DataTransferInterface, port)

}

func StartTCPMgmtServer() (net.Listener, error) {
	lis, err := net.Listen("tcp", ":"+config.MainConfig.MgmtPort)
	log.Info("Starting bacsrv protoapi on addr: ", lis.Addr())
	if err != nil {
		return nil, err
	}
	return lis, nil
}
