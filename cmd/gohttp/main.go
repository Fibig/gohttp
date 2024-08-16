package main

import (
	"fmt"
	"net/http"

	"example.com/Fibig/gohttp/pkg/httpserver"
)

func main() {
	server, err := httpserver.NewHttpServer("localhost", 9000)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	server.Use("aaa", func(hr httpserver.HttpRequest) []byte {
		body := []byte("aaa")
		response, _ := httpserver.GetResponse(http.StatusOK, make(map[string]string), &body, httpserver.GzipCompress)
		return response
	})

	server.Static("/public", "public")

	server.Start()
}
