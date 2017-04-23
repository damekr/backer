// +build linux darwin

package backup

import (
	"crypto/md5"
	"encoding/hex"
	log "github.com/Sirupsen/logrus"
	"github.com/damekr/backer/common/dataproto"
	"io"
	"os"
	"path"
	"path/filepath"
)

var AbsoluteFilesPaths []string

type FileInfo struct {
	Path   string
	Size   int64
	Exists bool
}

func checkPathForRegularFile(path string, f os.FileInfo, err error) error {
	if f.Mode().IsRegular() {
		log.Debugf("Adding file %s to list", path)
		AbsoluteFilesPaths = append(AbsoluteFilesPaths, path)
	} else {
		log.Debug("Found not regular file: ", path)
	}
	return nil
}

// TODO: GENERAL: Proper handling files privilages.

// GetAbsolutePaths makes actually two things resolve files in given paths and checks if exist.
// TODO Refactor me :)
func GetAbsolutePaths(paths []string) []string {
	log.Debug("Checking absolutive paths for: ", paths)
	validatedPaths := ValidatePaths(paths)
	for i := range validatedPaths {
		err := filepath.Walk(validatedPaths[i], checkPathForRegularFile)
		if err != nil {
			log.Error(err)
		}
	}
	return AbsoluteFilesPaths
}

func ValidatePaths(paths []string) []string {
	validatedPaths := []string{}
	for _, path := range paths {
		log.Printf("Checking path %s", path)
		_, err := os.Stat(path)
		if err != nil {
			log.Printf("Path %s does not exist\n", path)
		} else {
			validatedPaths = append(validatedPaths, path)
		}
	}
	return validatedPaths
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

// func CreateTempDir(location string) {
// 	log.Debugf("Creating temporary catalouge to store temp data in: %s", location)
// 	err := os.MkdirAll(location, 0700)
// 	if err != nil {
// 		log.Errorf("Cannot create temporary catalogue for storing data, exiting...")
// 		os.Exit(5)
// 	}
// 	TempDir = location
// }

// func GetTempAvailableSpace() int64 {
// 	log.Debug("Checking if temporary directory has enough space for restore")
// 	fs := syscall.Statfs_t{}
// 	err := syscall.Statfs(TempDir, &fs)
// 	if err != nil {
// 		log.Error("Cannot check file system capacity")
// 	}
// 	free := int64(fs.Bfree) * int64(fs.Bsize)
// 	log.Debug("Available space in temporary directory: ", free)
// 	return free
// }

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

func ReadFileHeader(fileLocation string) (*dataproto.FileTransferInfo, error) {
	var fileHeader dataproto.FileTransferInfo
	info, err := os.Stat(fileLocation)
	if err != nil {
		log.Errorf("File %s does not exist", fileLocation)
		return &fileHeader, err
	}
	checksum, err := calculateMD5Sum(fileLocation)
	if err != nil {
		log.Error("Was not able to calculate checksum, setting 0")
		checksum = "0"
	}
	fileHeader.Checksum = checksum
	fileHeader.Location = fileLocation
	fileHeader.Mode = info.Mode()
	fileHeader.Size = info.Size()
	fileHeader.Name = path.Base(fileLocation)
	return &fileHeader, nil
}
