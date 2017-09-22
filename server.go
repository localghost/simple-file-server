package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"context"
)

type Server struct {
	server *http.Server

	servedDir string

	done chan bool
}

func NewServer(address, servedDir string) *Server {
	return &Server{
		server: &http.Server{Addr: address},
		servedDir: servedDir,
		done: make(chan bool, 1),
	}
}

func (s *Server) Start() {
	s.registerHandlers()
	s.serve()
}

func (s* Server) registerHandlers() {
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir(servedDir))))
	http.HandleFunc("/filelist", s.filelistHandler)
	http.HandleFunc("/shutdown", s.shutdownHandler)
	http.HandleFunc("/health", func(response http.ResponseWriter, req *http.Request) {})
}

func (s* Server) serve() {
	go func() {
		log.Printf("Serving %s\n", s.servedDir)
		log.Printf("Listening on %s\n", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil {
			log.Printf("Finished listening: %s\n", err)
		}
	}()

	<- s.done

	if err := s.server.Shutdown(context.Background()); err != nil {
		log.Printf("Error while shutting down the server: %s\n", err)
	}
}

func (s* Server) filelistHandler(response http.ResponseWriter, req *http.Request) {
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

func (s* Server) shutdownHandler(http.ResponseWriter, *http.Request) {
	s.done <- true
}
