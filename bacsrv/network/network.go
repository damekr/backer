package network

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"path/filepath"

	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/storage"
	"github.com/damekr/backer/bacsrv/transfer"
	"github.com/damekr/backer/common"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "network"})

var sessionID uint64

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

func StartTCPMgmtServer() (net.Listener, error) {
	lis, err := net.Listen("tcp", ":"+config.MainConfig.MgmtPort)
	log.Info("Starting bacsrv protoapi on addr: ", lis.Addr())
	if err != nil {
		return nil, err
	}
	return lis, nil
}

type DataNetwork struct {
	Params         *common.ConnParameters
	Sessions       map[uint64]*transfer.Session
	Storage        storage.Storage
	CreateMetadata bool
}

func StartTCPDataServer(storage storage.Storage, writeSessionMetadata bool) {
	params := common.NewConnParameters()
	dataNetwork := DataNetwork{
		Params:         params,
		Sessions:       make(map[uint64]*transfer.Session),
		Storage:        storage,
		CreateMetadata: writeSessionMetadata,
	}
	dataNetwork.SetIP(config.MainConfig.DataTransferInterface)
	dataNetwork.SetPort(config.MainConfig.DataPort)
	err := dataNetwork.startServer()
	if err != nil {
		log.Error("Cannot listen TCP Data DataNetwork, err: ", err.Error())
	}
	log.Infof("Starting bacsrv data server on addr: %s, port: %s", config.MainConfig.DataTransferInterface, config.MainConfig.DataPort)

}

func (d DataNetwork) startServer() error {
	list := fmt.Sprintf("%s:%s", d.Params.Server, d.Params.Port)
	ln, err := net.Listen("tcp", list)
	if err != nil {
		log.Println("Could not listen on port, error: ", err)
		return err
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println("Could not accept client connection")
			continue
		}
		c, _ := conn.(*net.TCPConn)
		d.connectionHandler(c)
	}
	return nil
}

func (d DataNetwork) connectionHandler(conn *net.TCPConn) {
	sessionID++
	session := transfer.NewSession(sessionID, d.Params, conn, d.Storage)
	err := session.Negotiate(common.PROTOVERSION)
	if err != nil {
		log.Println("Protocol revision number mismatch", err)
		return
	}
	err = session.Authenticate(common.PASSWORD)
	if err != nil {
		log.Println("Authentication failed")
		return
	}
	d.Sessions[sessionID] = session
	for {
		err = session.SessionDispatcher()
		if err != nil {
			break
		}
	}
	if d.CreateMetadata {
		if err := d.createMetadata(session.Metadata); err != nil {
			log.Println("Could not create metadata, err: ", err.Error())
		}
	}
	delete(d.Sessions, sessionID)
}

func (d DataNetwork) createMetadata(metadata transfer.SessionMetaData) error {
	log.Println("Creating metadata")
	jsonData, err := json.Marshal(metadata)
	if err != nil {
		return err
	}
	file, err := os.Create(filepath.Join(metadata.Saveset, filepath.Base(metadata.Saveset+".json")))
	if err != nil {
		return err
	}
	_, err = file.Write(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func (d DataNetwork) SetPort(port string) {
	d.Params.Port = port
}

func (d DataNetwork) SetIP(ip string) {
	d.Params.Server = ip
}
