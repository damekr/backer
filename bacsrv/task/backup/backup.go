package backup

import (
	"fmt"
)

type Backup struct {
	ClientIP string
	Paths    []string
	Progress int
}

// TODO - maybe make tasks like: backupDefinition, BackupTask, PreBackupTask, PostBackupTask

func CreateBackup(clientIP string, paths []string) *Backup {
	return &Backup{
		ClientIP: clientIP,
		Paths:    paths,
	}
}

func (b *Backup) Run() {
	fmt.Println("Starting backup of client: ", b.ClientIP)

}

func (b *Backup) Stop() {
	fmt.Println("Stopping")
}

func (b *Backup) Setup(path string) {
	fmt.Println("Setup: ", path)
}
