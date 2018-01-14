package network

import (
	"encoding/json"
	"fmt"
	"net"
	"path/filepath"

	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/db"
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
	Sessions       map[uint64]*transfer.MainSession
	Storage        storage.Storage
	CreateMetadata bool
	Database       db.DB
}

func StartTCPDataServer(storage storage.Storage, writeSessionMetadata bool) {
	params := common.NewConnParameters()
	dataNetwork := DataNetwork{
		Params:         params,
		Sessions:       make(map[uint64]*transfer.MainSession),
		Storage:        storage,
		CreateMetadata: writeSessionMetadata,
		Database:       db.Get(),
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
		jsonData, err := json.Marshal(session.Metadata)
		if err != nil {
			log.Errorln("Could not marshal json with metadata")
		} else {
			if err := d.Database.WriteBackupMetadata(jsonData, filepath.Base(session.Metadata.SavesetPath), session.Metadata.ClientName); err != nil {
				log.Println("Could not create metadata, err: ", err.Error())
			}
		}

	}
	delete(d.Sessions, sessionID)
}

func (d DataNetwork) SetPort(port string) {
	d.Params.Port = port
}

func (d DataNetwork) SetIP(ip string) {
	d.Params.Server = ip
}
