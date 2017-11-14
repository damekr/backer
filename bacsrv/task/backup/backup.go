package backup

import (
	"fmt"
	"github.com/damekr/backer/bacsrv/config"
)

type Backup struct {
	Client   config.ClientDefinition
	Progress int
}

// TODO - maybe make tasks like: backupDefinition, BackupTask, PreBackupTask, PostBackupTask

func CreateBackup(client config.ClientDefinition) *Backup {
	return &Backup{
		Client: client,
	}
}

func (b *Backup) Run() {
	fmt.Println("Starting backup of client: ", b.Client)

}

func (b *Backup) Stop() {
	fmt.Println("Stopping")
}

func (b *Backup) Setup(path string) {
	fmt.Println("Setup: ", path)
}
