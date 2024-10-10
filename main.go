package main

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	for _, portEntry := range config.EntryList {
		go proxyPort(portEntry)
	}
	if configJSON, err := json.MarshalIndent(config, "", "  "); err != nil {
		log.Fatalln("Internal error ", err)
	} else {
		log.Printf("Current config: %s", string(configJSON))
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
}
