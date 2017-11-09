package backup

//
//type Backup struct {
//	Client config.clientDefinition
//	Progress int
//}
//// TODO - maybe make tasks like: backupDefinition, BackupTask, PreBackupTask, PostBackupTask
//
//func CreateBackup(client config.clientDefinition) *Backup {
//	return &Backup{
//		Client: client,
//	}
//}
//
//
//func (b *Backup) Start(c chan int){
//	fmt.Println("Starging backup of client: ", b.Client)
//	for {
//		b.Progress ++
//		time.Sleep(time.Second * 1)
//		c <- b.Progress
//		if b.Progress == 5 {
//			break
//		}
//	}
//
//}

//func (b *Backup) Stop() {
//	fmt.Println("Stopping")
//}
//
//func (b *Backup) SetupBackup(path string) {
//	fmt.Println("Setup: ", path)
//}
