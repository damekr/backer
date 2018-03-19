package server

import (
	"fmt"
	"github.com/d8x/bftp/common"
	"log"
	"net"
)

func Listen(params *common.ConnParameters) (net.Listener, error) {
	list := fmt.Sprintf("%s:%v", params.Server, params.Port)
	ln, err := net.Listen("tcp", list)
	if err != nil {
		log.Println("Could not listen on port, error: ", err)
		return nil, err
	}
	return ln, nil
}
