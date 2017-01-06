package restapi

import (
	"encoding/json"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/backer/bacsrv/config"
	"github.com/backer/bacsrv/status"
	"github.com/gorilla/mux"
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

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/", Index)
	router.HandleFunc("/status", StatusIndex)
	router.HandleFunc("/clients", ShowClients)
	log.Debug("Starting Server RESTAPi on port: ", config.MgmtPort)
	log.Fatal(http.ListenAndServe(":"+config.MgmtPort, router))
}

func Index(w http.ResponseWriter, r *http.Request) {
	bacsrv := Bacsrv{
		Name:      "Backer Server",
		Version:   "0.1",
		GoVersion: "1.6.2",
	}
	w.Header().Set("Content-Type", ContentType)
	if err := json.NewEncoder(w).Encode(bacsrv); err != nil {
		log.Error("Cannot encode server metadata into json")
	}

}

func StatusIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", ContentType)
	status := status.GetSeverStatus()
	if err := json.NewEncoder(w).Encode(status); err != nil {
		log.Error("Cannot encode memory information")
	}

	//fmt.Fprintln(w, "Todo Index!")
}

func ShowClients(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Todo show:")
}
