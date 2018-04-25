package fs

import (
	"io"

	"github.com/damekr/backer/pkg/bftp"
	"github.com/sirupsen/logrus"
)

var (
	log = logrus.WithFields(logrus.Fields{"prefix": "fs"})
)

type FileSystem interface {
	CreateFile(metadata bftp.FileMetadata) error
	WriteToFile(metadata bftp.FileMetadata) (io.WriteCloser, error)
	CreateDir(dirMetadata bftp.DirMetadata) error
	ReadFile(path string) (io.ReadCloser, error)
	ReadFileMetadata(path string) (*bftp.FileMetadata, error)
	ReadDirsMetadata(dirPaths []string) ([]*bftp.DirMetadata, error)
	CheckIfFileExists(path string) bool
}

type BackupObjects struct {
	Files []string
	Dirs  []string
}
