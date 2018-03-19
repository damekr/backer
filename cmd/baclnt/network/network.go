package network

import (
	"fmt"
	"log"
	"net"
	"os"

	"github.com/damekr/backer/cmd/baclnt/transfer"
	"github.com/damekr/backer/pkg/bftp"
)

var sessionID uint64

type Client struct {
	Params *bftp.ConnParameters
}

func CreateTransferClient() Client {
	params := bftp.NewConnParameters()
	return Client{
		Params: params,
	}
}

func (c Client) Connect(serverIP string, serverPort string) (*transfer.MainSession, error) {
	c.Params.Server = serverIP
	c.Params.Port = serverPort
	conn, err := connectToHost(serverIP, serverPort)
	if err != nil {
		log.Println("Cannot connect, error: ", err)
		return nil, err
	}
	mainSession := transfer.NewSession(sessionID, c.Params, conn)
	err = mainSession.Negotiate(bftp.PROTOVERSION)
	if err != nil {
		mainSession.Conn.Close()
		mainSession.Conn = nil
		return nil, err
	}
	err = mainSession.Authenticate(bftp.PASSWORD)
	if err != nil {
		mainSession.Conn.Close()
		mainSession.Conn = nil
		return nil, err
	}
	sessionID++
	return mainSession, nil
}

func connectToHost(host string, port string) (net.Conn, error) {
	server := fmt.Sprintf("%s:%s", host, port)
	connection, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not connect to", server, port, err)
		return nil, err
	}
	return connection, nil
}
