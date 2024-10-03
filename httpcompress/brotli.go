package httpcompress

import (
	"github.com/andybalholm/brotli"
	"log"
	"net/http"
	"strings"
)

type brotliResponseWriter struct {
	http.ResponseWriter
	Writer *brotli.Writer
}

func (b brotliResponseWriter) Write(p []byte) (int, error) {
	return b.Writer.Write(p)
}

type BrotliNode struct {
	Node
}

func (node *BrotliNode) Process(writer http.ResponseWriter, request *http.Request) {
	hit := node.GetChain().GetLocal("hit").(bool)
	if hit || !strings.Contains(request.Header.Get("Accept-Encoding"), "br") {
		node.Next(writer, request)
		return
	}
	log.Println("Hit Brotli")
	writer.Header().Set("Content-Encoding", "br")
	writer.Header().Set("Vary", "Accept-Encoding")

	bw := brotli.NewWriter(writer)
	defer bw.Close()

	bwr := brotliResponseWriter{ResponseWriter: writer, Writer: bw}
	node.GetChain().PutLocal("hit", true)

	node.Next(bwr, request)
}
