package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/clientsconfig"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/manager"
	"github.com/damekr/backer/bacsrv/operationshandler"
	"github.com/damekr/backer/bacsrv/repository"
	"github.com/damekr/backer/bacsrv/restapi"
	"github.com/damekr/backer/bacsrv/transfer"
)

var commit string

var configFlag = flag.String("config", "", "Configuration file")

func init() {
	flag.StringVar(configFlag, "c", "", "Configuration file")
}

func setLogger(srvConfig *config.ServerConfig) {
	log.SetFormatter(&log.TextFormatter{})
	switch srvConfig.LogOutput {

	case "STDOUT":
		log.SetOutput(os.Stdout)
	case "SYSLOG":
		//TODO
	}
	if srvConfig.Debug {
		log.SetLevel(log.DebugLevel)
	}
}

func mainLoop(srvConfig *config.ServerConfig) (string, error) {
	log.Debug("Entering into main loop...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	startDataServer(srvConfig)
	startRestApi(srvConfig)
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

func startRestApi(srvConfig *config.ServerConfig) {
	// Starting a new goroutine
	//paths := []string{"/home/damian/test"}
	//go api.SendBackupRequest(paths)
	go restapi.StartServerRestAPI(srvConfig)
}

func startDataServer(srvConfig *config.ServerConfig) {
	// It should have channel communication to close connection after stopping
	// Starging a new goroutine
	go transfer.InitTransferServer(srvConfig)
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

func getConfig(path string) *config.ServerConfig {
	config.SetConfigPath(path)
	srvConfig := config.GetServerConfig()
	return srvConfig

}

func main() {
	log.Printf("COMMIT: %s", commit)
	setFlags()
	srvConfig := getConfig(*configFlag)
	setLogger(srvConfig)
	srvConfig.ShowConfig()
	clientsconfig.InitClientsConfig(srvConfig)
	_, err := repository.CreateRepository()
	if err != nil {
		log.Println("Cannot create repository")
	}
	restore := operationshandler.Restore{
		Saveset: "/opt/bacsrc/ala.100MB",
	}
	restoreMessage := &operationshandler.RestoreTriggerMessage{
		ClientName:    "localhost",
		RestoreConfig: restore,
	}
	err = manager.SendRestoreTriggerMessage(restoreMessage)
	if err != nil {
		log.Error("BUG")
	}
	// mainLoop(srvConfig)
	//fmt.Println("REPO", repo.Location)
	//fmt.Printf("Repository status: %#v\n", repo.GetCapacityStatus())
	//clientBucket := repository.CreateClient("minitx")
	//fmt.Printf("Client %#v\n", clientBucket)

}
