package main

import (
	"fmt"
	"github.com/pkg/errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type fileListHandler struct {
	servedDir string
}

type fileListRequest struct {
	servedDir string

	fileType   string
	recursive  bool
	startsWith string

	response http.ResponseWriter
}

func NewFileListHandler(servedDir string) http.Handler {
	return &fileListHandler{servedDir: servedDir}
}

func (h *fileListHandler) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	query.Add("type", "file")
	query.Add("recursive", "yes")

	fileType := query.Get("type")
	if err := h.checkParamater("type", fileType, "any", "file", "dir"); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(response, err)
		return
	}

	recursive := query.Get("recursive")
	if err := h.checkParamater("recursive", recursive, "yes", "no"); err != nil {
		response.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(response, err)
		return
	}

	listRequest := fileListRequest{
		servedDir:  h.servedDir,
		fileType:   fileType,
		recursive:  recursive == "yes",
		startsWith: query.Get("startswith"),
		response:   response,
	}
	listRequest.Handle()
}

func (h *fileListHandler) checkParamater(name string, value string, validValues ...string) error {
	for _, validValue := range validValues {
		if value == validValue {
			return nil
		}
	}
	return errors.Errorf("'%s' is not a valid option for '%s' parameter\n", name, value)
}

func (r *fileListRequest) Handle() {
	dir := filepath.Join(r.servedDir, r.startsWith)

	if info, err := os.Stat(dir); err != nil || !info.IsDir() {
		r.response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(r.response, "'%s' is not an existing directory\n", r.startsWith)
		return
	}

	if r.recursive {
		r.printRecursive(dir)
	} else {
		r.printFlat(dir)
	}
}

func (r *fileListRequest) printRecursive(startDir string) {
	walkFunc := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Printf("Error when walking served directory: %s\n", err)
			return nil
		}
		r.printPath(path, info)
		return nil
	}
	filepath.Walk(startDir, walkFunc)
}

func (r *fileListRequest) printFlat(startDir string) {
	if files, err := ioutil.ReadDir(startDir); err == nil {
		for _, file := range files {
			r.printPath(filepath.Join(startDir, file.Name()), file)
		}
	} else {
		r.response.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(r.response, err)
	}
}

func (r *fileListRequest) printPath(path string, info os.FileInfo) {
	switch r.fileType {
	case "any":
		r.writePath(path)
	case "file":
		if !info.IsDir() {
			r.writePath(path)
		}
	case "dir":
		if info.IsDir() {
			r.writePath(path)
		}
	}
}

func (r *fileListRequest) writePath(path string) {
	fmt.Fprintln(r.response, strings.TrimPrefix(path, r.servedDir))
}
