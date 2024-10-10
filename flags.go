package main

import (
	"Horu/confbox"
	"Horu/internal"
	"flag"
	"log"
	"os"
)

var config internal.Config

func init() {
	confPath := flag.String("c", "/etc/horu/config.yaml", "config file path(yaml).")
	demo := flag.Bool("demo", false, "Give a demo config file.")
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

	if err := confbox.Load(*confPath, &config); err != nil {
		log.Fatalln(err)
	}
}
