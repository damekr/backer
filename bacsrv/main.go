package main

import (
	"os/signal"
	// "fmt"
	"flag"
	"os"
	// "github.com/backer/bacsrv/repository"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/backer/bacsrv/api"
	"github.com/backer/bacsrv/transfer"
	// "github.com/backer/bacsrv/config"
	"fmt"

	"github.com/backer/bacsrv/repository"
)

var configFlag = flag.String("config", "", "Configuration file")

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	flag.StringVar(configFlag, "c", "", "Configuration file")
}

func mainLoop() (string, error) {
	log.Debug("Entering into main loop...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	startDataServer()
	startApi()
	for {
		select {
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

func startApi() {
	paths := []string{"/home/damian/test"}
	go api.SendBackupRequest(paths)
}

func startDataServer() {
	// It should have channel communication to close connection after stopping
	go transfer.InitTransferServer()
}

func checkConfigFile(configPath string) error {
	// It works for one file, as viper supports directory with given extensions,
	// it shall be extended by this feature
	_, err := os.Stat(configPath)
	if err == nil {
		return nil
	}
	return err
}

func setFlags() {
	flag.Parse()
	if *configFlag == "" {
		log.Error("Please provide config file with proper flag")
		os.Exit(2)
	}
	if checkConfigFile(*configFlag) != nil {
		log.Error("Provided config path is not a file, exiting...")
		os.Exit(3)

	}
	log.Debug("Config file: ", *configFlag)

}

func main() {
	setFlags()
	// config.InitClientsConfig()
	repo, err := repository.CreateRepository()
	if err != nil {
		fmt.Println("Cannot create repository")
	}
	fmt.Println("REPO", repo.Location)
	//fmt.Printf("Repository status: %#v\n", repo.GetCapacityStatus())
	//clientBucket := repository.CreateClient("minitx")
	//fmt.Printf("Client %#v\n", clientBucket)
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
