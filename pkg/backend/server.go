package backend

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/2bytes/8k/internal/config"
	"github.com/2bytes/8k/pkg/frontend"
	"github.com/2bytes/8k/pkg/storage"
	"github.com/2bytes/8k/pkg/storage/inmemory"
	"github.com/2bytes/8k/util"
)

// Server defines all of the server configuration used throughout
type Server struct {
	config  *config.Config
	Storage storage.Mediator
	UIData  *frontend.Data
}

// NewServer creates a new instance of the server when provided with a storage backend
func NewServer(storageBackend storage.Mediator) Server {

	s := Server{
		config:  config.Get(),
		Storage: storageBackend,
	}

	s.UIData = &frontend.Data{
		Title:        s.config.Title,
		AccentColour: s.config.AccentColour,
		MaxBytes:     s.config.MaxBytes,
		MaxItems:     s.config.MaxItemsStored,
		TTL:          s.config.TTL,
		BaseAddress:  s.config.FormatBaseAddress(),
		RandomPath:   util.GenerateZBase32RandomPath(s.config.PathLength),
	}

	return s
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {

	tpl, err := template.ParseFiles(*config.UIFileHTML)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to load index file template", http.StatusInternalServerError)
		return
	}

	if err := tpl.Execute(w, s.UIData); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) uploadFile(w http.ResponseWriter, r *http.Request) {

	fileName := r.URL.Path[1:]

	encodedURL, err := url.Parse(fileName)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error parsing path", http.StatusBadRequest)
		return
	}

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "error parsing body", http.StatusBadRequest)
		return
	}

	if len(data) > s.config.MaxBytes {
		fmt.Printf("Uploaded data larger than max allowed bytes (%d) got: %d\n", s.config.MaxBytes, len(data))
		http.Error(w, fmt.Sprintf("request data larger than max allowed (%d bytes)", s.config.MaxBytes), http.StatusRequestEntityTooLarge)
	}

	err = s.Storage.Store(fileName, data)

	if err != nil {
		fmt.Println(err)
	}

	switch err {
	case inmemory.ErrorStorageFull:
		http.Error(w, err.Error(), http.StatusInsufficientStorage)
		return
	case storage.ErrorZeroLengthData, storage.ErrorZeroLengthFileName:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	case util.ErrorHashingFailed:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	case inmemory.ErrorKeyUsed:
		http.Error(w, "The path you have chosen, has been travelled before. Choose another.", http.StatusBadRequest)
		return
	}

	w.Write([]byte(s.UIData.BaseAddress + encodedURL.Path))
}

func (s *Server) serveUploaded(w http.ResponseWriter, r *http.Request) {

	data, err := s.Storage.Load(r.URL.Path[1:])

	if err != nil {
		fmt.Println(err)
	}

	switch err {
	case storage.ErrorNotFound:
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	case storage.ErrorZeroLengthData, storage.ErrorZeroLengthFileName:
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Write(data)
}

// HandleRequest is the request handler and router for incoming requests
// from the frontend
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {

	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/" {
			s.serveIndex(w, r)
		} else {
			if r.URL.Path == "/favicon.ico" {
				util.GenerateFaviconPNG(w)
				return
			}
			s.serveUploaded(w, r)
		}
	case http.MethodPost:
		s.uploadFile(w, r)
	default:
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	}
}
