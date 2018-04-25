package restore

import (
	"context"

	"github.com/damekr/backer/api/protoclnt"
	"github.com/damekr/backer/cmd/bacsrv/db"
	"github.com/damekr/backer/cmd/bacsrv/network"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/sirupsen/logrus"
)

var log = logrus.WithFields(logrus.Fields{"prefix": "task:restore"})

type Restore struct {
	ClientIP       string `json:"clientIP"`
	AssetID        int    `json:"assetID"`
	RestoreOptions bftp.RestoreOptions
	AssetsMetadata bftp.AssetMetadata
	Progress       int    `json:"-"`
	Status         bool   `json:"status"`
	BucketLocation string `json:"bucketLocation"`
}

func Create(clientIP string, assetID int, options bftp.RestoreOptions) *Restore {
	return &Restore{
		ClientIP:       clientIP,
		AssetID:        assetID,
		RestoreOptions: options,
	}
}

func (r *Restore) Run() {
	log.Println("Running restore of client: ", r.ClientIP)
	conn, err := network.EstablishGRPCConnection(r.ClientIP)
	if err != nil {
		log.Errorf("Cannot connect to address %s", r.ClientIP)
		return
	}
	// TODO Consider close grpc connection before restore gets done
	defer conn.Close()
	c := protoclnt.NewBaclntClient(conn)

	// Sending to client restore request with options
	response, err := c.Restore(context.Background(),
		&protoclnt.RestoreRequest{
			Ip:             r.ClientIP,
			WholeBackup:    r.RestoreOptions.WholeBackup,
			RestoreObjects: r.RestoreOptions.ObjectsPaths,
			BasePath:       r.RestoreOptions.BasePath,
			AssetID:        int32(r.AssetID),
		})

	if err != nil {
		log.Warningf("Could not get response from restore request, err: ", err)
		r.Status = false
		return
	}
	r.Status = true
	log.Println("Response: ", response)
}

func (r *Restore) Stop() {
	log.Println("Stopping")
}

// Setup configures restore job, should be splited into different kind of setups(singleDir, wholeBackup etc.).
func (r *Restore) Setup() error {
	backupMetadata, err := db.DB().ReadAssetMetadata(r.AssetID)
	if err != nil {
		return err
	}
	r.AssetsMetadata = *backupMetadata

	log.Debugln("Files to be restored metadata: ", r.AssetsMetadata.FilesMetadata)
	log.Debugln("Dirs to be restored metadata: ", r.AssetsMetadata.DirsMetadata)

	return nil
}
