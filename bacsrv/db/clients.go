// +build ignore

package db

import (
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/boltdb/bolt"
	"github.com/damekr/backer/bacsrv/config"
)

var (
	errClientExist        = errors.New("Client exists in database")
	errClientDoesNotExist = errors.New("Client does not exist")
	errClientResource     = errors.New("Cannot add client resources")
)

// AddClient adds client into database, assumes that CID is uniq.
func AddClient(clnt *config.Client) error {
	/* ClientResources Schema
	name: <name>,
	address: <ip_address>,
	cid: <client_identification number>,
	backupconfigid: <Id number>,
	platform: <linux/windows/darwin>,

	*/
	log.Debugf("Adding client %s into db", clnt.Name)
	db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(clientsBucket))
		if buf, err := json.Marshal(clnt); err != nil {
			return err
		} else if err := b.Put([]byte(clnt.CID), buf); err != nil {
			return err
		}
		return nil
	})

	log.Debugf("Client %s has been successfully added into DB", clnt.Name)
	return nil
}

func GetClient(cid string) (*config.Client, error) {
	client := &config.Client{}
	db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(clientsBucket))
		v := b.Get([]byte(cid))
		err := json.Unmarshal(v, client)
		if err != nil {
			log.Error(err)
			return err
		}
		fmt.Printf("The answer is: %#v\n", client)
		return nil
	})

	return client, nil
}

func ShowAllClients() error {
	log.Debug("Showing all clients")
	db.View(func(tx *bolt.Tx) error {
		// Assume bucket exists and has keys
		b := tx.Bucket([]byte(clientsBucket))
		b.ForEach(func(k, v []byte) error {
			fmt.Printf("key=%s, value=%s\n", k, v)
			return nil
		})
		return nil
	})
	return nil
}

func RemoveClient(name string) error {
	return nil

}
