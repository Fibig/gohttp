package httpserver

import (
	"fmt"
	"net"
)

type HttpServer struct {
	Host   string
	Port   uint16
	Routes *[]Route
}

var (
	CLRF             = "\r\n"
	HttpVersion      = "1.1"
	HttpServerName   = "FibigHttpServerYeah"
	RequestBodyLimit = 1024 * 100
)

func NewHttpServer(host string, port uint16) (*HttpServer, error) {
	return &HttpServer{Host: host, Port: port, Routes: &[]Route{}}, nil
}

func (server *HttpServer) Start() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%d", server.Host, server.Port))
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
		go handleRequest(conn, server)
	}
}
