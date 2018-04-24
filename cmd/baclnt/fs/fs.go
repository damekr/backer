package fs

import (
	"io"

	"github.com/damekr/backer/pkg/bftp"
)

type FileSystem interface {
	CreateFile(metadata bftp.FileMetadata) error
	CreateDir(dirMetadata bftp.DirMetadata) error
	ReadFile(path string) (io.ReadCloser, error)
	WriteFile(metadata bftp.FileMetadata) (io.WriteCloser, error)
	ReadFileMetadata(path string) (*bftp.FileMetadata, error)
	ReadDirMetadata(path string) (*bftp.DirMetadata, error)
	CheckIfFileExists(path string) bool
}

type BackupObjects struct {
	Files []string
	Dirs  []string
}
