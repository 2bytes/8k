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
}

// NewServer creates a new instance of the server when provided with a storage backend
func NewServer(storageBackend storage.Mediator) Server {

	s := Server{
		config:  config.Get(),
		Storage: storageBackend,
	}

	return s
}

// BaseAddress returns the servers public base address as configured.
func (s *Server) BaseAddress() string {
	return s.config.FormatBaseAddress()
}

func newUIData(c *config.Config) *frontend.Data {
	ui := &frontend.Data{
		Title:        c.Title,
		AccentColour: c.AccentColour,
		MaxBytes:     c.MaxBytes,
		MaxItems:     c.MaxItemsStored,
		TTL:          c.FormattedTime(),
		BaseAddress:  c.FormatBaseAddress(),
		RandomPath:   util.GenerateZBase32RandomPath(c.PathLength),
		Version:      config.Version,
	}

	return ui
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {

	tpl, err := template.ParseFiles(*config.UIFileHTML)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to load index file template", http.StatusInternalServerError)
		return
	}

	ui := newUIData(s.config)

	for s.Storage.Contains(ui.RandomPath) {
		fmt.Printf("Path %s in use, generating another\n", ui.RandomPath)
		ui.RandomPath = util.GenerateZBase32RandomPath(s.config.PathLength)
	}

	if err := tpl.Execute(w, ui); err != nil {
		fmt.Println(err)
	}
}

func (s *Server) uploadFile(w http.ResponseWriter, r *http.Request) {

	fileName := r.URL.Path[1:]

	if fileName == "" {
		fileName = util.GenerateZBase32RandomPath(config.Get().PathLength)
	}

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
		return
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

	w.Write([]byte(s.config.FormatBaseAddress() + encodedURL.Path + "\n"))
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
				util.WriteFaviconPNG(w)
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
