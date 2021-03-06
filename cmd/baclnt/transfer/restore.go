package transfer

import (
	"bufio"
	"encoding/gob"
	"io"

	"github.com/damekr/backer/cmd/baclnt/fs"
	"github.com/damekr/backer/pkg/bftp"
)

type RestoreSession struct {
	MainSession *MainSession
	FileSystem  fs.FileSystem
}

func CreateRestoreSession(mainSession *MainSession, fileSystem fs.FileSystem) *RestoreSession {
	return &RestoreSession{
		MainSession: mainSession,
		FileSystem:  fileSystem,
	}
}

func (r *RestoreSession) receiveAssetMetadata(options bftp.RestoreOptions) (*bftp.AssetMetadata, error) {
	log.Debugln("Receiving asset metadata")
	assetMetadata := new(bftp.AssetMetadata)
	connDecoder := gob.NewDecoder(r.MainSession.Conn)
	err := connDecoder.Decode(&assetMetadata)
	if err != nil {
		log.Println("Could not decode FileMetadata struct from server, err: ", err)
		return nil, err
	}
	log.Debugln("Received asset metadata: ", assetMetadata)
	// Here modifying assets according to restore options, not sure if it should be here
	modifiedAssets := r.modifyAssetMetadataWithRestoreOptions(assetMetadata, options)
	err = r.sendBackModifiedAssetMetadata(modifiedAssets)
	if err != nil {
		log.Errorln("Could not send modified asset")
	}
	return modifiedAssets, nil
}

func (r *RestoreSession) modifyAssetMetadataWithRestoreOptions(assetsMetadata *bftp.AssetMetadata, options bftp.RestoreOptions) *bftp.AssetMetadata {
	if options.WholeBackup && options.BasePath == "" {
		log.Debugln("Whole backup, the same path")
		return assetsMetadata
	}
	if !options.WholeBackup && options.BasePath == "" {
		log.Debugln("Single object, the same path")
		var newFilesMetadataList []bftp.FileMetadata
		var newDirMetadataList []bftp.DirMetadata
		for _, o := range options.ObjectsPaths {
			for _, f := range assetsMetadata.FilesMetadata {
				if f.NameWithPath == o {
					log.Debugln("Adding file: ", f)
					newFilesMetadataList = append(newFilesMetadataList, f)
				}
			}
			for _, k := range assetsMetadata.DirsMetadata {
				if k.Path == o {
					log.Debugln("Adding dir: ", k)
					newDirMetadataList = append(newDirMetadataList, k)
				}
			}
		}
		assetsMetadata.DirsMetadata = newDirMetadataList
		assetsMetadata.FilesMetadata = newFilesMetadataList
		log.Debugln("Metadata: ", assetsMetadata)
		return assetsMetadata
	}
	return nil
}

func (r *RestoreSession) sendBackModifiedAssetMetadata(metadata *bftp.AssetMetadata) error {
	log.Debugln("Sending back modified asset message")
	fileAEnc := gob.NewEncoder(r.MainSession.Conn)
	if err := fileAEnc.Encode(&metadata); err != nil {
		log.Errorln("Could not send acknowledge, err: ", err)
		return err
	}
	return nil
}

func (r *RestoreSession) createDirs(dirsMetadata []bftp.DirMetadata) error {
	log.Debugln("Creating directories")
	var err error
	for _, v := range dirsMetadata {
		err = r.FileSystem.CreateDir(v)
		if err != nil {
			log.Errorln("Could not create directory, error: ", err)
		}
	}
	return err
}

func (r *RestoreSession) RestoreFile(metadata bftp.FileMetadata) error {
	log.Printf("Downloading file: %r to: %r", metadata.LocationOnServer, metadata.NameWithPath)

	// Sending remote file path using FileMetadata struct

	err := r.sendRequestForFile(&metadata)
	if err != nil {
		log.Errorln("Could not send file metadata for restore, err: ", err.Error())
	}

	// Receiving empty as ack
	err = r.MainSession.receiveEmptyAckMessage()
	if err != nil {
		return err
	}

	if metadata.NameWithPath == "" {
		//	FilePath is empty so does not exist or an error on server side
		return bftp.FileDoesNotExist
	}
	err = r.FileSystem.CreateFile(metadata)
	if err != nil {
		log.Println("Cannot create file to write, err: ", err)
		return err
	}
	fileWritter, err := r.FileSystem.WriteToFile(metadata)
	if err != nil {
		log.Errorln("Could not create file writer, err: ", err)
	}
	defer fileWritter.Close()

	// Downloading file part
	log.Println("Server has such file starting downloading file to local path: ", metadata.NameWithPath)
	err = r.downloadFile(fileWritter, metadata)
	if err != nil {
		log.Println("Could not download file, error: ", err)
		return err
	}

	//Sending file acknowledge
	fileTransferAcknowledge, err := r.sendTransferAcknowledge(metadata.NameWithPath)
	if err != nil {
		log.Errorln("Could not send transfer acknowledge, err: ", err)
	}
	log.Println("File transfer acknowledge: ", fileTransferAcknowledge)

	return nil
}

func (r *RestoreSession) sendRequestForFile(metadata *bftp.FileMetadata) error {
	log.Println("Sending request for file located in remote path: ", metadata.NameWithPath)
	ftran := gob.NewEncoder(r.MainSession.Conn)
	err := ftran.Encode(&metadata)
	if err != nil {
		log.Errorln("Could not encode FileMetadata struct")
		return err
	}
	return nil
}

func (r *RestoreSession) receiveFileMetadata() (*bftp.FileMetadata, error) {
	log.Println("Waiting for accepting transferring file")
	fileRecTransfer := new(bftp.FileMetadata)
	frect := gob.NewDecoder(r.MainSession.Conn)
	err := frect.Decode(&fileRecTransfer)
	if err != nil {
		log.Println("Could not decode FileMetadata struct from server, err: ", err)
		return nil, err
	}
	return fileRecTransfer, nil
}

//TODO it is the same in backup also, consider put it into main session struct
func (r *RestoreSession) sendTransferAcknowledge(filePath string) (*bftp.FileAcknowledge, error) {
	fileMetadata, err := r.FileSystem.ReadFileMetadata(filePath)
	if err != nil {
		log.Errorln("Could not read file metadata, err: ", err)
	}
	fileSizeAckn := new(bftp.FileAcknowledge)
	fileSizeEncoder := gob.NewEncoder(r.MainSession.Conn)
	fileSizeAckn.Size = fileMetadata.FileSize
	log.Debugln("Filetransfer acknowledge struct: ", fileSizeAckn)
	if err := fileSizeEncoder.Encode(&fileSizeAckn); err != nil {
		log.Warningln("Could not encode FileAcknowledge struct, err: ", err)
		return nil, err
	}
	log.Println("Sent file size acknowledge")
	return fileSizeAckn, nil
}

func (r *RestoreSession) downloadFile(fileWritter io.Writer, metadata bftp.FileMetadata) error {
	log.Println("Starting downloading file, size: ", metadata.FileSize)
	reader := bufio.NewReader(r.MainSession.Conn)
	var readFromConnection int64
	var wroteToFile int64
	buffer := make([]byte, r.MainSession.Transfer.Buffer)
	if metadata.FileSize < int64(r.MainSession.Transfer.Buffer) {
		log.Println("Shrinking buffer to file size: ", metadata.FileSize)
		buffer = make([]byte, metadata.FileSize)
	}
	for {
		read, err := reader.Read(buffer)
		if err != nil {
			log.Println("Could not read from connection reader, error: ", err)
			break
		}
		readFromConnection += int64(read)
		wrote, err := fileWritter.Write(buffer[:read])
		if err != nil {
			log.Println("Could not write to file writter, error: ", err)
			break
		}
		wroteToFile += int64(wrote)
		if wroteToFile == metadata.FileSize {
			log.Println("Wrote all data to file")
			break
		}

	}
	//fileWritter.Flush()
	return nil
}
