package restapi

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/job"
	"github.com/damekr/backer/bacsrv/status"
	"github.com/gorilla/mux"
)

// ContentType specified type of data sending from application
const ContentType = "text/json"

// Bacsrv represents information about server
type Bacsrv struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	GoVersion string `json:"goversion"`
}

// StartServerRestAPI starts http server on given port
func StartServerRestAPI(srvConfig *config.ServerConfig) {
	router := NewRouter()
	log.Debug("Starting Server RESTAPi on port: ", srvConfig.RestAPIPort)
	log.Fatal(http.ListenAndServe(":"+srvConfig.RestAPIPort, router))
}

// Index shows basic information about server
func Index(w http.ResponseWriter, r *http.Request) {
	bacsrv := Bacsrv{
		Name:      "Backer Server",
		Version:   "0.1",
		GoVersion: "1.7.5",
	}
	log.Printf("BACSRV %#v", bacsrv)
	w.Header().Set("Content-Type", ContentType)
	if err := json.NewEncoder(w).Encode(bacsrv); err != nil {
		log.Error("Cannot encode server metadata into json")
	}

}

// StatusIndex responds current server status
func StatusIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ContentType)
	serverStatus := status.GetSeverStatus()
	if err := json.NewEncoder(w).Encode(serverStatus); err != nil {
		log.Error("Cannot encode memory information")
	}

}

// ShowClients responds all information about added clients
func ShowClients(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ContentType)
	if err := json.NewEncoder(w).Encode(job.GetAllIntegratedClients()); err != nil {
		log.Error("Cannot encode clients structs")
	}
}

// HelloMessage determines hostname about client with given name or ip
type HelloMessage struct {
	Hostname string `json:"hostname"`
}

// ShowClientStatus simply connects over grpc to client with specified name and read hostname
func ShowClientStatus(w http.ResponseWriter, r *http.Request) {
	// // TODO It must be somehow specified what will be used, could be a name and during the process read ip address and then send request
	vars := mux.Vars(r)
	clientName := vars["clientName"]
	log.Debug("Received arguments: ", clientName)
	clientHostname, err := job.SendHelloMessageToClient(clientName)
	if err != nil {
		errorMessage := &HelloMessage{
			Hostname: "Error, Cannot find given client",
		}

		if err := json.NewEncoder(w).Encode(errorMessage); err != nil {
			log.Error("Cannot encode memory information")
		}
		return
	}
	message := &HelloMessage{
		Hostname: clientHostname,
	}
	if err := json.NewEncoder(w).Encode(message); err != nil {
		log.Error("Cannot encode memory information")
	}

}

type IntegrationMessage struct {
	ClientName string `json:"clientName"`
	Status     string `json:"integrationStatus"`
	BackupID   string `json:"backupIdentyfication"`
}

func IntegrateNewClient(w http.ResponseWriter, r *http.Request) {
	var integrateMessage IntegrationMessage
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Errorf("Cannot read Integration Body")
	}
	if err := r.Body.Close(); err != nil {
		log.Errorf("Unexpected end of body in integration request")
	}
	if err := json.Unmarshal(body, &integrateMessage); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Errorf("Cannot decode json, check header")
		}
	}
	log.Printf("Received integration message: %#v", integrateMessage)
	fmt.Fprint(w, "OK")

}

// TODO Exclude body handler to on generic method, to prevent repeating the same stuff

/*
TriggerClientBackkup json message
{
"paths": ["sad", "sd"],
"excludedPaths": ["asd", "sd"],
"retentionTime": "365"
}
*/

func TriggerClientBackup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientName := vars["clientName"]
	log.Debug("Received arguments: ", clientName)
	var backupConfigMassage config.BackupConfig
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Error("Cannot read body of Trigger message backup")
	}
	if err := r.Body.Close(); err != nil {
		log.Errorf("Unexpected end of body in backup trigger request")
	}
	if err := json.Unmarshal(body, &backupConfigMassage); err != nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(422) // unprocessable entity
		if err := json.NewEncoder(w).Encode(err); err != nil {
			log.Errorf("Cannot decode json, check header")
		}
	}
	clientBackupMessage := &config.BackupTriggerMessage{
		ClientName:   clientName,
		BackupConfig: backupConfigMassage,
	}
	log.Printf("Received backup message: %#v", backupConfigMassage)
	log.Printf("Full backup message %#v", clientBackupMessage)
	// manager.SendBackupTriggerMessage(clientBackupMessage)
	fmt.Fprint(w, "OK")
}
