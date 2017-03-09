package archiver

import (
	log "github.com/Sirupsen/logrus"
	"os"
	"path/filepath"
	"syscall"
)

var TempDir string

type FileInfo struct {
	Path   string
	Size   int64
	Exists bool
}

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
