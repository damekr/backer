package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/api"
	"github.com/damekr/backer/baclnt/config"
)

var configFlag = flag.String("config", "", "Configuration file")

func init() {
	flag.StringVar(configFlag, "c", "", "Configuration file")
}

func setLogger(clntConfig *config.ClientConfig) {
	log.SetFormatter(&log.TextFormatter{})
	switch clntConfig.LogOutput {

	case "STDOUT":
		log.SetOutput(os.Stdout)
	case "SYSLOG":
		//TODO
	}
	if clntConfig.Debug {
		log.SetLevel(log.DebugLevel)
	}
}

func mainLoop(clntConfig *config.ClientConfig) (string, error) {
	log.Debug("Entering into main loop...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	startProtoApi(clntConfig)
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

func startProtoApi(config *config.ClientConfig) {
	go api.ServeServer(config)
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
	clntConfig := config.ReadConfigFile(*configFlag)
	clntConfig.ShowConfig()
	log.Info("Starting baclnt application...")
	srv, err := mainLoop(clntConfig)
	if err != nil {
		log.Error("Cannot start client application, error: ", err.Error())
		os.Exit(1)
	}
	log.Info(srv)

	// fmt.Println("OK")
	// // startInterfaceClient()
	// host := "localhost"
	// port := 27001
	// paths := []string{
	//     "/tmp",
	//     "/home/damian/dupa",
	// }
	// archivename := "tmp.tar"
	// connection := TransferConnection{
	//     Port: port,
	//     Host: host,
	// }

	// backup := BackupConfig{
	//     TRConn: connection,
	// }
	// backup.CreateArchive(paths, archivename)

}
