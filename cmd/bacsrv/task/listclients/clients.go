package listclients

import (
	"fmt"

	"github.com/damekr/backer/cmd/bacsrv/db"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "task:listclients"})

type ListClients struct {
	Names []string `json:"clientName"`
}

func Create() *ListClients {
	return &ListClients{}
}

func (b *ListClients) Run() {
	log.Println("Getting clients")
	database := db.DB()
	names, err := database.ReadClientsNames()
	if err != nil {
		log.Errorln("Could not read clients names, err: ", err)
	} else {
		b.Names = names
	}

}

func (b *ListClients) Stop() {
	fmt.Println("Stopping")
}
