package backup

import (
	"fmt"
)

type Backup struct {
	ClientIP       string
	Paths          []string
	Progress       int
	ValidatedPaths []string
}

// TODO - maybe make tasks like: backupDefinition, BackupTask, PreBackupTask, PostBackupTask

func CreateBackup(clientIP string, paths []string) *Backup {
	return &Backup{
		ClientIP: clientIP,
		Paths:    paths,
	}
}

func (b *Backup) Run() {
	//fmt.Println("Starting backup of client: ", b.ClientIP)
	//log.Println("Pinging client: ", b.ClientIP)
	//conn, err := network.EstablishGRPCConnection(b.ClientIP)
	//if err != nil {
	//	log.Warningf("Cannot connect to address %s", b.ClientIP)
	//
	//}
	//defer conn.Close()
	//c := protosrv.NewBacsrvClient(conn)
	//r, err := c.Backup(context.Background(), &protosrv.BackupRequest{Paths: b.Paths})
	//if err != nil {
	//	log.Warningf("Could not get client name: %v", err)
	//}
	//log.Debugf("Received client validated paths: %s", r.Validpaths)
	b.ValidatedPaths = []string{"/etc", "/var"}

}

func (b *Backup) Stop() {
	fmt.Println("Stopping")
}

func (b *Backup) Setup(paths []string) {
	fmt.Println("Setup: ", paths)
}
