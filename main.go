package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

var servedDir string
var listenPort uint
var listenAddress string

func init() {
	flag.StringVar(&servedDir, "served-dir", "", "Directory to serve contents from, defaults to CWD")
	flag.UintVar(&listenPort, "listen-port", 3000, "Port to listen on, defaults to 3000")
	flag.StringVar(&listenAddress, "listen-address", "", "Address to listen on, defaults to any (0.0.0.0)")
	flag.Parse()

	if servedDir == "" {
		var err error
		if servedDir, err = os.Getwd(); err != nil {
			log.Fatalf("Could not discover current working directory: %s\n", err)
		}
	}
}

func main() {
	log.Println("server starting")

	if err := NewServer(fmt.Sprintf("%s:%d", listenAddress, listenPort), servedDir).Start(); err != nil {
		log.Fatalf("server failed: %s\n", err)
	}

	log.Println("server exiting")
}
