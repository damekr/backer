package main

import (
    "fmt"
    // "github.com/backer/baclnt/interface"

	"net/rpc"
	"log"
	"net"
	"net/rpc/jsonrpc"
)

func startInterfaceClient(){
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
    // startInterfaceClient()
    host := "localhost"
    port := 27001
    paths := []string{
        "/tmp",
        "/home/damian/dupa",
    }
    archivename := "tmp.tar"
    connection := TransferConnection{
        Port: port,
        Host: host,
    }
    
    backup := BackupConfig{
        TRConn: connection,
    }
    backup.CreateArchive(paths, archivename)

    
}