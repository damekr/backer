// +build linux darwin

package fs

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

type FileInfo struct {
	Path   string
	Size   int64
	Exists bool
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

type FS struct {
	AbsoluteFilesPaths []string
}

func (f *FS) checkPathForRegularFile(path string, file os.FileInfo, err error) error {
	if file.Mode().IsRegular() {
		log.Debugf("Adding file %s to list", path)
		f.AbsoluteFilesPaths = append(f.AbsoluteFilesPaths, path)
	} else {
		log.Debug("Found not regular file: ", path)
	}
	return nil
}

func (f *FS) GetAbsolutePaths(paths []string) []string {
	log.Debug("Checking absolutive paths for: ", paths)
	validatedPaths := f.ValidatePaths(paths)
	for i := range validatedPaths {
		err := filepath.Walk(validatedPaths[i], f.checkPathForRegularFile)
		if err != nil {
			log.Error(err)
		}
	}
	return f.AbsoluteFilesPaths
}

func (f *FS) ValidatePaths(paths []string) []string {
	var validatedPaths []string
	for _, vPath := range paths {
		log.Printf("Checking vPath %s", vPath)
		_, err := os.Stat(vPath)
		if err != nil {
			log.Printf("Path %s does not exist\n", vPath)
		} else {
			validatedPaths = append(validatedPaths, vPath)
		}
	}
	return validatedPaths
}

//func GetFilesInformations(paths []string) []FileInfo {
//	log.Debug("Getting files informations")
//	filesInfo := []FileInfo{}
//	absFilePaths := GetAbsolutePaths(paths)
//	for _, f := range absFilePaths {
//		fileInfo := new(FileInfo)
//		info, err := os.Stat(f)
//		if err != nil {
//			log.Error("Cannot open file: %s", f)
//		}
//		fileInfo.Exists = true
//		log.Debug("Adding path ", f)
//		fileInfo.Path = f
//		log.Debug("File size: ", info.Size())
//		fileInfo.Size = info.Size()
//		filesInfo = append(filesInfo, *fileInfo)
//	}
//	log.Debug("Size of list, ", len(filesInfo))
//	return filesInfo
//}

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

func ReadFileHeader(fileLocation string) (*FileTransferInfo, error) {
	var fileHeader FileTransferInfo
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

type FileSystem struct {
	Path string
}

func NewFS(path string) *FileSystem {
	return &FileSystem{
		Path: path,
	}
}

func (f *FileSystem) CreateFile(name string) (*os.File, error) {
	file, err := os.Create(path.Join(f.Path, name))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func (f *FileSystem) OpenFile(name string) (*os.File, error) {
	file, err := os.Open(path.Join(f.Path, name))
	if err != nil {
		return nil, err
	}
	return file, nil
}

func CheckIfFileExists(fullFilePath string) bool {
	if _, err := os.Stat(fullFilePath); err != nil {
		fmt.Print(err)
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func GetFileSize(fullPath string) int64 {
	file, err := os.Open(fullPath)
	defer file.Close()
	fstat, err := file.Stat()
	if err != nil {
		log.Println("Cannot do stat on file, returning 0")
		return 0
	}
	return fstat.Size()
}
