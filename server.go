package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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

func (s *Server) registerHandlers() {
	http.Handle("/files/", http.StripPrefix("/files", http.FileServer(http.Dir(servedDir))))
	http.HandleFunc("/filelist", s.filelistHandler)
	http.HandleFunc("/shutdown", s.shutdownHandler)
	http.HandleFunc("/health", func(response http.ResponseWriter, req *http.Request) {})
}

func (s *Server) serve() error {
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

func (s *Server) filelistHandler(response http.ResponseWriter, req *http.Request) {
	query := req.URL.Query()
	query.Add("type", "file")
	query.Add("recursive", "yes")
	fileType := query.Get("type")
	recursive := query.Get("recursive")

	printFile := func(path string, info os.FileInfo) {
		switch fileType {
		case "any":
			fmt.Fprintln(response, strings.TrimPrefix(path, s.servedDir))
		case "file":
			if !info.IsDir() {
				fmt.Fprintln(response, strings.TrimPrefix(path, s.servedDir))
			}
		case "dir":
			if info.IsDir() {
				fmt.Fprintln(response, strings.TrimPrefix(path, s.servedDir))
			}
		}
	}

	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error when walking served directory: %s\n", err)
			return nil
		}
		printFile(path, info)
		return nil
	}

	if recursive == "yes" {
		filepath.Walk(s.servedDir, walkFunc)
	} else if recursive == "no" {
		if files, err := ioutil.ReadDir(s.servedDir); err == nil {
			for _, file := range files {
				printFile(file.Name(), file)
			}
		} else {
			response.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(response, err)
		}
	} else {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(response, "'%s' is not a valid option for 'recursive' parameter\n", recursive)
	}
}

func (s *Server) shutdownHandler(http.ResponseWriter, *http.Request) {
	s.shutdown <- true
}
