package httpserver

import (
	"bytes"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
)

type HttpRequest struct {
	HttpVersion string
	Path        string
	Method      string
	Headers     map[string]string
	BodyRaw     []byte
}

func NewHttpRequest(data []byte) (*HttpRequest, error) {
	// dataSplit seperates start line and header from the body
	dataSplit := bytes.SplitN(data, []byte(CLRF+CLRF), 2)

	// dataAboveBody helps to seperate the start line and the headers
	dataAboveBody := bytes.Split(dataSplit[0], []byte(CLRF))

	// raw data
	startLineRaw := dataAboveBody[0]
	headersRaw := dataAboveBody[1:]
	bodyRaw := dataSplit[1]

	// parse startLine data
	startLine := strings.Split(string(startLineRaw), " ")
	httpVersion, path, method := startLine[0], startLine[1], startLine[2]

	// parse headerRaw into Header Map
	headers := make(map[string]string)
	for _, v := range headersRaw {
		header := strings.Split(strings.ReplaceAll(string(v), " ", ""), ":")
		headers[header[0]] = header[1]
	}

	// check if content length header is set
	if headers[HeaderContentLength] == "" {
		return &HttpRequest{
			HttpVersion: httpVersion,
			Path:        path,
			Method:      method,
			Headers:     headers,
			BodyRaw:     bodyRaw,
		}, nil
	}

	// check for content length
	contentLength, err := strconv.Atoi(headers[HeaderContentLength])
	if err != nil {
		return nil, err
	}

	if len(startLineRaw)+len(headersRaw)+len([]byte(CLRF+CLRF))+contentLength > RequestBodyLimit {
		return nil, fmt.Errorf("payload too large")
	}

	return &HttpRequest{
		HttpVersion: httpVersion,
		Path:        path,
		Method:      method,
		Headers:     headers,
		BodyRaw:     bodyRaw,
	}, nil
}

func handleRequest(conn net.Conn, server *HttpServer) {
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
	response := []byte{}

	// Check if route exists
	isRouteExisting := false
	for _, route := range *server.Routes {
		fmt.Printf("%s ? %s\n", route.Path, request.Path[1:])
		if strings.HasPrefix(request.Path[1:], route.Path) {
			fmt.Println("Route match: " + route.Path)
			response = route.Func(*request)
			isRouteExisting = true
			break
		}
	}

	// 404 Not Found
	if !isRouteExisting {
		notFoundBody := []byte("404 NOT FOUND")
		response, err = GetResponse(http.StatusNotFound, map[string]string{}, &notFoundBody, GzipCompress)
	}

	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	conn.Write(response)
}
