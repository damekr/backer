package test

import (
	"fmt"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/job"
	"github.com/damekr/backer/bacsrv/task/backup"
	"github.com/damekr/backer/bacsrv/task/prebackup"
)

//func ShowConfig() {
//	fmt.Println("TEST MAIN:  ", config.MainConfig.DataPort)
//	fmt.Println("TEST CLIENTS: ", config.MainClientsConfig.AllClients)
//	fmt.Println("TEST BACKUPS: ", config.MainBackupsConfig.AllBackups)
//}

func ShowClientsConfig() {
	fmt.Println("TEST CLIENTS: ", config.AllClients)
}

func StartBackup() {
	client := config.ClientDefinition{
		Name: "SAS",
	}
	//client2 := config.clientDefinition{
	//	Name: "221431",
	//}
	preb := prebackup.CreatePreBackup(client)
	backup1 := backup.CreateBackup(client)
	backup1.Setup("asda")
	job1 := job.New("backup")
	job1.AddTask(preb)
	job1.AddTask(backup1)
	job1.Start()

	//work2 := job.Create("backup", client2, 12)
	//backupJob2 := work2.(*job.backupDefinition)
	//backupJob2.Setup("asdsad")
	//
	//go work2.Start(prog2)

}
