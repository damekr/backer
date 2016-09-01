package main

import (
	// "net/rpc"
	// "log"
	// "fmt"
    // "github.com/backer/common"
	"time"
)


func CheckClientAvailability(){
	client := Client{address: "localhost:8222"}
	backupTime := 10
	actualTime := 0
	for {
		client.initConnection()
		time.Sleep(1 * time.Second)
		client.PingClient()
		actualTime++
		if actualTime == backupTime{
			client.initConnection()
			client.RunBackup()
		}
	}
    
    // client.RunBackup()

}



func main(){
	go InitTransferServer()
    CheckClientAvailability()
	

}
