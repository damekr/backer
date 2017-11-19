package transfer

import (
	"fmt"
	"log"
	"net"
	"os"
)

var sessionID uint64

func connectToHost(host string, port int) (net.Conn, error) {
	server := fmt.Sprintf("%s:%v", host, port)
	connection, err := net.Dial("tcp", server)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not connect to", server, port, err)
		return nil, err
	}
	return connection, nil
}

// BFTPConnect initialize connection to remote host according to given parameters
func BFTPConnect(params *ConnParameters) (*Session, error) {
	conn, err := connectToHost(params.Server, params.Port)
	if err != nil {
		log.Println("Cannot connect")
		return nil, err
	}
	cSession := NewSession(sessionID, params, conn)
	err = cSession.negotiate(PROTOVERSION)
	if err != nil {
		cSession.Conn.Close()
		cSession.Conn = nil
		return nil, err
	}
	err = cSession.authenticate(PASSWORD)
	if err != nil {
		cSession.Conn.Close()
		cSession.Conn = nil
		return nil, err
	}
	sessionID++
	return cSession, nil
}
