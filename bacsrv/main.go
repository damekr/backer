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
    backup := &common.Backup{Paths: []string{"ALA","ASL"}}
    fmt.Println("Doing backup of: ", backup)
    err = client.Call("Args.ShowPaths", backup, &result)
    if err != nil {
        log.Fatalf("Error: %s", err)
    }
    log.Printf("Return: %s, %s", result, backup.Paths)
    for _, v := range backup.Paths{
        log.Println(v)
    }
}
