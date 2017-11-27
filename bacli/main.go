package main

import (
	"os"
	"runtime"

	"github.com/damekr/backer/bacli/cmds"
	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var commit string
var log = logrus.WithFields(logrus.Fields{"prefix": "main"})

func init() {
	logrus.SetFormatter(&prefixed.TextFormatter{})
	logrus.SetLevel(logrus.DebugLevel)

	cmds.Execute()
}

func main() {
	log.Debugf("main %#v", os.Args)
	log.Debugf("bacli, compiled with %v on %v/%v", runtime.Version(), runtime.GOOS, runtime.GOARCH)

}
