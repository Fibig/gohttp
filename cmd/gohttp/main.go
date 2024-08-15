package main

import (
	"fmt"

	"example.com/Fibig/gohttp/pkg/httpserver"
)

func main() {
	server, err := httpserver.NewHttpServer("localhost", 9000)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	server.Start()
}
