package httpcompress

import (
	"compress/gzip"
	"log"
	"net/http"
	"strings"
)

type gzipResponseWriter struct {
	http.ResponseWriter
	Writer *gzip.Writer
}

func (g gzipResponseWriter) Write(b []byte) (int, error) {
	return g.Writer.Write(b)
}

type GZipNode struct {
	Node
}

func (node *GZipNode) Process(writer http.ResponseWriter, request *http.Request) {
	hit := node.GetChain().GetLocal("hit").(bool)
	if hit || !strings.Contains(request.Header.Get("Accept-Encoding"), "gzip") {
		node.Next(writer, request)
		return
	}
	log.Println("Hit Gzip")
	writer.Header().Set("Content-Encoding", "gzip")
	writer.Header().Set("Vary", "Accept-Encoding")

	gz := gzip.NewWriter(writer)
	defer gz.Close()

	gw := gzipResponseWriter{ResponseWriter: writer, Writer: gz}
	node.GetChain().PutLocal("hit", true)

	node.Next(gw, request)
}
