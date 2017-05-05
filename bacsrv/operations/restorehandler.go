package operations

import (
	"net"
)

type Restore struct {
	Saveset     string `json:"saveset"`
	SavesetSize int64  `json:"savesetSize"`
}

type RestoreTriggerMessage struct {
	ClientName    string `json:"clientName"`
	RestoreConfig Restore
}

func (r *Restore) DoRestore(conn *net.Conn) error {
	return nil
}
