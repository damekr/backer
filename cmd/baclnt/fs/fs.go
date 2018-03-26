package fs

import (
	"io"

	"github.com/damekr/backer/pkg/bftp"
)

type FileSystem interface {
	CreateFile(metadata bftp.FileMetadata) error
	ReadFile(path string) (io.ReadCloser, error)
	WriteFile(metadata bftp.FileMetadata) (io.WriteCloser, error)
	ReadFileMetadata(path string) (*bftp.FileMetadata, error)
	CheckIfFileExists(path string) bool
}
