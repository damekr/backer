package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"

	"github.com/damekr/backer/baclnt/transfer"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"

	// "github.com/damekr/backer/baclnt/fs"
	"github.com/damekr/backer/baclnt/api"
	"github.com/damekr/backer/baclnt/config"
)

var configFlag = flag.String("config", "", "Configuration file")
var commit string
var log = logrus.WithFields(logrus.Fields{"prefix": "main"})

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
	startProtoAPI()
	startTransferServer()
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

func startProtoAPI() {
	go api.Start()
}

func startTransferServer() {
	go transfer.StartBFTPServer()
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
	log.Debugf("COMMIT: %s", commit)
	config.ReadInConfig(*configFlag)
	setLogger()
	config.MainConfig.ShowConfig()
	config.MainConfig.ShowConfig()
	log.Info("Starting baclnt application...")
	srv, err := mainLoop()
	if err != nil {
		log.Error("Cannot start client application, error: ", err.Error())
		os.Exit(1)
	}
	log.Info(srv)

}
