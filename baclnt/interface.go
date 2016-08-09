package main

import (
    "log"
    "net"
    "net/http"
    "net/rpc"
    "github.com/backer/common"
)


func serveInterface(){
    backup := new(common.Args)
    err := rpc.Register(backup)
    if err != nil{
        log.Fatalf("ERROR %s", err)
    }
    rpc.HandleHTTP()
    l, e := net.Listen("tcp", ":8080")
    if e != nil {
        log.Fatalf("Error tcp: %s", e)
    }
    log.Println("Serving RPC Handler")
    err = http.Serve(l, nil)
    if err != nil {
        log.Fatalf("Error serving: %s", err)
    }
}