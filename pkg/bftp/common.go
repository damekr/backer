package bftp

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
	"os"
	"time"
)

// TODO Exclude whole bftp protocol to pkg directory
const (
	AuthenticationPassed = "passed"
	AuthenticationFailed = "failed"
)

var (
	PORT         = "8090"
	SERVER       = "0.0.0.0"
	PROTOVERSION = "0.1"
	PASSWORD     = "john"
	KEY          = "the-key-has-to-be-32-bytes-long!"
	BUFFERSIZE   = 1024

	// Protocol Operations
	TGET = "GET"
	TPUT = "PUT"
	//
)

// Global Errors
var (
	AuthenticationFailedError         = errors.New("BFTP: Authentication failed")
	FileDoesNotExist                  = errors.New("BFTP: File does not exist")
	ServerDoesNotSupportSuchOperation = errors.New("BFTP: Server does not support such operation")
	ProtocolVersionMismatch           = errors.New("BFTP: Protocol version mismatch")
)

type ConnParameters struct {
	Server string
	Port   string
}

type Negotiate struct {
	ClientName   string
	ProtoVersion string
}

type Authenticate struct {
	CiperText []byte
}

type Transfer struct {
	TransferType  string
	AssetID       int
	ObjectsNumber int
	Buffer        int
}

type AssetMetadata struct {
	ClientName    string `json:"clientName"`
	ID            int    `json:"ID"`
	BucketPath    string `json:"bucketLocation"`
	SavesetPath   string `json:"savesetLocation"`
	FilesMetadata []FileMetadata
	DirsMetadata  []DirMetadata
}

type RestoreOptions struct {
	WholeBackup  bool     `json:"wholeBackup"`
	ObjectsPaths []string `json:"restoreObjectsPaths"`
	BasePath     string   `json:"basePath"`
}

type FileMetadata struct {
	Name             string      `json:"fileName"`
	FullPath         string      `json:"fullPathOfFile"`
	FileSize         int64       `json:"fileSize"`
	ModTime          time.Time   `json:"modTime"`
	UID              int         `json:"userID"`
	GID              int         `json:"groupID"`
	Mode             os.FileMode `json:"fileMode"`
	DirMode          os.FileMode `json:"dirMode"`
	Checksum         string      `json:"fileMD5Checksum"`
	LocationOnServer string      `json:"locationOnServer"`
	BackupTime       string      `json:"backupTime"`
}

type DirMetadata struct {
	Path       string      `json:"dirPath"`
	ModTime    time.Time   `json:"modTime"`
	UID        int         `json:"userID"`
	GID        int         `json:"groupID"`
	Mode       os.FileMode `json:"dirMode"`
	BackupTime string      `json:"backupTime"`
}

type FileTransferInfo struct {
}

type FileAcknowledge struct {
	Size int64
}

type EmtpyAck struct {
}

type BFTPAcknowledgeMessage struct {
	Status bool
}

// NewConnParameters creates struct with parameters used for connection
func NewConnParameters() *ConnParameters {
	connParameters := ConnParameters{}
	connParameters.Server = SERVER
	connParameters.Port = PORT
	return &connParameters
}

func Encrypt(password []byte, key []byte) ([]byte, error) {

	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	return gcm.Seal(nonce, nonce, password, nil), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
	c, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}
