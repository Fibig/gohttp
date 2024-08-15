package httpserver

import (
	"fmt"
	"net/http"
)

type HTTPResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func (r HTTPResponse) String() string {
	response := fmt.Sprintf("HTTP/%s %d %s\r\n", HttpVersion, r.StatusCode, http.StatusText(r.StatusCode))
	for k, v := range r.Headers {
		response += fmt.Sprintf("%s:%s\r\n", k, v)
	}
	response += "\r\n" + string(r.Body)
	return response
}
