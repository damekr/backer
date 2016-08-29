package main

import (
    "fmt"
    // "github.com/backer/baclnt/interface"

	"net/rpc"
	"log"
	"net"
	"net/rpc/jsonrpc"
)

func startInterfaceServer(){
    client := new(Client)
    server := rpc.NewServer()
    server.Register(client)
    l, e := net.Listen("tcp", ":8222")
    if e != nil {
        log.Fatal("Listen error: ", e)
    }
    for {
        conn, err := l.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go server.ServeCodec(jsonrpc.NewServerCodec(conn))
    }
}


func main(){
    fmt.Println("OK")
    startInterfaceServer()
    // TODO Change NewArchive method to get name and make archive according to given name
    // archive := NewArchive([]string{"/tmp"},"Archiwum")
    // archive.MakeArchive("/home/damian/tmp.tar")

    // InitConnection()
    
}