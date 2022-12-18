package server

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type gzipWriter struct {
	http.ResponseWriter           // original
	Writer              io.Writer // gzip wrapper writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	// w.Writer будет отвечать за gzip-сжатие, поэтому пишем в него
	return w.Writer.Write(b)
}

func gzipSend(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
			next.ServeHTTP(writer, request)
			return
		}

		gz, err := gzip.NewWriterLevel(writer, gzip.BestSpeed)
		if err != nil {
			_, _ = writer.Write([]byte(err.Error()))
			return
		}
		defer gz.Close()

		writer.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(gzipWriter{ResponseWriter: writer, Writer: gz}, request)
	})
}

func gzipReceive(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		if request.Header.Get(`Content-Encoding`) == `gzip` {
			gz, err := gzip.NewReader(request.Body)
			if err != nil {
				http.Error(writer, err.Error(), http.StatusInternalServerError)
			}
			request.Body = gz
		}

		next.ServeHTTP(writer, request)
	})
}
