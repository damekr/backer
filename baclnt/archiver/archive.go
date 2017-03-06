package archiver

import (
	"archive/tar"
	log "github.com/Sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
	"strconv"

	"time"
)

var archiveName string

func generateArchiveName() {
	hostname, err := os.Hostname()
	if err != nil {
		log.Errorf("Cannot get hostname to set archive name")
	}
	now := time.Now()
	nanos := now.UnixNano()
	shorted := strconv.FormatInt(nanos, 10)
	archiveName = hostname + shorted + ".tar"
}

type Archive struct {
	Name  string
	Paths []string
	Size  int64
}

func NewArchive(Paths []string) *Archive {
	generateArchiveName()
	return &Archive{
		Paths: Paths,
		Size:  0,
	}
}

func (a *Archive) validatePaths() {
	validatedPaths := []string{}
	for _, path := range a.Paths {
		log.Printf("Checking path %s", path)
		_, err := os.Stat(path)
		if err != nil {
			log.Printf("Path %s does not exist\n", path)
		} else {
			validatedPaths = append(validatedPaths, path)
		}
	}
	a.Paths = validatedPaths
}

func addFile(tw *tar.Writer, location string) error {
	file, err := os.Open(location)
	if err != nil {
		return err
	}
	defer file.Close()
	if stat, err := file.Stat(); err == nil {
		header := new(tar.Header)
		header.Name = location
		header.Size = stat.Size()
		header.Mode = int64(stat.Mode())
		header.ModTime = stat.ModTime()
		if err := tw.WriteHeader(header); err != nil {
			return err
		}
		if _, err := io.Copy(tw, file); err != nil {
			return err
		}
	}
	return nil
}

// getFilesAbsPaths checks if given file exists and sets abs path to be included in tar archive
func getFilesAbsPaths(paths []string) []string {
	fileList := []string{}
	for i := range paths {
		err := filepath.Walk(paths[i], func(path string, f os.FileInfo, err error) error {
			if f.Mode().IsRegular() {
				fileList = append(fileList, path)
			}
			return nil
		})
		if err != nil {
			log.Fatalln(err)
		}
	}
	for _, file := range fileList {
		log.Println("File: ", file)
	}
	return fileList
}

// MakeArchive is able to create an archive from given paths
func (a *Archive) MakeArchive() string {
	a.validatePaths()
	log.Printf("Making tar package from %s \n", a.Paths)
	absFilePaths := getFilesAbsPaths(a.Paths)
	tarAbsolutePath := TempDir + "/" + archiveName
	log.Debugf("Absoulte path of package %s", tarAbsolutePath)
	tarfile, err := os.Create(tarAbsolutePath)
	if err != nil {
		log.Panicf("Location %s does not exists", tarAbsolutePath)
	}
	defer tarfile.Close()
	tw := tar.NewWriter(tarfile)
	defer tw.Close()
	for i := range absFilePaths {
		if err := addFile(tw, absFilePaths[i]); err != nil {
			log.Fatalln(err)
		}
	}
	fi, err := tarfile.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	a.Size = fi.Size()
	log.Printf("File size: %d Bytes", fi.Size())
	return tarfile.Name()
}
