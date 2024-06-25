package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

const (
	CRLF = "\r\n"
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

	requestBuffer := make([]byte, 1024)
	n, err := connection.Read(requestBuffer)
	if err != nil {
		fmt.Println("Error reading request: ", err.Error())
		os.Exit(1)
	}

	request := strings.Split(string(requestBuffer[:n]), CRLF)
	requestStatusLine := strings.Split(request[0], " ")
	fmt.Println(requestStatusLine)

	requestUri := strings.Split(requestStatusLine[1], "/")

	if requestStatusLine[1] == "/" {
		connection.Write([]byte("HTTP/1.1 200 OK" + CRLF + CRLF))
	} else if requestUri[1] == "echo" {
		responseBody := requestUri[2]
		responseHeaders := fmt.Sprintf("HTTP/1.1 200 OK"+CRLF+"Content-Type: text/plain"+CRLF+"Content-Length: %d"+CRLF+CRLF, len(responseBody))
		connection.Write([]byte(responseHeaders + responseBody))
	} else if requestUri[1] == "user-agent" {
		for index, line := range request {
			if strings.HasPrefix(request[index], "User-Agent: ") {
				responseBody := strings.TrimPrefix(line, "User-Agent: ")
				responseHeaders := fmt.Sprintf("HTTP/1.1 200 OK"+CRLF+"Content-Type: text/plain"+CRLF+"Content-Length: %d"+CRLF+CRLF, len(responseBody))
				connection.Write([]byte(responseHeaders + responseBody))
			}
		}
	} else {
		connection.Write([]byte("HTTP/1.1 404 Not Found" + CRLF + CRLF))
	}
}
