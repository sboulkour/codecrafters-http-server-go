package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
)

const (
	CRLF = "\r\n"
)

func main() {

	directoryPtr := flag.String("directory", "/tmp", "Directory to serve")
	flag.Parse()
	fmt.Println("Working directory: ", *directoryPtr)

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		fmt.Println("Failed to bind to port 4221")
		os.Exit(1)
	}

	for {
		connection, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		go handleConnection(connection, directoryPtr)
	}

}

func handleConnection(connection net.Conn, directoryPtr *string) {
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
	} else {
		switch requestUri[1] {
		case "echo":
			responseBody := requestUri[2]
			responseHeaders := fmt.Sprintf("HTTP/1.1 200 OK"+CRLF+"Content-Type: text/plain"+CRLF+"Content-Length: %d"+CRLF+CRLF, len(responseBody))
			connection.Write([]byte(responseHeaders + responseBody))
			connection.Close()
		case "user-agent":
			for index, line := range request {
				if strings.HasPrefix(request[index], "User-Agent: ") {
					responseBody := strings.TrimPrefix(line, "User-Agent: ")
					responseHeaders := fmt.Sprintf("HTTP/1.1 200 OK"+CRLF+"Content-Type: text/plain"+CRLF+"Content-Length: %d"+CRLF+CRLF, len(responseBody))
					connection.Write([]byte(responseHeaders + responseBody))
					connection.Close()
				}
			}
		case "files":
			if len(requestUri[2]) > 2 {
				filePath := filepath.Join(*directoryPtr, requestUri[2])

				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					connection.Write([]byte("HTTP/1.1 404 Not Found" + CRLF + CRLF))
					fmt.Println("Error reading file: ", err.Error())
					connection.Close()
					return
				}
				responseBody := string(fileContent)
				responseHeaders := fmt.Sprintf("HTTP/1.1 200 OK"+CRLF+"Content-Type: application/octet-stream"+CRLF+"Content-Length: %d"+CRLF+CRLF, len(responseBody))
				connection.Write([]byte(responseHeaders + responseBody))
				connection.Close()
			} else {
				connection.Write([]byte("HTTP/1.1 404 Not Found" + CRLF + CRLF))
				connection.Close()
			}
		default:
			connection.Write([]byte("HTTP/1.1 404 Not Found" + CRLF + CRLF))
			connection.Close()
		}

	}
}
