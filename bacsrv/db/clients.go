package db

import (
	"errors"
	"fmt"
	"github.com/HouzuoGuo/tiedot/db"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/bacsrv/clientsconfig"
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
		log.Error("Something happend, clients doc has not been created")
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
	docID, err := clntCol.Insert(clnt)
	if err != nil {
		log.Error("Cannot add client resources, error: ", err.Error())
		return errClientResource
	}
	log.Debug("Client resources have been sucessfully added with doc ID: ", docID)
	return nil
}

// func RemoveClient(name string) error {
// 	dbConnection, err := GetClientDBConnection()
// 	if err != nil {
// 		log.Error("Cannot open clients db to remove client: ", name)
// 		return err
// 	}
// 	log.Debugf("Removing client: %s from clients DB", name)

// }

// func findClientDocByName(name string, dbConn *db.DB) (*db.Col, error) {
// 	log.Debug("Looking for client doc using name: ", name)
// 	for
// }
