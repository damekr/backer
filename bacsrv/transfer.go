package main

import (
	"fmt"
	"net"
	"strings"
	"strconv"
)
const BUFFERSIZE = 1024

func InitTransferServer(){
	listener, err := net.Listen("tcp", "localhost:27001")
	if err != nil {
		fmt.Println("An error" + err.Error())
	}

	bufferFileName := make([]byte, 64)
	bufferFileSize := make([]byte, 10)
	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("An error: "+err.Error())
		}else{
			fmt.Println("New connection estabilished")
		}
		sizeBytes, err := connection.Read(bufferFileSize)
		if err != nil {
			fmt.Println("An error during read: "+err.Error())
		}
		fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
		fmt.Printf("Received: %d\n", fileSize)

		nameBytes, err := connection.Read(bufferFileName)
		if err != nil {
			fmt.Println("An error during read: "+err.Error())
		}
		fileName := strings.Trim(string(bufferFileName), ":")
		fmt.Printf("Received: %s\n", fileName)

		// file, err := os.Create("tmp.tar")
		// if err != nil{
		// 	fmt.Println("An error during openning file")
		// }
		// defer file.Close()
		// n, err := io.Copy(connection, file)
		// if err == io.EOF {
		// 	fmt.Println("An error during copy to file. "+ err.Error())
		// }
		fmt.Println(sizeBytes, "bytes received")
		fmt.Println(nameBytes, "bytes received")
		connection.Close()
	}
}


// func Listener(){
// 	connection, err := net.Dial("tcp", "localhost:27001")
// 	if err != nil {
// 		fmt.Println("Error listetning: ", err)
// 		os.Exit(1)
// 	}
// 	defer connection.Close()
// 	fmt.Println("Server started! Waiting for connections...")
		// fmt.Println("Client connected")
        // bufferFileName := make([]byte, 64)
		// bufferFileSize := make([]byte, 10)
	
		// connection.Read(bufferFileSize)
		// fileSize, _ := strconv.ParseInt(strings.Trim(string(bufferFileSize), ":"), 10, 64)
	
		// connection.Read(bufferFileName)
		// fileName := strings.Trim(string(bufferFileName), ":")
		// fmt.Println("File name: ", fileName)
		// if fileName == "" {
		// 	fileName = "tmp.tar"
		// }

		// newFile, err := os.Create(fileName)
	
		// if err != nil {
		// 	panic(err)
		// }
		// defer newFile.Close()
		// var receivedBytes int64
	
		// for {
		// 	if (fileSize - receivedBytes) < BUFFERSIZE {
		// 		io.CopyN(newFile, connection, (fileSize - receivedBytes))
		// 		connection.Read(make([]byte, (receivedBytes+BUFFERSIZE)-fileSize))
		// 		break
		// 	}
		// io.CopyN(newFile, connection, BUFFERSIZE)
		// receivedBytes += BUFFERSIZE
		// }
		// fmt.Println("Received file completely!")


