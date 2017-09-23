package main

import (
	"context"
	"log"
	"net/http"
)

type server struct {
	server *http.Server

	servedDir string

	shutdown chan bool
	err      chan error
}

func NewServer(address, servedDir string) *server {
	return &server{
		server:    &http.Server{Addr: address},
		servedDir: servedDir,
		shutdown:  make(chan bool, 1),
		err:       make(chan error, 1),
	}
}

func (s *server) Start() error {
	s.registerHandlers()
	return s.serve()
}

func (s *server) registerHandlers() {
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir(s.servedDir))))
	http.Handle("/filelist", NewFileListHandler(s.servedDir))
	http.HandleFunc("/shutdown", s.shutdownHandler)
	http.HandleFunc("/health", func(response http.ResponseWriter, req *http.Request) {})
}

func (s *server) serve() error {
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
	case err := <-s.err:
		return err
	case <-s.shutdown:
		if err := s.server.Shutdown(context.Background()); err != nil {
			log.Printf("Error while shutting down the server: %s\n", err)
			return err
		}
	}

	return nil
}

func (s *server) shutdownHandler(http.ResponseWriter, *http.Request) {
	s.shutdown <- true
}
