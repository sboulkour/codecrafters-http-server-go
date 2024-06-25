package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	connection, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	request := make([]byte, 1024)
	connection.Read(request)
	if strings.HasPrefix(string(request), "GET / HTTP/1.1") {
		connection.Write([]byte("HTTP/1.1 200 OK\r\n\r\n"))
	} else {
		connection.Write([]byte("HTTP/1.1 404 Not Found\r\n\r\n"))
	}
}
