package main

import (
	"flag"
	"fmt"

	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/network"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"

	"os"
	"os/signal"
	"syscall"

	"github.com/damekr/backer/bacsrv/api"
	"github.com/damekr/backer/bacsrv/storage"
)

var commit string
var log = logrus.WithFields(logrus.Fields{"prefix": "main"})

var configFlag = flag.String("config", "", "Configuration file")

func init() {
	flag.StringVar(configFlag, "c", "", "Configuration file")

}

func setLogger() {
	logrus.SetFormatter(&prefixed.TextFormatter{})
	switch config.MainConfig.LogOutput {

	case "STDOUT":
		logrus.SetOutput(os.Stdout)
	case "SYSLOG":
		//TODO
	}
	if config.MainConfig.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}
}

func mainLoop() (string, error) {
	log.Debug("Entering into main loop...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	startDataServer()
	startProtoApi()
	// startRestApi(srvConfig)
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

//
//func startRestApi(srvConfig *config.ServerConfig) {
//	// Starting a new goroutine
//	//paths := []string{"/home/damian/test"}
//	//go api.SendBackupRequest(paths)
//	go restapi.StartServerRestAPI(srvConfig)
//}
//
func startProtoApi() {
	// Starting a new goroutine
	go api.Start()
}

func startDataServer() {
	// It should have channel communication to close connection after stopping
	// Starting a new goroutine
	go network.StartTCPDataServer(storage.DefaultStorage, true)

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

func initStorage(storageType string) {
	err := storage.Create(storageType)
	if err != nil {
		log.Panic("Cannot create storage")
		os.Exit(1)
	}

}

func initConfigs(mainConfigPath string) error {
	err := config.ReadInServerConfig(mainConfigPath)
	if err != nil {
		log.Println("Please setup configuration file")
		return err
	}
	return nil
}

//func test() {
//	bucket := storage.DefaultStorage.CreateBucket("TESTCLIENT_BUCKET")
//	log.Println("BucketLocation: ", bucket.Location)
//	saveset := bucket.CreateSaveset()
//	log.Println("SavesetLocation: ", saveset.Location)
//}

func main() {
	fmt.Println("COMMIT: ", commit)
	setFlags()
	err := initConfigs(*configFlag)
	if err != nil {
		log.Panicln("Cannot init bacsrv configurations")
	}
	setLogger()
	initStorage(config.MainConfig.Storage.Type)
	//test()
	//config.InitClientsConfig(srvConfig)
	//config.InitBackupConfig(srvConfig)
	//initRepository()
	//initClientsBuckets()
	//serverTest
	mainLoop()
	//fmt.Println("REPO", repo.Location)
	//fmt.Printf("Storage status: %#v\n", repo.GetCapacityStatus())
	//clientBucket := storage.CreateClient("minitx")
	//fmt.Printf("clientDefinition %#v\n", clientBucket)

}
