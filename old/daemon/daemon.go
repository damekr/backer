package daemon

import (
	"fmt"
	"github.com/damekr/backer/bacsrv/api"
	"github.com/takama/daemon"
	"log"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	daemon.Daemon
}

var stdlog, errlog *log.Logger

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: bacsrvd start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "start":
			service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	for {
		select {
		case killSignal := <-interrupt:
			fmt.Println("Got signal: ", killSignal)
			fmt.Println("Stoping daemon...")
			if killSignal == os.Interrupt {
				return "Daemon was interrupted by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

}

func (service *Service) Start() {
	paths := []string{"/tmp", "/var"}
	go api.SendBackupRequest(paths)
}

func (service *Service) Stop() (string, error) {
	return "nil", nil
}

func (service *Service) Status() (string, error) {
	return "nil", nil
}
