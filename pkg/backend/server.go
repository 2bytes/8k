package backend

import (
	"8192bytes/internal/flags"
	"8192bytes/pkg/frontend"
	"8192bytes/pkg/storage"
	"8192bytes/pkg/storage/inmemory"
	"8192bytes/util"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
)

// Server defines all of the server configuration used throughout
type Server struct {
	UI         frontend.UI
	PathLength int
	BindTLS    bool
	MaxBytes   int
	Storage    storage.Mediator
}

// NewServer creates a new instance of the server when provided with a storage backend
func NewServer(storageBackend storage.Mediator) Server {

	s := Server{
		UI: frontend.UI{
			Title:        *flags.PageTitle,
			AccentColour: *flags.AccentColour,
			ProtoTLS:     *flags.PublicProtoTLS || *flags.BindTLS,
			Address:      *flags.PublicAddress,
			Port:         *flags.PublicPort,
		},
		MaxBytes:   *flags.MaxBytes,
		PathLength: *flags.PathLength,
		Storage:    storageBackend,
	}

	return s
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {

	tpl, err := template.ParseFiles(*flags.UIFileHTML)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to load index file template", http.StatusInternalServerError)
		return
	}

	pd := frontend.Data{
		Title:        s.UI.Title,
		AccentColour: s.UI.AccentColour,
		Proto:        s.UI.ProtoString(),
		Address:      s.UI.Address,
		Port:         strconv.Itoa(s.UI.Port),
		MaxBytes:     s.MaxBytes,
		RandomPath:   util.GenerateZBase32RandomPath(s.PathLength),
	}

	if err := tpl.Execute(w, pd); err != nil {
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

	if len(data) > s.MaxBytes {
		fmt.Printf("Uploaded data larger than max allowed bytes (%d) got: %d\n", s.MaxBytes, len(data))
		http.Error(w, fmt.Sprintf("request data larger than max allowed (%d bytes)", s.MaxBytes), http.StatusRequestEntityTooLarge)
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

	w.Write([]byte(s.UI.BaseAddress() + encodedURL.Path))
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
