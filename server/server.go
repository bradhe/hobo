package server

import (
	"encoding/json"
	"net/http"

	"github.com/bradhe/location-search/models"
	"github.com/bradhe/location-search/search"

	"github.com/gorilla/mux"
)

type Server struct {
	client *search.Client
	mux    *mux.Router
}

type NotFoundHandler struct{}

func (n NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RenderError(w, GetError("not_found", r.Method, r.URL.Path))
}

type MethodNotAllowedHandler struct{}

func (m MethodNotAllowedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	RenderError(w, GetError("method_not_allowed", r.URL.Path, r.Method))
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger.Infof("%s %s", r.Method, r.URL.Path)
	s.mux.ServeHTTP(w, r)
}

func LocationParam(r *http.Request) string {
	vals := r.URL.Query()
	return vals.Get("q")
}

type SearchResponse struct {
	Cities []models.City `json:"cities"`
}

func (s *Server) Search(w http.ResponseWriter, r *http.Request) {
	cities, err := s.client.Search(LocationParam(r))

	if err != nil {
		logger.WithError(err).Error("search failed")
		RenderError(w, GetError("internal_server_error"))
	} else if err := json.NewEncoder(w).Encode(SearchResponse{cities}); err != nil {
		panic(err)
	}
}

func (s *Server) ListenAndServe(addr string) error {
	server := http.Server{
		Addr:    addr,
		Handler: s,
	}

	return server.ListenAndServe()
}

func New(client *search.Client) *Server {
	s := new(Server)
	s.client = client

	r := mux.NewRouter()
	r.NotFoundHandler = NotFoundHandler{}
	r.MethodNotAllowedHandler = MethodNotAllowedHandler{}
	r.HandleFunc("/search", s.Search).Methods("GET")
	s.mux = r

	return s
}
