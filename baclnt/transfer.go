package main

import (
	"fmt"
	"net"
	"strconv"
	"os"
	"io"
	)

const BUFFERSIZE = 1024

type Connection struct {
    Port    int8
    Host    string
    Timeout int
    BUFFERSIZE int
}

func NewConnection(Port int8, Host string, Timeout int, BufferSize int) *Connection{
    return &Connection{
        Port: Port,
        Host: Host,
        Timeout: Timeout,
        BUFFERSIZE: BufferSize,
    }
}

func InitConnection() {
	connection, err := net.Dial("tcp", "localhost:27001")
	if err != nil {
		fmt.Println("An error during connection: "+ err.Error())
	}
	fileName, fileSize := readFileMetadata()
	outSize, err := connection.Write([]byte(fileSize))
	if err != nil {
		fmt.Println("An error occured: "+err.Error())
	}
	fmt.Println(outSize, "bytes sent Name") 
	outName, err := connection.Write([]byte(fileName))
	if err != nil {
		fmt.Println("An error occured: "+err.Error())
	}
	fmt.Println(outName, "bytes sent size")
	defer connection.Close()
    // Start sending file
	sendBuffer := make([]byte, BUFFERSIZE)
	file := readFile()
	defer file.Close()
	for {
		_, err := file.Read(sendBuffer)
		if err == io.EOF {
			break
		}
		connection.Write(sendBuffer)
	}
	fmt.Println("File has been sent, closing connection!")
	return 
}

func readFile() *os.File{
	file, err := os.Open("/home/damian/tmp.tar")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	return file
}

func readFileMetadata() (string, string){
	file, err := os.Open("/home/damian/tmp.tar")
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fileInfo, err := file.Stat()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
	fileName := fillString(fileInfo.Name(), 64)
	fmt.Printf("File Size: %s\nFile Name: %s\n", fileSize, fileName)
	return fileName, fileSize
}

// func sendFileToClient(connection net.Conn, f string) {
// 	fmt.Println("A client has connected!")
// 	defer connection.Close()
// 	file, err := os.Open(f)
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fileInfo, err := file.Stat()
// 	if err != nil {
// 		fmt.Println(err)
// 		return
// 	}
// 	fileSize := fillString(strconv.FormatInt(fileInfo.Size(), 10), 10)
// 	fileName := fillString(fileInfo.Name(), 64)
// 	fmt.Printf("File Size: %s, FileName: %s ", fileSize, fileName)
// 	fmt.Println("Sending filename and filesize!")
// 	connection.Write([]byte(fileSize))
// 	connection.Write([]byte(fileName))
// 	fmt.Println("Sending filename and filesize!")
// 	sendBuffer := make([]byte, BUFFERSIZE)
// 	fmt.Println("Start sending file!")
// 	for {
// 		_, err = file.Read(sendBuffer)
// 		if err == io.EOF {
// 			break
// 		}
// 		connection.Write(sendBuffer)
// 	}
// 	fmt.Println("File has been sent, closing connection!")
// 	return
// }

func fillString(retunString string, toLength int) string {
	for {
		lengtString := len(retunString)
		if lengtString < toLength {
			retunString = retunString + ":"
			continue
		}
		break
	}
	return retunString
}