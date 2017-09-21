package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func buildFileList(response http.ResponseWriter, req *http.Request) {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error when walking served directory: %s\n", err)
			return nil
		}
		if !info.IsDir() {
			fmt.Fprintln(response, strings.TrimPrefix(path, servedDir))
		}
		return nil
	}
	filepath.Walk(servedDir, walkFunc)
}

func startServer(done chan bool) *http.Server {
	server := &http.Server{Addr: fmt.Sprintf("%s:%d", listenAddress, listenPort)}

	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir(servedDir))))
	http.HandleFunc("/filelist", buildFileList)
	http.HandleFunc("/shutdown", func(http.ResponseWriter, *http.Request) {
		done <- true
	})

	go func() {
		log.Printf("Serving %s\n", servedDir)
		log.Printf("Listening on %s:%d\n", listenAddress, listenPort)
		if err := server.ListenAndServe(); err != nil {
			log.Printf("Finished listening: %s\n", err)
		}
	}()

	return server
}

func main() {
	log.Println("Server starting")

	done := make(chan bool, 1)
	server := startServer(done)
	<-done

	if err := server.Shutdown(context.Background()); err != nil {
		log.Fatalf("Failed to shutdown the server: %s\n", err)
	}

	log.Println("Server exiting")
}
