package main

import (
	"net/rpc/jsonrpc"
	"net"
	"log"
	"time"
    "github.com/backer/common"
)

type Client struct {
    auth    map[string]string
    address string
    conn    *net.Conn

}


func (c *Client)initConnection(){
    conn, err := net.Dial("tcp", c.address)
    if err != nil {
        log.Fatal("Client is not available")
    }
    c.conn = &conn
}

func (c *Client) Ping(){
    conn := c.conn
    jsonconn := jsonrpc.NewClient(*conn)
    var reply common.Reply
    var args *common.Args
    now := time.Now()
    sec := now.Unix()
    args = &common.Args{A: sec}
    e := jsonconn.Call("Client.Ping", args, &reply)
    if e != nil {
        log.Fatal("Call error: ", e)
    }
    log.Print("Response: ", reply.C)
    if reply.C != sec {
        log.Fatal("Client is not available")
    }
   
}

func (c *Client) RunBackup(){
    conn, err := net.Dial("tcp", c.address)
    if err != nil {
        log.Fatal("Client is not available")
    }
    var reply common.Reply
    var args common.Args
    jsonconn := jsonrpc.NewClient(conn)
    args = common.Args{Path: "LAA"}
    e := jsonconn.Call("Client.ExecuteBackup", args, &reply)
    if e != nil {
        log.Fatal("Call error: ", e)
    }
    log.Print("Response: ", reply.C)
}

func (c *Client) RunRestore(destination string){

}