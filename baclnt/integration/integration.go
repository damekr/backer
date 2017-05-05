// +build linux darwin

package integration

import (
	"bytes"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/baclnt/config"
	"os/exec"
	"runtime"
	"strings"
)

type ClientInfo struct {
	Name     string `json:"name"`
	CID      string `json:"cid"`
	Platform string `json:"platform"`
}

func init() {
	log.Debugln("Initializing integration module")
}

func updateMainConfigStructWithCID(cid string) {
	config.ClntConfig.CID = cid
}

func GetClientInfo() *ClientInfo {
	cid, err := readCID()
	if err != nil {
		log.Error("Unabled to read CID, setting value 0")
		cid = "0"
	}
	return &ClientInfo{
		Name:     config.GetExternalName(),
		CID:      cid,
		Platform: runtime.GOOS,
	}

}

func readCID() (string, error) {
	var out bytes.Buffer
	hostidcmd := exec.Command("hostid")
	hostidcmd.Stdout = &out
	err := hostidcmd.Run()
	if err != nil {
		log.Error("Cannot get hostid, error: ", err.Error())
		return "", err
	}
	hostid := strings.TrimRight(out.String(), "\n")
	updateMainConfigStructWithCID(hostid)
	return hostid, nil
}
