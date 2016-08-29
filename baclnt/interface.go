package main

import (
    "log"
)

// Args received by connection
type Args struct {
    A int64
}

// Reply returned value to server
type Reply struct {
    C int64
}
// Client just simply data type
type Client int

// Ping method returns recieved value to show that the client is available
func (c *Client) Ping(args *Args, reply *Reply) error {
    log.Println("Values:", args.A)
    reply.C = args.A
    return nil
}

func (c *Client) Error(args *Args, reply *Reply) error {
    panic("ERROR")
}