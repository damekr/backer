package transfer

import (
	"bufio"
	"encoding/gob"

	"github.com/d8x/bftp/storage"
	"github.com/damekr/backer/cmd/bacsrv/db"
	"github.com/damekr/backer/pkg/bftp"
	"github.com/sirupsen/logrus"
)

type RestoreSession struct {
	MainSession *MainSession
}

var logRestore = logrus.WithFields(logrus.Fields{"prefix": "transfer:restore"})

func CreateRestoreSession(mainSession *MainSession) *RestoreSession {
	return &RestoreSession{
		MainSession: mainSession,
	}
}

func (r *RestoreSession) HandleRestoreSession(assetID int) error {
	logRestore.Debugln("Handling incomming TGET transfer type")
	logRestore.Infoln("Client requesting asset ID: ", assetID)
	err := r.sendAssetMetadata(assetID)
	if err != nil {
		log.Errorln("ERROR in sending asset metadata, err: ", err)
		return err
	}
	modifiedAssetMetadata, err := r.receiveModifiedAssetMetadata()
	if err != nil {
		log.Errorln(err)
	}
	log.Debugln("Requesting assets according to metadata: ", modifiedAssetMetadata)

	for i := 0; i < len(modifiedAssetMetadata.FilesMetadata); i++ {

		fileT := new(bftp.FileMetadata)
		fileTEmpty := new(bftp.FileMetadata)

		// Decoding metadata of file to be transferred
		requestedAssetDecoder := gob.NewDecoder(r.MainSession.Conn)
		err := requestedAssetDecoder.Decode(&fileT)
		if err != nil {
			logRestore.Errorln("Could not decode FileMetadata struct, error: ", err)
			fileTEnc := gob.NewEncoder(r.MainSession.Conn)
			if err := fileTEnc.Encode(&fileTEmpty); err != nil {
				logRestore.Errorln("Could not encode empty FileMetadata struct, err: ", err)
				return err
			}
			return err
		}

		logRestore.Debugf("Checking if file %s exists", fileT.LocationOnServer)
		if !storage.CheckIfFileExists(fileT.LocationOnServer) {
			logRestore.Debugf("File: %s does not exist", fileT.LocationOnServer)
			fileTEncNotExist := gob.NewEncoder(r.MainSession.Conn)
			if err := fileTEncNotExist.Encode(&fileTEmpty); err != nil {
				logRestore.Errorln("Could not encode empty FileMetadata struct, err: ", err)
				return err
			}
		}

		// Sending empty as acknowledge
		err = r.MainSession.sendEmptyAckMessage()
		if err != nil {
			return err
		}
		logRestore.Debugln("Handling transfer with sending file to client, file: ", fileT.LocationOnServer)

		// Sending file
		err = r.uploadFile(fileT.LocationOnServer, fileT.FileSize)
		if err != nil {
			logRestore.Errorln("Could not send file, err: ", err.Error())
		}

		// Receiving file acknowledge
		fileSize := new(bftp.FileAcknowledge)
		fileSizeEncoder := gob.NewDecoder(r.MainSession.Conn)
		if err := fileSizeEncoder.Decode(&fileSize); err != nil {
			logRestore.Errorln("Could not decode FileAcknowledge struct, err: ", err)
			return err
		}
		logRestore.Debugln("Received file acknowledge, file size: ", fileSize.Size)
	}
	return nil
}

func (r *RestoreSession) sendAssetMetadata(assetID int) error {
	logRestore.Debugln("Sending asset metadata: ", assetID)
	connEncoder := gob.NewEncoder(r.MainSession.Conn)
	asset, err := db.DB().ReadAssetMetadata(assetID)
	if err != nil {
		log.Errorln("Cannot read asset metadata, err: ", err)
	}
	if err := connEncoder.Encode(&asset); err != nil {
		logRestore.Errorln("Could not encode asset struct, err: ", err)
		return err
	}
	return nil
}

func (r *RestoreSession) receiveModifiedAssetMetadata() (*bftp.AssetMetadata, error) {
	modifiedAssetMetadata := new(bftp.AssetMetadata)
	connDecoder := gob.NewDecoder(r.MainSession.Conn)
	if err := connDecoder.Decode(&modifiedAssetMetadata); err != nil {
		logRestore.Errorln("Could not decode modifiedAsset struct, err: ", err)
		return modifiedAssetMetadata, err
	}
	return modifiedAssetMetadata, nil
}

func (r *RestoreSession) uploadFile(localFilePath string, size int64) error {
	logRestore.Debugln("Starting sending fileReaderCloser: ", localFilePath)
	fileReaderCloser, err := r.MainSession.Storage.ReadFile(localFilePath)
	if err != nil {
		logRestore.Errorln("Cannot create localfile to write, err: ", err)
		return err
	}
	defer fileReaderCloser.Close()
	writer := bufio.NewWriter(r.MainSession.Conn)
	reader := bufio.NewReader(fileReaderCloser)
	var readFromFile int64
	var wroteToConnection int64
	buffer := make([]byte, r.MainSession.Transfer.Buffer)
	if size < int64(r.MainSession.Transfer.Buffer) {
		logRestore.Debugln("Shrinking buffer to file size: ", size)
		buffer = make([]byte, size)
	}
	for {
		read, err := reader.Read(buffer)
		if err != nil {
			logRestore.Errorln("Could not read from fileReaderCloser reader, error: ", err)
			break
		}
		readFromFile += int64(read)
		wrote, err := writer.Write(buffer[:read])
		if err != nil {
			logRestore.Errorln("Could not write to connection buffer, error: ", err)
			break
		}
		wroteToConnection += int64(wrote)
		if wroteToConnection == size {
			logRestore.Debugln("Wrote all data to connection")
			break
		}
	}
	writer.Flush()
	return nil
}
