package main

import (
	"net/rpc"
	"log"
	"fmt"
    "github.com/backer/common"
)

var result int



func main(){

    client, err := rpc.DialHTTP("tcp", ":8080")
    if err != nil {
        log.Fatalf("Error in dialing. %s", err)
    }
    backup := &common.Backup{Paths: "ALA"}
    fmt.Println("Backup: ", backup)
    err = client.Call("Args.ShowPaths", backup, &result)
    if err != nil {
        log.Fatalf("Error: %s", err)
    }
    log.Printf("Return: %s, %s", result, backup.Paths)
}
