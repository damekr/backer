package main

import (
	"net/rpc/jsonrpc"
	"net"
	"log"
	"time"
    "github.com/backer/common"
)



func makeConnection(){
    conn, err := net.Dial("tcp", "localhost:8222")

    if err != nil {
        panic(err)
    }
    defer conn.Close()
    c := jsonrpc.NewClient(conn)
    var reply common.Reply
    var args *common.Args
    now := time.Now()
    sec := now.Unix()
    args = &common.Args{sec}
    e := c.Call("Client.Ping", args, &reply)
    if e != nil {
        log.Fatal("Call error: ", e)
    }
    log.Print("Response: ", reply.C)
    if reply.C != sec {
        log.Fatal("Client is not available")
    }
   
}
