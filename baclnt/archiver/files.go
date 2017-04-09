// +build linux darwin

package archiver

import (
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/common/dataproto"
	"os"
	"path"
	"path/filepath"
	"syscall"
)

var TempDir string

type FileInfo struct {
	Path   string
	Size   int64
	Exists bool
}

// GetAbsolutePaths makes actually two things resolve files in given paths and checks if exist.
// TODO Refactor me :)
// FIXME: If an element is empty fails
func GetAbsolutePaths(paths []string) []string {
	log.Debug("Checking absolutive paths for: ", paths)
	fileList := []string{}
	for i := range paths {
		err := filepath.Walk(paths[i], func(path string, f os.FileInfo, err error) error {
			if f.Mode().IsRegular() {
				log.Debugf("Adding file %s to list", path)
				fileList = append(fileList, path)
			} else {
				log.Debug("Found not regular file: ", path)
			}
			return nil
		})
		if err != nil {
			log.Error(err)
		}
	}
	return fileList
}

func GetFilesInformations(paths []string) []FileInfo {
	log.Debug("Getting files informations")
	filesInfo := []FileInfo{}
	absFilePaths := GetAbsolutePaths(paths)
	for _, f := range absFilePaths {
		fileInfo := new(FileInfo)
		info, err := os.Stat(f)
		if err != nil {
			log.Error("Cannot open file: %s", f)
		}
		fileInfo.Exists = true
		log.Debug("Adding path ", f)
		fileInfo.Path = f
		log.Debug("File size: ", info.Size())
		fileInfo.Size = info.Size()
		filesInfo = append(filesInfo, *fileInfo)
	}
	log.Debug("Size of list, ", len(filesInfo))
	return filesInfo
}

func CreateTempDir(location string) {
	log.Debugf("Creating temporary catalouge to store temp data in: %s", location)
	err := os.MkdirAll(location, 0700)
	if err != nil {
		log.Errorf("Cannot create temporary catalogue for storing data, exiting...")
		os.Exit(5)
	}
	TempDir = location
}

func GetTempAvailableSpace() int64 {
	log.Debug("Checking if temporary directory has enough space for restore")
	fs := syscall.Statfs_t{}
	err := syscall.Statfs(TempDir, &fs)
	if err != nil {
		log.Error("Cannot check file system capacity")
	}
	free := int64(fs.Bfree) * int64(fs.Bsize)
	log.Debug("Available space in temporary directory: ", free)
	return free
}

func ReadFileHeader(fileLocation string) (*dataproto.FileTransferInfo, error) {
	var fileHeader dataproto.FileTransferInfo
	info, err := os.Stat(fileLocation)
	if err != nil {
		log.Errorf("File %s does not exist", fileLocation)
		return &fileHeader, err
	}
	fileHeader.Location = fileLocation
	fileHeader.Mode = info.Mode()
	fileHeader.Size = info.Size()
	fileHeader.Name = path.Base(fileLocation)

	return &fileHeader, nil
}
