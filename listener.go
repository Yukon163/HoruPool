package main

import (
	"Horu/httpcompress"
	"Horu/internal"
	"fmt"
	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"runtime"
)

func convertMaps(maps map[internal.SrcUrl]internal.DstUrl) map[internal.DstUrl]*httputil.ReverseProxy {
	proxyMap := make(map[internal.DstUrl]*httputil.ReverseProxy)
	for _, dst := range maps {
		parse, _ := url.Parse(string(dst))
		reverseProxy := httputil.NewSingleHostReverseProxy(parse)
		runtime.KeepAlive(reverseProxy)
		proxyMap[dst] = reverseProxy
	}
	runtime.KeepAlive(proxyMap)
	return proxyMap
}

func proxyPort(entry internal.PortEntry) {
	runtime.LockOSThread()

	maps := entry.Maps
	proxyMaps := convertMaps(maps)
	mux := http.NewServeMux()
	addr := fmt.Sprintf(":%d", entry.Port)
	http3Server := http3.Server{
		Addr:    addr,
		Handler: mux,
		QUICConfig: &quic.Config{
			Allow0RTT:       false,
			EnableDatagrams: true,
		},
	}
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("get", r.Host, r.ProtoMajor)

		srcUrl := internal.SrcUrl(r.Host)
		dstUrl, exists := maps[srcUrl]
		if !exists {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		r.Host = string(dstUrl)

		chain := new(httpcompress.Chain)

		brotliNode := new(httpcompress.BrotliNode)
		zlibNode := new(httpcompress.ZlibNode)
		gzipNode := new(httpcompress.GZipNode)
		brotliNode.BindChain(chain)
		zlibNode.BindChain(chain)
		gzipNode.BindChain(chain)
		chain.Init(brotliNode, zlibNode, gzipNode, httpcompress.FinalNode(func(w http.ResponseWriter, r *http.Request) {
			reverseServer := chain.GetLocal("reverseServer").(*httputil.ReverseProxy)
			runtime.Gosched()
			reverseServer.ServeHTTP(w, r)
		}))
		chain.PutLocal("hit", false)
		chain.PutLocal("reverseServer", proxyMaps[dstUrl])
		brotliNode.Process(w, r)
		runtime.GC()
	}))

	err := http3Server.ListenAndServeTLS(config.CertFile, config.KeyFile)

	log.Fatalf("proxy port %d listen http3 fail: %v", entry.Port, err)
}
