package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/backer/baclnt/api"
	"os"
	"os/signal"
	"syscall"
)

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func mainLoop() (string, error) {
	log.Debug("Entering into main loop...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)
	go api.ServeServer()
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

func main() {
	log.Info("Starting baclnt application...")
	srv, err := mainLoop()
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
