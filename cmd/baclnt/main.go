package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/damekr/backer/cmd/baclnt/fs"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"

	// "github.com/damekr/backer/baclnt/fs"
	"github.com/damekr/backer/cmd/baclnt/config"
	"github.com/damekr/backer/cmd/baclnt/grpc"
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
	go grpc.Start()
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

func testFileRead() {
	localFS := fs.NewLocalFileSystem()
	// dataBuff := make([]byte, 1024)
	// file, _ := os.Open("/tmp/file.txt")
	// defer file.Close()
	// metadata, err := localFS.ReadFileMetadata("/tmp/file.txt")
	// metadata.FullPath = "/"
	// metadata.Name = "myfile2.txt"
	// err = localFS.CreateFile(metadata)
	// if err != nil {
	// 	log.Errorln("Cannot create file, err: ", err)
	// }
	// writer, err := localFS.WriteFile(metadata)
	// defer writer.Close()
	// if err != nil {
	// 	log.Errorln("Cannot create writer, err: ", err)
	// }
	//
	// written, err := io.Copy(writer, file)
	// if err != nil {
	// 	log.Errorln("Error during copy: ", err)
	// }
	// log.Infoln("Written: ", written)
	//
	// reader, err := localFS.ReadFile("/tmp/myfile.txt")
	// defer reader.Close()
	// data, err := ioutil.ReadAll(reader)
	// if err != nil {
	// 	log.Errorln("Cannot read from file, err: ", err)
	// }
	//
	// fmt.Println(string(data))
	files, err := localFS.ReadBackupObjectsLocations([]string{"/tmp"})
	if err != nil {
		log.Errorln("Error: ", err)
	}
	fmt.Println(files)

}

func main() {
	setFlags()
	log.Debugf("COMMIT: %s", commit)
	config.ReadInConfig(*configFlag)
	setLogger()
	// testFileRead()
	config.MainConfig.ShowConfig()
	log.Info("Starting baclnt application...")
	srv, err := mainLoop()
	if err != nil {
		log.Error("Cannot start client application, error: ", err.Error())
		os.Exit(1)
	}
	log.Info(srv)

}
