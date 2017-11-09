package test

import (
	"fmt"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/job"
	"github.com/damekr/backer/bacsrv/task/backup"
	"time"
)

func ShowConfig(){
	fmt.Println("TEST MAIN:  ", config.MainConfig.DataPort)
	fmt.Println("TEST CLIENTS: ", config.MainClientsConfig.AllClients)
	fmt.Println("TEST BACKUPS: ", config.MainBackupsConfig.AllBackups)
}



func StartBackup(){
	client := config.Client{
		Name: "SAS",
	}
	//client2 := config.Client{
	//	Name: "221431",
	//}
	prog1 := make(chan int)

	backup1 := backup.CreateBackup(client)
	backup1.SetupBackup("asda")
	job1 := job.New(backup1)
	job1.AddTask(backup1)


	//prog2 := make(chan int)
	//
	//
	//go work.Start(prog1)
	//
	//
	//work2 := job.Create("backup", client2, 12)
	//backupJob2 := work2.(*job.Backup)
	//backupJob2.SetupBackup("asdsad")
	//
	//go work2.Start(prog2)
	//
	for {
		select {
		case pr1 := <-prog1:
			fmt.Println("Prog1: ", pr1)
		default:
			fmt.Println("Default")
		}
		time.Sleep(time.Millisecond * 500)
	}



}