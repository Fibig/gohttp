package httpserver

import (
	"maps"
	"strconv"
	"time"

	"github.com/gabriel-vasile/mimetype"
)

const (
	HeaderContentLength   = "Content-Length"
	HeaderContentType     = "Content-Type"
	HeaderServer          = "Server"
	HeaderDate            = "Date"
	HeaderContentEncoding = "Content-Encoding"
)

func getResponseHeaders(contentType, contentLength string, additionalHeaders map[string]string) string {
	// set must need headers
	headers := map[string]string{
		HeaderServer:        HttpServerName,
		HeaderContentType:   contentType,
		HeaderContentLength: contentLength,
		HeaderDate:          time.Now().UTC().String(),
	}

	// add to/overwrite headers with additional headers
	maps.Copy(headers, additionalHeaders)

	// transform all headers to string
	headersString := ""
	for k, v := range headers {
		headersString += k + ":" + v + CLRF
	}

	return headersString
}

func getContentType(data []byte) string {
	return mimetype.Detect(data).String()
}

func getContentLength(data []byte) string {
	return strconv.Itoa(len(data))
}
