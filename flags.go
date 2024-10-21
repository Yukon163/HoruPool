package main

import (
	"Horu/confbox"
	"Horu/internal"
	"flag"
	"log"
	"os"
	"path/filepath"
	"time"
)

var config internal.HoruConfig

const HoruLogDir = "/etc/horu/log"

func init() {
	confPath := flag.String("c", "./config.yaml", "config file path(yaml).")
	demo := flag.Bool("demo", false, "Give a demo config file.")
	logToFile := flag.Bool("l", true, "Redirect log to /etc/horu/log/hh-mm.log.")
	flag.Parse()

	if *demo {
		if err := confbox.Save("./demo.yaml", internal.HoruConfig{
			EntryList: []internal.PortEntry{
				{
					Port: 80,
					Maps: map[internal.SrcConfig]internal.DstUrl{
						internal.SrcConfig{
							HTTP3:           false,
							SSL:             false,
							KeyFile:         "",
							CertFile:        "",
							SrcHost:         "blog.akua.fan",
							Allow0RTT:       false,
							EnableDatagrams: false,
						}: "http://localhost:8080",
						internal.SrcConfig{
							HTTP3:           false,
							SSL:             false,
							KeyFile:         "",
							CertFile:        "",
							SrcHost:         "ip.akua.fan",
							Allow0RTT:       false,
							EnableDatagrams: false,
						}: "http://localhost:5001",
					},
				},
				{
					Port: 443,
					Maps: map[internal.SrcConfig]internal.DstUrl{
						internal.SrcConfig{
							HTTP3:           true,
							SSL:             true,
							KeyFile:         "path/to/your/server.key",
							CertFile:        "path/to/your/server.crt",
							SrcHost:         "blog.akua.fan",
							Allow0RTT:       true,
							EnableDatagrams: true,
						}: "http://localhost:8080",
						internal.SrcConfig{
							HTTP3:           false,
							SSL:             true,
							KeyFile:         "path/to/your/server.key",
							CertFile:        "path/to/your/server.crt",
							SrcHost:         "ip.akua.fan",
							Allow0RTT:       false,
							EnableDatagrams: false,
						}: "http://localhost:5001",
					},
				},
			},
			EnableCompress: true,
			BrotliLevel:    6,
			GzipLevel:      -1,
			ZlibLevel:      -1,
		}); err != nil {
			log.Fatalln("Internal error ", err)
			return
		}
		log.Println("Demo file generated! (./demo.yaml)")
		os.Exit(1)
	}
	err := os.MkdirAll(HoruLogDir, 0755)
	if err != nil {
		log.Fatalf("error creating horu directory: %v", err)
	}
	if *logToFile {
		now := time.Now()
		logFileName := now.Format("15-04") + ".log"
		logFilePath := filepath.Join("/etc/horu/log/", logFileName)

		file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			log.Fatalf("error opening log file: %v", err)
		}

		log.SetOutput(file)
	}

	if err := confbox.Load(*confPath, &config); err != nil {
		log.Fatalln(err)
	}
}
