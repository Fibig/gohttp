package httpserver

import (
	"fmt"
	"maps"
	"net"
	"net/http"
	"strconv"
)

type HttpServer struct {
	host string
	port uint16
}

var (
	CLRF             = "\r\n"
	HttpVersion      = "1.1"
	HttpServerName   = "FibigHttpServerYeah"
	RequestBodyLimit = 1024 * 100
)

func NewHttpServer(host string, port uint16) (*HttpServer, error) {
	return &HttpServer{host: host, port: port}, nil
}

func (server *HttpServer) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.host, server.port))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	defer listener.Close()

	fmt.Println("Listening on:", listener.Addr().String())

	for {
		// Accept incoming connections
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Handle client connection in a goroutine
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	defer conn.Close()

	// Create a buffer to read data into
	buffer := make([]byte, RequestBodyLimit)

	// read into buffer
	bufferEnd, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Parse Request and create instance of HttpRequest Type
	request, err := NewHttpRequest(buffer[:bufferEnd])
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Routing
	responseBody := make([]byte, 0)
	if request.Path == "/" {
		responseBody = []byte("<html><h1>This is root</h1></html>")
	} else {
		responseBody = []byte("<html><h1>This is everything else</h1></html>")
	}

	response, err := getResponse(http.StatusOK, map[string]string{}, &responseBody, GzipCompress)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn.Write(response)
}

func getResponse(status int, additionalHeaders map[string]string, responseBody *[]byte, compressionMethod func(*[]byte) (*[]byte, map[string]string, error)) ([]byte, error) {
	// apply compression
	responseBodyCompressed, compressionHeader, err := compressionMethod(responseBody)
	if err != nil {
		return nil, err
	}

	// overwrite compression headers with additional headers so that additional headers have more priority
	maps.Copy(compressionHeader, additionalHeaders)

	// build together the http message response and compress with the given compressionMethod
	response := []byte(getResponseStartLine(status) + getResponseHeaders(getContentType(*responseBody), getContentLength(*responseBodyCompressed), compressionHeader) + CLRF)
	response = append(response, *responseBodyCompressed...)

	return response, nil
}

func getResponseStartLine(status int) string {
	return "HTTP/" + HttpVersion + " " + strconv.Itoa(status) + " " + http.StatusText(status) + CLRF
}
