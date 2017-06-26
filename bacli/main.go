package main

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacli/cmds"
	"os"
	"runtime"
)

var commit string

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	cmds.Execute()
}

func main() {

	log.Debugf("main %#v", os.Args)
	log.Debugf("bacli, compiled with %v on %v/%v", runtime.Version(), runtime.GOOS, runtime.GOARCH)

}
