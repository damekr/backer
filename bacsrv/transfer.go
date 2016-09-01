package main

import (
	"fmt"
	"net"
	"strings"
	"strconv"
	"io"
	"os"
)

// BUFFERSIZE determines how big is piece of data that will be send in one frame
const BUFFERSIZE = 1024

func InitTransferServer(){
	listener, err := net.Listen("tcp", "localhost:27001")
	if err != nil {
		fmt.Println("An error" + err.Error())
	}
	
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("An error: "+err.Error())
		}else{
			fmt.Println("New connection estabilished")
			fileSize := GetFileSize(connection)
			fileName := GetFileName(connection)
			// Part of receiving file
			ReceiveFile(fileSize, fileName, connection)

		}

		
	}
}

// ReceiveFile is able to read data from buffer and save them in created file.
// It also checks if retrived file is equeal to sent earlier in first chunks of data.
func ReceiveFile(fileSize int64, fileName string, connection net.Conn){
	newFile, err := os.Create(fileName)
	if err != nil {
		panic(err.Error())
	}
	defer newFile.Close()
	var receivedBytes int64
	for {
		if (fileSize - receivedBytes) < BUFFERSIZE {
			io.CopyN(newFile, connection, (fileSize - receivedBytes))
			connection.Read(make([]byte, (receivedBytes +BUFFERSIZE) - fileSize))
			break
		}
		io.CopyN(newFile, connection, BUFFERSIZE)
		receivedBytes += BUFFERSIZE
		}
	fmt.Println("Received file completely!")
}

// GetFileSize gets file size from given buffer --> remember to send and receive data in proper order 
func GetFileSize(connection net.Conn) int64{
	bufferFileSize := make([]byte, 10)
	_, err := connection.Read(bufferFileSize)
	if err != nil {
		fmt.Println("An error during read: "+err.Error())
	}
	fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	fmt.Printf("Received: %d\n", fileSize)
	return fileSize
}

// GetFileName returns filename from given connection --> remember to send data and read
func GetFileName(connection net.Conn) string{
	bufferFileName := make([]byte, 64)
	_, err := connection.Read(bufferFileName)
	if err != nil {
		fmt.Println("An error during read: "+err.Error())
	}
	fileName := strings.Trim(string(bufferFileName), ":")
	fmt.Printf("Received: %s\n", fileName)
	return fileName

}

