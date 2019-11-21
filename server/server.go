package server

import (
	"encoding/json"
	"net/http"

	"github.com/bradhe/location-search/models"
	"github.com/bradhe/location-search/search"
)

type Server struct {
	client *search.Client
	mux    *http.ServeMux
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	mux := http.NewServeMux()
	mux.HandleFunc("/search", s.Search)
	s.mux = mux

	return s
}
