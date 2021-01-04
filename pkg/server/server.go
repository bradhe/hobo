package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/bradhe/hobo/pkg/config"
	"github.com/bradhe/hobo/pkg/models"
	"github.com/bradhe/hobo/pkg/search"

	"github.com/gorilla/mux"
)

type Server struct {
	search search.Search
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

func (s *Server) GetSearch(w http.ResponseWriter, r *http.Request) {
	loc := LocationParam(r)
	logger.WithField("location", loc).Info("performing search")

	cities, err := s.search.Search(loc)

	if err != nil {
		logger.WithError(err).Error("search failed")
		RenderError(w, GetError("internal_server_error"))
	} else if err := json.NewEncoder(w).Encode(SearchResponse{cities}); err != nil {
		logger.WithError(err).Error("failed to serialize results")
		RenderError(w, GetError("internal_server_error"))
	}
}

func (s *Server) GetHealthCheck(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, `{"status": "OK"}`)
}

func (s *Server) ListenAndServe(addr string) error {
	server := http.Server{
		Addr:    addr,
		Handler: s,
	}

	return server.ListenAndServe()
}

func New(conf *config.Config) *Server {
	s := &Server{
		search: search.New(conf),
	}

	r := mux.NewRouter()
	r.NotFoundHandler = NotFoundHandler{}
	r.MethodNotAllowedHandler = MethodNotAllowedHandler{}
	r.HandleFunc("/search", s.GetSearch).Methods("GET")
	r.HandleFunc("/_health/check", s.GetHealthCheck).Methods("GET")
	s.mux = r

	return s
}
