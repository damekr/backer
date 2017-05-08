package db

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/clientsconfig"
	"github.com/damekr/backer/tiedot/db"
	"github.com/damekr/backer/tiedot/dberr"
	"strings"
)

var (
	errClientExist        = errors.New("Client exists in database")
	errClientDoesNotExist = errors.New("Client does not exist")
	errClientResource     = errors.New("Cannot add client resources")
)

func init() {
	log.Debug("Initializing db module")
}

func getClientsCollection(dbconn *db.DB) (*db.Col, error) {
	log.Debug("Opening Clients DB Doc")
	coll := dbconn.Use(clientsDocName)
	if coll == nil {
		log.Error("Something happened, clients doc has not been created")
		return coll, fmt.Errorf("Clients doc does not exist")
	}
	return coll, nil

}

// AddClient adds client into database, assumes that CID is uniq.
func AddClient(client *clientsconfig.Client) error {
	/* ClientResources Schema
	name: <name>,
	address: <ip_address>,
	cid: <client_identification number>,
	backupconfigid: <Id number>,
	platform: <linux/windows/darwin>,

	*/
	dbConnection, err := OpenDB()
	if err != nil {
		log.Error("Cannot open DB")
		return err
	}
	defer dbConnection.Close()
	clntCol, err := getClientsCollection(dbConnection)
	if err != nil {
		return err
	}
	log.Debugf("Adding client: %s into database ", client.Name)

	clnt := map[string]interface{}{
		"name":           client.Name,
		"address":        client.Address,
		"cid":            client.CID,
		"backupconfigid": client.BackupID,
		"platform":       client.Platform,
	}
	log.Debug("Adding client: %s resources", client.Name)
	// TODO: Check if client is already in the DB
	docID, err := clntCol.Insert(clnt)
	if err != nil {
		log.Error("Cannot add client resources, error: ", err.Error())
		return errClientResource
	}
	log.Debug("Client resources have been sucessfully added with doc ID: ", docID)
	return nil
}

func RemoveClient(name string) error {
	dbConnection, err := OpenDB()
	if err != nil {
		log.Error("Cannot open clients db to remove client: ", name)
		return err
	}
	defer dbConnection.Close()
	log.Debugf("Removing client: %s from clients DB", name)
	clntCol, err := getClientsCollection(dbConnection)
	if err != nil {
		return err
	}
	id, err := findClientIDByName(name, clntCol)
	if err != nil {
		return err
	}
	if id == 0 {
		return errClientDoesNotExist
	}
	log.Debug("Removing client with ID: ", id)
	if err := clntCol.Delete(id); dberr.Type(err) == dberr.ErrorNoDoc {
		log.Warning("Client with given name is not available in DB")
		return errClientDoesNotExist
	}
	return nil

}

func findClientIDByName(name string, clntCol *db.Col) (int, error) {
	// TODO: Specify if clients names must be case sensitive
	log.Debug("Looking for client doc using name: ", name)
	var clntID int
	m := make(map[string]string)
	clntCol.ForEachDoc(func(id int, docContent []byte) (willMoveOn bool) {
		if err := json.Unmarshal(docContent, &m); err != nil {
			log.Error("Could not unmarshal db doc in finding client by id: ", err.Error())
		}
		log.Debug("Processing client: ", m)
		if strings.ToLower(m["name"]) == strings.ToLower(name) {
			clntID = id
			return false
		}
		return true
	})
	return clntID, nil
}
