// +build linux darwin

package fs

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"os"
	"path"
	"path/filepath"

	"github.com/damekr/backer/pkg/bftp"
	log "github.com/sirupsen/logrus"
)

const (
	defaultDirPerm = os.FileMode(0744)
)

type LocalFileSystem struct {
}

func NewLocalFileSystem() LocalFileSystem {
	return LocalFileSystem{}
}

func (l LocalFileSystem) CreateFile(metadata bftp.FileMetadata) error {
	log.Debugln("Creating file: ", metadata)
	err := l.createFileDir(metadata)
	if err != nil {
		log.Errorln("Cannot rebuild file dir, err: ", err)
		return err
	}
	file, err := os.Create(path.Join(metadata.FullPath, metadata.Name))
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
	return nil
}

func (l LocalFileSystem) createFileDir(metadata bftp.FileMetadata) error {
	log.Debugln("Rebuilding needed directory for file full path: ", metadata.FullPath)
	return os.MkdirAll(metadata.FullPath, metadata.DirMode)
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

func (l LocalFileSystem) WriteFile(metadata bftp.FileMetadata) (io.WriteCloser, error) {
	log.Debugln("Writing file: ", metadata)
	file, err := os.OpenFile(path.Join(metadata.FullPath, metadata.Name), os.O_APPEND|os.O_WRONLY, os.ModeAppend)
	if err != nil {
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
	fileMetadata.Name = file.Name()
	fileMetadata.FileSize = fileInfo.Size()
	fileMetadata.ModTime = fileInfo.ModTime()
	fileMetadata.Mode = fileInfo.Mode()
	fileMetadata.FullPath = path.Dir(filePath)
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

func (l LocalFileSystem) ExpandDirsForFiles(paths []string) ([]string, error) {
	log.Debug("Checking absolute paths for: ", paths)
	validatedPaths := l.validatePaths(paths)
	var filesToBeBackedUp []string
	for i := range validatedPaths {
		err := filepath.Walk(validatedPaths[i], func(path string, info os.FileInfo, err error) error {
			if info.Mode().IsRegular() {
				log.Debugf("Adding file %s to list", path)
				filesToBeBackedUp = append(filesToBeBackedUp, path)
			} else if info.Mode()&os.ModeSymlink != 0 {
				log.Debugln("Found symlink: ", path)
			} else if info.Mode().IsDir() {
				log.Debugln("Found dir: ", path)
			} else {
				log.Debug("Found not regular file: ", path)
			}
			return nil
		})
		if err != nil {
			return filesToBeBackedUp, err
		}
	}
	return filesToBeBackedUp, nil
}

func (m LocalFileSystem) validatePaths(paths []string) []string {
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
