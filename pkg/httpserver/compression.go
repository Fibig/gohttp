package httpserver

import (
	"bytes"
	"compress/gzip"
)

func GzipCompress(data *[]byte) (*[]byte, map[string]string, error) {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)

	_, err := w.Write(*data)
	if err != nil {
		return nil, map[string]string{}, err
	}
	err = w.Close()
	if err != nil {
		return nil, map[string]string{}, err
	}

	compressedData := b.Bytes()

	return &compressedData, map[string]string{HeaderContentEncoding: "gzip"}, nil
}

func NoCompression(data *[]byte) (*[]byte, map[string]string, error) {
	return data, map[string]string{}, nil
}
