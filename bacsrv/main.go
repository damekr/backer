package main

import (
	// "net/rpc"
	// "log"
	"fmt"
	"os"
    // "github.com/backer/common"
	"github.com/takama/daemon"
	bacsrvd "github.com/backer/bacsrv/daemon"
	
)


// func CheckClientAvailability(){
// 	client := api.Client{address: "localhost:8222"}
// 	backupTime := 10
// 	actualTime := 0
// 	for {
// 		client.initConnection()
// 		time.Sleep(1 * time.Second)
// 		client.Ping()
// 		actualTime++
// 		if actualTime == backupTime{
// 			client.initConnection()
// 			client.RunBackup()
// 		}
// 	}
    
//     // client.RunBackup()

// }

// Data engine and Mgmt engine should be started in 
// Seperated gorutines.


const (

    name = "bacsrvd"
    description = "Daemon for Backup & Restore Service"

)



func main(){
	srv, err := daemon.New(name, description)
	if err != nil {
		fmt.Println("An error during starting daemon")
		os.Exit(1)
	}
	service := &bacsrvd.Service{srv}
	status, err := service.Manage()
	if err != nil {
		fmt.Println(status, "Error: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
	// go InitTransferServer()
    // CheckClientAvailability()
	// Config := readConfigFile()
	// Config.showConfig()


}