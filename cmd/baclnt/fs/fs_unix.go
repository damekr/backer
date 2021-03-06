// +build linux darwin

package fs

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/damekr/backer/pkg/bftp"
)

type LocalFileSystem struct {
}

func NewLocalFileSystem() LocalFileSystem {
	return LocalFileSystem{}
}

func (l LocalFileSystem) CreateFile(metadata bftp.FileMetadata) error {
	log.Debugln("Creating file: ", metadata.NameWithPath)
	file, err := os.Create(metadata.NameWithPath)
	defer file.Close()
	if err != nil {
		return err
	}
	err = file.Chmod(metadata.Mode)
	if err != nil {
		return err
	}
	// Not all filesystems supports
	// err = file.Chown(metadata.UID, metadata.GID)
	// if err != nil {
	// 	return err
	// }
	log.Debugf("File %s has been created", metadata.NameWithPath)
	return nil
}

func (l LocalFileSystem) CreateDir(metadata bftp.DirMetadata) error {
	err := l.createFileDir(metadata)
	if err != nil {
		log.Errorln("Cannot create directory, err: ", err)
		return err
	}
	return nil
}

func (l LocalFileSystem) createFileDir(metadata bftp.DirMetadata) error {
	log.Debugln("Creating directory: ", metadata.Path)
	return os.MkdirAll(metadata.Path, metadata.Mode)
}

func (l LocalFileSystem) ReadFile(filePath string) (io.ReadCloser, error) {
	log.Debugln("Reading file: ", filePath)
	file, err := os.Open(filePath)
	if err != nil {
		log.Errorln("Cannot open file for reading, err: ", err)
		return nil, err
	}
	return io.ReadCloser(file), nil
}

func (l LocalFileSystem) WriteToFile(metadata bftp.FileMetadata) (io.WriteCloser, error) {
	log.Debugln("Opening to write file: ", metadata.NameWithPath)
	file, err := os.OpenFile(metadata.NameWithPath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
		log.Errorln("Could not open file to create writer, err: ", err)
		return nil, err
	}
	return io.WriteCloser(file), nil
}

func (l LocalFileSystem) ReadFileMetadata(filePath string) (*bftp.FileMetadata, error) {
	file, err := os.Open(filePath)
	defer file.Close()
	if err != nil {
		return nil, err
	}
	fileInfo, err := file.Stat()
	if err != nil {
		return nil, err
	}
	fileMetadata := new(bftp.FileMetadata)
	fileMetadata.NameWithPath = file.Name()
	fileMetadata.FileSize = fileInfo.Size()
	fileMetadata.ModTime = fileInfo.ModTime()
	fileMetadata.Mode = fileInfo.Mode()
	return fileMetadata, nil
}

func (l LocalFileSystem) CheckIfFileExists(fullFilePath string) bool {
	if _, err := os.Stat(fullFilePath); err != nil {
		log.Errorln(err)
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func (l LocalFileSystem) ReadDirsMetadata(dirPaths []string) ([]*bftp.DirMetadata, error) {
	var dirsMetadata []*bftp.DirMetadata
	var err error
	for _, v := range dirPaths {
		dirMetadata, err := l.readDirMetadata(v)
		if err != nil {
			log.Errorln("Cannot read dir metadata, err: ", err)
		} else {
			dirsMetadata = append(dirsMetadata, dirMetadata)
		}
	}
	return dirsMetadata, err
}

func (l LocalFileSystem) readDirMetadata(path string) (*bftp.DirMetadata, error) {
	dir, err := os.Open(path)
	defer dir.Close()
	if err != nil {
		return nil, err
	}
	dirInfo, err := dir.Stat()
	if err != nil {
		return nil, err
	}
	dirMetadata := new(bftp.DirMetadata)
	dirMetadata.Path = path
	dirMetadata.Mode = dirInfo.Mode()
	dirMetadata.ModTime = dirInfo.ModTime()
	dirMetadata.BackupTime = time.Now().String()

	return dirMetadata, nil
}

func (l LocalFileSystem) ReadBackupObjectsLocations(paths []string) (BackupObjects, error) {
	log.Debug("Reading backup objects with paths:", paths)
	backupObjects := BackupObjects{}
	validatedPaths := l.validatePaths(paths)
	for i := range validatedPaths {
		err := filepath.Walk(validatedPaths[i], func(path string, info os.FileInfo, err error) error {
			if info.Mode().IsRegular() {
				log.Debugf("Adding file %s to list", path)
				backupObjects.Files = append(backupObjects.Files, path)
			} else if info.Mode()&os.ModeSymlink != 0 {
				log.Debugln("Found symlink: ", path)
			} else if info.Mode().IsDir() {
				log.Debugln("Found dir: ", path)
				backupObjects.Dirs = append(backupObjects.Dirs, path)
			} else {
				log.Debug("Found not regular file: ", path)
			}
			return nil
		})
		if err != nil {
			return backupObjects, err
		}
	}
	return backupObjects, nil
}

func (l LocalFileSystem) validatePaths(paths []string) []string {
	var validatedPaths []string
	for _, p := range paths {
		log.Debugln("Checking path: ", p)
		_, err := os.Stat(p)
		if err != nil {
			log.Warningf("AbsolutePath %s does not exist\n", p)
		} else {
			validatedPaths = append(validatedPaths, p)
		}
	}
	return validatedPaths
}

func (l LocalFileSystem) handleSymlinkFile() {

}

func (l LocalFileSystem) handleHardlinkFile() {

}

func calculateMD5Sum(fileLocation string) (string, error) {
	log.Debugf("Calculating file %s md5 checksum", fileLocation)
	file, err := os.Open(fileLocation)
	if err != nil {
		return "", err
	}
	defer file.Close()
	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	hashInBytes := hash.Sum(nil)[:16]
	returnMD5String := hex.EncodeToString(hashInBytes)
	return returnMD5String, nil
}
