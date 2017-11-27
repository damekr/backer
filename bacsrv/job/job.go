package job

import (
	"github.com/damekr/backer/bacsrv/task"
)

type Job struct {
	Tasks []task.Task
	ID    int
	Name  string
}

var Jobs []*Job
var id = 0

// TODO Now can be many jobs, but must be possible to define many tasks for one particular job
var jobs2 map[*Job][]*task.Task

func New(name string) *Job {
	id++
	return &Job{
		ID:   id,
		Name: name,
	}
}

func (j *Job) AddTask(task task.Task) error {
	j.Tasks = append(j.Tasks, task)
	//switch task {
	//case task.(*fs.Backup):
	//	j.Tasks = append(j.Tasks, task)
	//
	//}
	return nil
}

func (j *Job) Start() error {
	for _, t := range j.Tasks {
		t.Run()
	}
	return nil

}

func GetAllTasks() []*Job {
	return Jobs
}

//type BackupJob struct {
//	BackupConfig	*config.Backup
//	ClientConfig	*config.Client
//}

//func (b *BackupJob) Run() error{
//	log.Info("Starting fs job of client: ", b.ClientConfig.Address)
//	validatedPaths, err := preBackupChecks(b.BackupConfig.Paths, b.ClientConfig.Address)
//	if err != nil {
//		log.Error("Cannot validate paths on client side")
//		return err
//	}
//	// TODO Here paths can be removed bases on excluded
//	log.Debugf("Got validated paths from client: %s starting fs...", validatedPaths)
//	err = outprotoapi.SendBackupRequest(validatedPaths, b.ClientConfig.Address)
//	if err != nil {
//		log.Error("Triggering fs failed!")
//		return err
//	}
//	log.Info("Backup has been triggered properly!")
//	return nil
//}
//
//
//func preBackupChecks(paths []string, clntAddr string) ([]string, error) {
//	log.Debug("Starting executing prebackup checks...")
//	log.Debug("Checking if client responds")
//	hostname, err := SendHelloMessageToClient(clntAddr)
//	if err != nil {
//		log.Errorf("ClientConfig %s does not responds", clntAddr)
//		return nil, err
//	}
//	log.Debug("ClientConfig sent it's own hostname, and it is: ", hostname)
//	checkedPaths, err := outprotoapi.CheckPaths(clntAddr, paths)
//	if err != nil {
//		log.Error("An error ocurred during checking paths, error: ", err.Error())
//		return nil, err
//	}
//	return checkedPaths, nil
//}
//
//// SendHelloMessageToClient is responsible for proxing restapi reqests to clients
//func SendHelloMessageToClient(clntAddress string) (string, error) {
//	clntHostname, err := outprotoapi.SayHelloToClient(clntAddress)
//	if err != nil {
//		log.Errorf("Given client on address %s is not available", clntAddress)
//		return "", err
//	}
//	return clntHostname, nil
//
//}
//
//// IntegrateClient performs client integration with all operatinos
//func IntegrateClient(client *config.Client) error {
//	log.Infof("Starting %s integration...", client.Name)
//	clntHostname, err := SendHelloMessageToClient(client.Address)
//	if err != nil {
//		log.Errorf("ClientConfig %s with address: %s does not respond", client.Name, client.Address)
//		return err
//	}
//	log.Debugf("Got hostname: %s from client side, performing integration", clntHostname)
//	remoteInformations, err := outprotoapi.SendIntegrationRequest(client)
//	if err != nil {
//		log.Error("Cannot get information from remote host: ", client.Address)
//	}
//	log.Debugf("Got information %#v about client", remoteInformations)
//	return nil
//}
//
////// StartBackup start fs on client with given configuration
////// This function should require only BackupJob Struct
////func StartBackup(backupConfig *config.Backup, clntAddr string) error {
////	log.Info("Starting fs of client: ", clntAddr)
////	validatedPaths, err := preBackupChecks(backupConfig.Paths, clntAddr)
////	if err != nil {
////		log.Error("Cannot validate paths on client side")
////		return err
////	}
////	// TODO Here paths can be removed bases on excluded
////	log.Debugf("Got validated paths from client: %s starting fs...", validatedPaths)
////	err = outprotoapi.SendBackupRequest(validatedPaths, clntAddr)
////	if err != nil {
////		log.Error("Triggering fs failed!")
////		return err
////	}
////	log.Info("Backup has been triggered properly!")
////	return nil
////}
//
//// GetAllIntegratedClients simply fetching clients from clients configuration file, at least now and shows them
//func GetAllIntegratedClients() []config.Client {
//	return config.GetAllClients()
//}
