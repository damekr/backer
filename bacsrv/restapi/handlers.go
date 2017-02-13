package restapi

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/config"
	"github.com/damekr/backer/bacsrv/manager"
	"github.com/damekr/backer/bacsrv/status"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

const ContentType = "text/json"

type Bacsrv struct {
	Name      string `json:"name"`
	Version   string `json:"version"`
	GoVersion string `json:"goversion"`
}

func init() {
	log.SetFormatter(&log.TextFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func StartServerRestApi(config *config.ServerConfig) {
	router := NewRouter()
	log.Debug("Starting Server RESTAPi on port: ", config.MgmtPort)
	log.Fatal(http.ListenAndServe(":"+config.MgmtPort, router))
}

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

func StatusIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ContentType)
	serverStatus := status.GetSeverStatus()
	if err := json.NewEncoder(w).Encode(serverStatus); err != nil {
		log.Error("Cannot encode memory information")
	}

}

func ShowClients(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Todo show:")
}

type HelloMessage struct {
	Hostname string `json:"hostname"`
}

func ShowClientStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientName := vars["clientName"]
	log.Debug("Received arguments: ", clientName)
	clientHostname, err := manager.HelloMessageManager(clientName)
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
	//TODO Change this endpoint
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
