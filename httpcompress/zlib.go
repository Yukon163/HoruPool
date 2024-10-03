package httpcompress

import (
	"compress/zlib"
	"log"
	"net/http"
	"strings"
)

type zlibResponseWriter struct {
	http.ResponseWriter
	Writer *zlib.Writer
}

func (z zlibResponseWriter) Write(b []byte) (int, error) {
	return z.Writer.Write(b)
}

type ZlibNode struct {
	Node
}

func (node *ZlibNode) Process(writer http.ResponseWriter, request *http.Request) {
	hit := node.GetChain().GetLocal("hit").(bool)
	if hit || !strings.Contains(request.Header.Get("Accept-Encoding"), "deflate") {
		node.Next(writer, request)
		return
	}
	log.Println("Hit Zlib")
	writer.Header().Set("Content-Encoding", "deflate")
	writer.Header().Set("Vary", "Accept-Encoding")

	zw := zlib.NewWriter(writer)
	defer zw.Close()

	zwr := zlibResponseWriter{ResponseWriter: writer, Writer: zw}
	node.GetChain().PutLocal("hit", true)

	node.Next(zwr, request)
}
