package dataproto

import (
	"encoding/gob"
	log "github.com/Sirupsen/logrus"
	"io"
	"net"
	"os"
)

func init() {
	log.Debug("Initializing transfer protocol")
}

// Transfer is a header of transfer connection, always is send at the beginning of data transfer
type Transfer struct {
	TType string
	From  string
}

type FileTransferInfo struct {
	Name     string
	Location string
	Size     int64
	UID      int
	GID      int
	Mode     os.FileMode
	Checksum string
}

// SendDataTypeHeader encoding data to be send over socket as first chunk of data
func SendDataTypeHeader(transfer *Transfer, conn net.Conn) error {
	log.Debug("Marshaling data for header transfer")
	enc := gob.NewEncoder(conn)
	err := enc.Encode(transfer)
	if err != nil {
		log.Errorf("Cannot encode transfer header with type: %s, and from: %s", transfer.TType, transfer.From)
		return err
	}
	return nil
}

// UnmarshalTransferHeader gets buffered data and decode them
func UnmarshalTransferHeader(conn net.Conn) (*Transfer, error) {
	log.Debug("Unmarshaling data from transfered header")
	var transfer Transfer
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&transfer)
	if err != nil {
		log.Error("Cannot decode transfer header data")
		return &transfer, err
	}
	return &transfer, nil
}

// SendFileInfoHeader sends each time before file trsansfer information about file being transfered
func SendFileInfoHeader(fileInfo *FileTransferInfo, conn net.Conn) error {
	log.Debugf("Sending file header:  %#v", fileInfo)
	enc := gob.NewEncoder(conn)
	err := enc.Encode(fileInfo)
	if err != nil {
		log.Error("Could not encode file info header")
		return err
	}
	return nil
}

// UnmarshalFileInfoHeader getting information from connection about file being transfered
func UnmarshalFileInfoHeader(conn net.Conn) (*FileTransferInfo, error) {
	log.Debug("Unmarshaling file info header")
	var fileInfo FileTransferInfo
	dec := gob.NewDecoder(conn)
	err := dec.Decode(&fileInfo)
	switch {
	case err == io.EOF:
		return &fileInfo, err

	case err != nil:
		log.Error("Cannot decode file info header, error: ", err.Error())
		return &fileInfo, err
	}
	s, err := conn.Write([]byte("\n"))
	if err != nil {
		log.Error("ERORR: ", err.Error())
	}
	log.Debug("Responded: ", s)
	return &fileInfo, nil
}
