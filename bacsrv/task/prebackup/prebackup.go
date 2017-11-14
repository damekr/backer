package prebackup

import (
	"fmt"
	"github.com/damekr/backer/bacsrv/config"
)

type PreBackup struct {
	Client   config.ClientDefinition
	Progress int
}

// TODO - maybe make tasks like: backupDefinition, BackupTask, PreBackupTask, PostBackupTask

func CreatePreBackup(client config.ClientDefinition) *PreBackup {
	return &PreBackup{
		Client: client,
	}
}

func (b *PreBackup) Run() {
	fmt.Println("Starting prebackup of client: ", b.Client)

}

func (b *PreBackup) Stop() {
	fmt.Println("Stopping")
}

func (b *PreBackup) Setup(path string) {
	fmt.Println("Setup: ", path)
}
