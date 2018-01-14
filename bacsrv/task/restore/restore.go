package restore

import (
	"context"
	"fmt"

	"github.com/damekr/backer/bacsrv/network"
	"github.com/damekr/backer/common/protoclnt"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "task:restore"})

type Restore struct {
	ClientIP       string   `json:"clientIP"`
	RequestedPaths []string `json:"requestedPaths"`
	ValidPaths     []string `json:"validPaths"`
	Progress       int      `json:"-"`
	Status         bool     `json:"status"`
	BucketLocation string   `json:"bucketLocation"`
}

func Create(clientIP string, paths []string) *Restore {
	return &Restore{
		ClientIP:       clientIP,
		RequestedPaths: paths,
	}
}

func (r *Restore) Run() {
	log.Println("Running backup of client client: ", r.ClientIP)
	conn, err := network.EstablishGRPCConnection(r.ClientIP)
	if err != nil {
		log.Errorf("Cannot connect to address %s", r.ClientIP)
		return
	}
	defer conn.Close()
	c := protoclnt.NewBaclntClient(conn)
	response, err := c.Restore(context.Background(), &protoclnt.RestoreRequest{Ip: r.ClientIP, Paths: r.RequestedPaths})
	if err != nil {
		log.Warningf("Could not get paths of client: %v", err)
		r.Status = false
		return
	}
	r.Status = true
	log.Println("Response: ", response)
}

func (r *Restore) Stop() {
	fmt.Println("Stopping")
}

func (r *Restore) Setup(paths []string) {
	r.RequestedPaths = paths
}
