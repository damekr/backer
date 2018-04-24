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
	modifiedAssets := r.modifyAssetMetadataWithRestoreOptions(assetMetadata, options)
	err = r.sendBackModifiedAssetMetadata(modifiedAssets)
	if err != nil {
		log.Errorln("Could not send modified asset")
	}
	return assetMetadata, nil
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
				if f.FullPath == o {
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

func (r *RestoreSession) GetFile(fileRemotePath, fileLocalPath string) error {
	log.Printf("Downloading file: %r to: %r", fileRemotePath, fileLocalPath)

	// Sending remote file path using FileMetadata struct
	fileMetadata := new(bftp.FileMetadata)
	fileMetadata.FullPath = fileRemotePath
	err := r.sendRequestForFile(fileMetadata)
	if err != nil {
		log.Errorln("Could not send file metadata for restore, err: ", err.Error())
	}

	// Receiving file size with the same struct as above
	receivedFileMetadata, err := r.receiveFileMetadata()
	if err != nil {
		log.Errorln("Could not receive file metadata, err: ", err.Error())
	}

	if receivedFileMetadata.FullPath == "" {
		//	FilePath is empty so does not exist or an error on server side
		return bftp.FileDoesNotExist
	}
	log.Println("Creating local file: ", fileLocalPath)
	err = r.FileSystem.CreateFile(*receivedFileMetadata)
	if err != nil {
		log.Println("Cannot create localfile to write, err: ", err)
		return err
	}
	// TODO, This part might not work
	fileWritter, err := r.FileSystem.WriteFile(*receivedFileMetadata)
	if err != nil {
		log.Errorln("Could not create file writter, err: ", err)
	}
	defer fileWritter.Close()

	// Downloading file part
	log.Println("Server has such file starting downloading file to local path: ", fileLocalPath)
	err = r.downloadFile(fileWritter, receivedFileMetadata)
	if err != nil {
		log.Println("Could not download file, error: ", err)
		return err
	}

	//Sending file acknowledge
	fileTransferAcknowledge, err := r.sendTransferAcknowledge(fileLocalPath)
	if err != nil {
		log.Errorln("Could not send transfer acknowledge, err: ", err)
	}
	log.Println("File transfer acknowledge: ", fileTransferAcknowledge)

	return nil
}

func (r *RestoreSession) sendRequestForFile(metadata *bftp.FileMetadata) error {
	//TODO When backup metadata will be created, then send here also backup id
	log.Println("Sending request for file located in remote path: ", metadata.FullPath)
	ftran := gob.NewEncoder(r.MainSession.Conn)
	err := ftran.Encode(&metadata)
	if err != nil {
		log.Errorln("Could not encode FileMetadata struct")
		return err
	}
	return nil
}

func (r *RestoreSession) receiveFileMetadata() (*bftp.FileMetadata, error) {
	log.Println("Waiting for accepting transfering file")
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

func (r *RestoreSession) downloadFile(fileWritter io.Writer, metadata *bftp.FileMetadata) error {
	log.Println("Starting downloading file, size: ", metadata.FileSize)
	reader := bufio.NewReader(r.MainSession.Conn)
	var readFromConnection int64
	var wroteToFile int64
	buffer := make([]byte, r.MainSession.Transfer.Buffer)
	if metadata.FileSize < int64(r.MainSession.Transfer.Buffer) {
		log.Println("Shrinking buffer to filesize: ", metadata.FileSize)
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
	// fileWritter.Flush()
	return nil
}
