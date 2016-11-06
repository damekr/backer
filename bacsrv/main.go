package main

import (
	"os/signal"
	// "fmt"
	"os"
	// "github.com/backer/bacsrv/repository"
	"github.com/backer/bacsrv/api"
	"github.com/backer/bacsrv/transfer"
	"syscall"
	log "github.com/Sirupsen/logrus"
	// "github.com/backer/bacsrv/config"
	"github.com/backer/bacsrv/repository"
	"fmt"
)

func init(){
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func mainLoop() (string, error){
	log.Debug("Entering into main loop...")
    interrupt := make(chan os.Signal, 1)
    signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	startDataServer()
	startApi()
    for {
        select{
            case killSignal := <-interrupt:
            log.Info("Got signal: ", killSignal)
            log.Info("Stopping application, exiting...")
            if killSignal == os.Interrupt {
                return "Application was interrupted by system signal", nil
            }
        return "Application was killed", nil
        }
    
	}
	 
}

func startApi(){
    paths := []string{"/home/damian/test"}
    go api.SendBackupRequest(paths)
}

func startDataServer(){
	// It should have channel communication to close connection after stopping
	go transfer.InitTransferServer() 
}


func main(){
	
	// config.InitClientsConfig()
	clientBucket := repository.CreateClient("minitx")
	fmt.Printf("Client %#v\n", clientBucket)
	// repo, err := repository.CreateRepository()
	// if err != nil{
	// 	fmt.Println("Cannot create repository")
	// }
	// fmt.Printf("Repository status: %#v\n", repo.GetCapacityStatus())
	// srv, err := mainLoop()
	// if err != nil {
	// 	fmt.Println("An error during starting daemon")
	// 	os.Exit(1)
	// }
	// fmt.Println(srv)
	

}