package main

import (
	"AkuaProxy/confbox"
	"AkuaProxy/data"
	"AkuaProxy/httpcompress"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/quic-go/quic-go/http3"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

var config data.Config

func init() {
	confPath := flag.String("c", "/simon/conf/proxy.yaml", "config file path(yaml).")
	demo := flag.Bool("demo", false, "Give a demo config file.")
	flag.Parse()

	if *demo {
		if err := confbox.Save("./demo.yaml", data.Config{
			PointLists: []data.PointList{
				{
					Port: 8080,
					Points: map[data.SrcUrl]data.DstUrl{
						"https://source1.com:443": "https://destination1.com",
						"https://source2.com":     "https://destination2.com",
					},
				},
				{
					Port: 80,
					Points: map[data.SrcUrl]data.DstUrl{
						"https://source3.com": "https://destination3.com",
						"https://source4.com": "https://destination4.com",
					},
				},
			},
			KeyFile:  "path/to/your/server.key",
			CertFile: "path/to/your/server.crt",
		}); err != nil {
			log.Fatalln("Internal error ", err)
			return
		}
		os.Exit(1)
	}

	if err := confbox.Load(*confPath, &config); err != nil {
		log.Fatalln(err)
	}
}

func main() {
	for _, pointList := range config.PointLists {
		go proxyPort(pointList)
	}
	configJSON, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatalln("Internal error ", err)
	}
	log.Println("Current config:")
	log.Println(string(configJSON))
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}

func proxyPort(list data.PointList) {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	points := list.Points

	proxyMap := make(map[data.DstUrl]*httputil.ReverseProxy)
	for _, dst := range points {
		parse, _ := url.Parse(string(dst))
		reverseProxy := httputil.NewSingleHostReverseProxy(parse)
		proxyMap[dst] = reverseProxy
	}

	err := http3.ListenAndServeTLS(fmt.Sprintf(":%d", list.Port), config.CertFile, config.KeyFile, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			srcUrl = data.SrcUrl(r.Host)
			dstUrl = points[srcUrl]
		)
		r.Host = string(dstUrl)
		reverseServer, exists := proxyMap[dstUrl]
		if exists {
			chain := new(httpcompress.Chain)

			brotliNode := new(httpcompress.BrotliNode)
			zlibNode := new(httpcompress.ZlibNode)
			gzipNode := new(httpcompress.GZipNode)
			brotliNode.BindChain(chain)
			zlibNode.BindChain(chain)
			gzipNode.BindChain(chain)
			chain.Init(brotliNode, zlibNode, gzipNode, httpcompress.FinalNode(func(w http.ResponseWriter, r *http.Request) {
				reverseServer := chain.GetLocal("reverseServer").(*httputil.ReverseProxy)
				reverseServer.ServeHTTP(w, r)
			}))
			chain.PutLocal("hit", false)
			chain.PutLocal("reverseServer", reverseServer)
			brotliNode.Process(w, r)
		}
	}))
	log.Fatalf("proxy port %d fail: %v", list.Port, err)
}
