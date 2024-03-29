package backend

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/2bytes/8k/internal/config"
	"github.com/2bytes/8k/pkg/frontend"
	"github.com/2bytes/8k/pkg/storage"
	"github.com/2bytes/8k/pkg/storage/inmemory"
	"github.com/2bytes/8k/util"
)

const (
	IPv4Prefix = "ip."
	IPv6Prefix = "ip6."
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

func (s *Server) newUIData(r *http.Request) *frontend.Data {
	ui := &frontend.Data{
		Title:        s.config.Title,
		AccentColour: s.config.AccentColour,
		MaxBytes:     s.config.MaxBytes,
		MaxItems:     s.config.MaxItemsStored,
		TTL:          s.config.FormattedTime(),
		BaseAddress:  util.FormatBaseAddress(s.config, r),
		RandomPath:   util.GenerateZBase32RandomPath(s.config.PathLength),
		Version:      config.Version,
	}

	return ui
}

func (s *Server) serveIndex(w http.ResponseWriter, r *http.Request) {

	var tpl *template.Template
	var err error

	if *config.UIFileHTML == "" {
		tpl, err = template.ParseFS(frontend.IndexFS, frontend.IndexFile)
	} else {
		tpl, err = template.ParseFiles(*config.UIFileHTML)
	}

	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to load index file template", http.StatusInternalServerError)
		return
	}

	ui := s.newUIData(r)

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

	w.Header().Add("Content-Type", "text/plain")
	w.Write([]byte(util.FormatBaseAddress(s.config, r) + encodedURL.Path + "\n"))
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

	w.Header().Add("Content-Type", "text/plain")
	w.Write(data)
}

// HandleRequest is the request handler and router for incoming requests
// from the frontend
func (s *Server) HandleRequest(w http.ResponseWriter, r *http.Request) {

	if strings.HasPrefix(r.Host, IPv6Prefix) {
		// TODO Format to strip port
		w.Write([]byte(r.RemoteAddr))
		return
	}

	if strings.HasPrefix(r.Host, IPv4Prefix) {
		w.Write([]byte(strings.Split(r.RemoteAddr, ":")[0]))
		return
	}

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

func (s *Server) HandleUtilityFunction(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Utility function called: %s\n", r.URL)
}
