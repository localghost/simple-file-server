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

	shutdown chan bool
	err      chan error
}

func NewServer(address, servedDir string) *Server {
	return &Server{
		server:    &http.Server{Addr: address},
		servedDir: servedDir,
		shutdown:  make(chan bool, 1),
		err:       make(chan error, 1),
	}
}

func (s *Server) Start() error {
	s.registerHandlers()
	return s.serve()
}

func (s* Server) registerHandlers() {
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir(servedDir))))
	http.HandleFunc("/filelist", s.filelistHandler)
	http.HandleFunc("/shutdown", s.shutdownHandler)
	http.HandleFunc("/health", func(response http.ResponseWriter, req *http.Request) {})
}

func (s* Server) serve() error {
	go func() {
		log.Printf("Serving %s\n", s.servedDir)
		log.Printf("Listening on %s\n", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil {
			log.Printf("Finished listening: %s\n", err)
			if err == http.ErrServerClosed {
				close(s.err)
			} else {
				s.err <- err
			}
		}
	}()

	select {
	case err := <- s.err:
		return err
	case <- s.shutdown:
		if err := s.server.Shutdown(context.Background()); err != nil {
			log.Printf("Error while shutting down the server: %s\n", err)
			return err
		}
	}

	return nil
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
	s.shutdown <- true
}
