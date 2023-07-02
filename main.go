package main

import (
	"net"
	"strings"
)

func main() {
	startServer()
}

func startServer() {

	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(connection net.Conn) {
	for {

		data := make([]byte, 1024)

		_, err := connection.Read(data)
		if err != nil {
			if err.Error() == "EOF" {
				return
			}
			panic(err)
		}

		request := parseHTTPRequest(string(data))

		response := HTTPResponse{
			Version: request.Version,
			Status:  "200 OK",
			Headers: request.Headers,
			Data:    request.Path,
		}

		response.send(connection)
		err = connection.Close()
		if err != nil {
			panic(err)
		}

		return
	}
}

type HTTPRequest struct {
	Method  string
	Path    string
	Version string
	Headers map[string]string
	Data    string
}

func parseHTTPRequest(data string) HTTPRequest {

	// Loop through lines
	lines := strings.Split(data, "\r\n")

	var httpRequest HTTPRequest
	httpRequest.Headers = map[string]string{}

	// Parse first line
	firstLine := strings.Split(lines[0], " ")
	httpRequest.Method = firstLine[0]
	httpRequest.Path = firstLine[1]
	httpRequest.Version = firstLine[2]

	// Skip first line
	lines = lines[1:]

	isData := false
	for _, line := range lines {

		// If line is empty, we are at the data
		if line == "" {
			isData = true
			continue
		}

		if isData {
			httpRequest.Data += line
			continue
		}

		// Parse header
		split := strings.Split(line, ": ")
		httpRequest.Headers[split[0]] = split[1]
	}

	return httpRequest
}

type HTTPResponse struct {
	Version string
	Status  string
	Headers map[string]string
	Data    string
}

func (r *HTTPResponse) String() string {

	response := r.Version + " " + r.Status + "\r\n"

	for key, value := range r.Headers {
		response += key + ": " + value + "\r\n"
	}

	response += "\r\n" + r.Data

	return response
}

func (r *HTTPResponse) send(conn net.Conn) {
	_, err := conn.Write([]byte(r.String()))
	if err != nil {
		panic(err)
	}
}
