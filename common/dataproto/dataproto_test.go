package dataproto


import (
	"net"
	"fmt"
	"testing"
	"io/ioutil"
	"crypto/rand"
	"os"
)

const PORT = ":12000"

type tempFile struct {
	Name string
	Size int
}

func newTempFile(name string, size int) tempFile{
	return tempFile{
		Name: name,
		Size: size,
	}
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// Note that err == nil only if we read len(b) bytes.
	if err != nil {
		return nil, err
	}

	return b, nil
}

func (t *tempFile) create() (string, error) {
	tmpFile, err := ioutil.TempFile("", t.Name)
	if err != nil{
		fmt.Errorf("Cannot create temp file")
		return "", err
	}
	data, err := generateRandomBytes(t.Size)
	if err != nil {
		fmt.Errorf("Could not generate random data")
		return "", err
	}
	s, err := tmpFile.Write(data)
	if err != nil {
		fmt.Errorf("Error while writing content to file")
		return "", err
	}
	fmt.Printf("Written %s bytes to file", s)
	return tmpFile.Name(), nil
}

func (t *tempFile) cleanup(fullpath string) error {
	fmt.Println("Cleaning up temp file")
	err := os.Remove(fullpath)
	return err
}

func startServer() net.Listener{
	ln, err := net.Listen("tcp", PORT)
	if err != nil {
		fmt.Print("Cannot listen, error: ", err.Error())
	}
	return ln
}

func sendTransferHeader(transferType string){
	conn, err := net.Dial("tcp", PORT)
	if err != nil {
		fmt.Print("Error while connecting")
	}
	defer conn.Close()
	transfer := New("localhost", conn)
	transfer.SendTypeHeader(transferType)
}


func TestTransfer_SendTypeHeader(t *testing.T) {
	ln := startServer()
	transferType := "backup"
	go func() {
		sendTransferHeader(transferType)
	}()
	defer ln.Close()
	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Errorf("Error while listening")
		}
		defer conn.Close()

		fmt.Println("Handling connection")
		transfer2 := New("client", conn)
		transferType, err := transfer2.ReceiveTypeHeader()
		if transferType != transferType{
			t.Fail()
		}
		if err != nil {
			t.Fatal(err)
		}

		return
	}
}


//func sendFile() {
//	conn, err := net.Dial("tcp", PORT)
//	if err != nil {
//		fmt.Print("Error while connecting")
//	}
//	defer conn.Close()
//	transfer := New("client", conn)
//	transfer.SendFile()
//}

func TestTransfer_SendFile(t *testing.T) {
	fileSize := 1024000 << 8 // 250MB

	tmpFile := tempFile{
		Name: "dummyFile",
		Size: fileSize,
	}
	fileName, err := tmpFile.create()
	if err != nil {
		t.Fatal("Could not create temp file")
	}
	t.Log("Created temp file: ", fileName)
	t.Log("Cleaning up")
	err = tmpFile.cleanup(fileName)
	if err != nil{
		t.Fatal(err)
	}
}