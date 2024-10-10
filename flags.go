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

var config internal.Config

const HoruDir = "/etc/horu"
const HoruLogDir = "/etc/horu/log"

func init() {
	confPath := flag.String("c", "/etc/horu/config.yaml", "config file path(yaml).")
	demo := flag.Bool("demo", false, "Give a demo config file.")
	logToFile := flag.Bool("l", true, "Redirect log to /etc/horu/log/hh-mm.log.")
	flag.Parse()

	if *demo {
		if err := confbox.Save("./demo.yaml", internal.Config{
			EntryList: []internal.PortEntry{
				{
					Port: 8080,
					Maps: map[internal.SrcUrl]internal.DstUrl{
						"source1.com": "https://destination1.com",
						"source2.com": "https://destination2.com",
					},
				},
				{
					Port: 80,
					Maps: map[internal.SrcUrl]internal.DstUrl{
						"source3.com": "https://destination3.com",
						"source4.com": "https://destination4.com",
					},
				},
			},
			KeyFile:  "path/to/your/server.key",
			CertFile: "path/to/your/server.crt",
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
