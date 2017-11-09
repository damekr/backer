package main

import (
	"flag"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/test"
	"os"

)

var commit string

var configFlag = flag.String("config", "", "Configuration file")

func init() {
	flag.StringVar(configFlag, "c", "", "Configuration file")
}
//
//func setLogger(srvConfig *config.ServerConfig) {
//	log.SetFormatter(&log.TextFormatter{})
//	switch srvConfig.LogOutput {
//
//	case "STDOUT":
//		log.SetOutput(os.Stdout)
//	case "SYSLOG":
//		//TODO
//	}
//	if srvConfig.Debug {
//		log.SetLevel(log.DebugLevel)
//	}
//}
//
//func mainLoop(srvConfig *config.ServerConfig) (string, error) {
//	log.Debug("Entering into main loop...")
//	interrupt := make(chan os.Signal, 1)
//	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
//	startDataServer(srvConfig)
//	startProtoApi(srvConfig)
//	// startRestApi(srvConfig)
//	serverTest()
//	for {
//		select {
//		case killSignal := <-interrupt:
//			log.Info("Got signal: ", killSignal)
//			log.Info("Stopping application, exiting...")
//			if killSignal == os.Interrupt {
//				return "Application was interrupted by system signal", nil
//			}
//			return "Application was killed", nil
//		}
//
//	}
//}
//
//func startRestApi(srvConfig *config.ServerConfig) {
//	// Starting a new goroutine
//	//paths := []string{"/home/damian/test"}
//	//go api.SendBackupRequest(paths)
//	go restapi.StartServerRestAPI(srvConfig)
//}
//
//func startProtoApi(srvConfig *config.ServerConfig) {
//	// Starting a new goroutine
//	go inprotoapi.ServeServer(srvConfig)
//}
//
//func startDataServer(srvConfig *config.ServerConfig) {
//	// It should have channel communication to close connection after stopping
//	// Starging a new goroutine
//	go transfer.StartTransferServer(srvConfig)
//}

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


//
//func initRepository() {
//	err := storage.InitRepository()
//	if err != nil {
//		log.Panic("Cannot create storage")
//		os.Exit(1)
//	}
//}

//func initClientsBuckets() {
//	err := storage.InitClientsBuckets()
//	if err != nil {
//		log.Panic("Cannot initialize clients buckets")
//		os.Exit(2)
//	}
//}
//
//func serverTest() {
//
//}

func initConfigs(mainConfigPath string) error{
	err := config.ReadInServerConfig(mainConfigPath)
	if err != nil {
		log.Println("Please setup configuration file")
		return err
	}
	log.Println("Clients config file path: ", config.MainConfig.ClientsConfigFilePath)
	err = config.ReadInClientsConfig(config.MainConfig.ClientsConfigFilePath, config.MainConfig.BackupsConfigFilePath,
		config.MainConfig.SchedulesConfigFilePath)
	if err != nil {
		log.Println("Could not read clients config, error: ",  err)
		return err
	}
	return nil
}


func main() {
	log.Printf("COMMIT: %s", commit)
	setFlags()
	err := initConfigs(*configFlag)
	if err != nil {
		log.Panicln("Cannot init bacsrv configurations")
	}
	test.ShowConfig()
	test.StartBackup()
	//config.InitClientsConfig(srvConfig)
	//config.InitBackupConfig(srvConfig)
	//initRepository()
	//initClientsBuckets()
	//serverTest()
	//mainLoop(srvConfig)

	//fmt.Println("REPO", repo.Location)
	//fmt.Printf("Storage status: %#v\n", repo.GetCapacityStatus())
	//clientBucket := storage.CreateClient("minitx")
	//fmt.Printf("Client %#v\n", clientBucket)

}
