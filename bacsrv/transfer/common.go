package transfer

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"errors"
	"io"
)

var (
	PORT         = 8090
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
	AuthenticationFailed              = errors.New("BFTP: Authentication failed")
	FileDoesNotExist                  = errors.New("BFTP: File does not exist")
	ServerDoesNotSupportSuchOperation = errors.New("BFTP: Server does not support such operation")
	ProtocolVersionMismatch           = errors.New("BFTP: Protocol version mismatch")
)

//

type ConnParameters struct {
	Server string
	Port   int
}

type Negotiate struct {
	ProtoVersion string
}

type Authenticate struct {
	CiperText []byte
}

type Transfer struct {
	TransferType string
	Buffer       int
}

type FileTransfer struct {
	//TODO Name is not used anywhere yet and is not send
	Name     string
	FullPath string
	FileSize int64
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
