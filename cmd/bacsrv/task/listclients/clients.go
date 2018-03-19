package listclients

import (
	"fmt"

	"github.com/damekr/backer/cmd/bacsrv/db"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "task:listclients"})

//TODO Move it, its internal operation hence should not be here

type ListClients struct {
	Names []string `json:"clientName"`
}

func Create() *ListClients {
	return &ListClients{}
}

func (b *ListClients) Run() {
	log.Println("Getting clients")
	database := db.Get()
	b.Names = database.GetClientsNames()
}

func (b *ListClients) Stop() {
	fmt.Println("Stopping")
}
