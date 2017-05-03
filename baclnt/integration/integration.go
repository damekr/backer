// +build linux darwin

package integration

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/config"
	"os/exec"
	"strings"
)

func init() {
	log.Debugln("Initializing integration module")
}

func updateMainConfigStructWithCID(cid string) {
	config.ClntConfig.CID = cid
}

// GetCID returns hostid
func GetCID() string {
	var out bytes.Buffer
	hostidcmd := exec.Command("hostid")
	hostidcmd.Stdout = &out
	err := hostidcmd.Run()
	if err != nil {
		log.Error("Cannot get hostid, error: ", err.Error())
		return ""
	}
	hostid := strings.TrimRight(out.String(), "\n")
	updateMainConfigStructWithCID(hostid)
	return hostid
}
