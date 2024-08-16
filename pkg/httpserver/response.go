package httpserver

import (
	"fmt"
	"maps"
	"net/http"
	"strconv"
)

type HttpResponse struct {
	StatusCode int
	Headers    map[string]string
	Body       []byte
}

func (r HttpResponse) String() string {
	response := fmt.Sprintf("HTTP/%s %d %s\r\n", HttpVersion, r.StatusCode, http.StatusText(r.StatusCode))
	for k, v := range r.Headers {
		response += fmt.Sprintf("%s:%s\r\n", k, v)
	}

	response += "\r\n" + string(r.Body)
	return response
}

func GetResponse(status int, additionalHeaders map[string]string, responseBody *[]byte, compressionMethod func(*[]byte) (*[]byte, map[string]string, error)) ([]byte, error) {
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
