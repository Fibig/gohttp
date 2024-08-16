package httpserver

import (
	"log"
	"net/http"
	"os"
	"strings"
)

type Route struct {
	Path              string
	Func              RouteFunction
	StaticContentPath string
}

type RouteFunction func(HttpRequest) []byte

func (server *HttpServer) Use(path string, routeFunction RouteFunction) {
	server.addRouteToRouter(&Route{Path: path, Func: routeFunction})
}

func (server *HttpServer) Static(path, staticContentPath string) {
	server.addRouteToRouter(&Route{Path: path, Func: newStaticContentRouteFunction(path, staticContentPath), StaticContentPath: staticContentPath})
}

func (server *HttpServer) addRouteToRouter(route *Route) {
	// remove beginning "/"
	if strings.Index(route.Path, "/") == 0 {
		route.Path = route.Path[1:]
	}

	// check if route exists
	for _, v := range *server.Routes {
		if v.Path == route.Path {
			log.Panic("route already exists")
		}
	}

	*server.Routes = append(*server.Routes, *route)
}

func newStaticContentRouteFunction(path, staticContentPath string) RouteFunction {
	return func(hr HttpRequest) []byte {
		// remove beginning "/"
		if strings.Index(path, "/") == 0 {
			path = path[1:]
		}

		// remove path from request
		requestedFilePath := strings.TrimPrefix(hr.Path, "/"+path)

		// add file path to staticCointentPath with beginning "/"
		fileLocation := "./" + staticContentPath + requestedFilePath

		// read file
		fileData, err := os.ReadFile(fileLocation)
		if err != nil {
			notFoundBody := []byte("404 NOT FOUND")
			notFoundResponse, _ := GetResponse(http.StatusNotFound, map[string]string{}, &notFoundBody, GzipCompress)
			return notFoundResponse
		}

		response, _ := GetResponse(http.StatusOK, make(map[string]string), &fileData, GzipCompress)
		return response
	}
}

func (route *Route) isRouteStatic() bool {
	return route.StaticContentPath != ""
}
