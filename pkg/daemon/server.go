package daemon

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Bitspark/slang/pkg/env"

	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

var SlangVersion string

type Server struct {
	Host   string
	Port   int
	router *mux.Router
	ctx    *context.Context
}

func addContext(ctx context.Context, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func NewServer(ctx *context.Context, env *env.Environment) *Server {
	r := mux.NewRouter().StrictSlash(true)
	srv := &Server{env.HTTP.Address, env.HTTP.Port, r, ctx}
	srv.mountWebServices()
	return srv
}
func (s *Server) Handler() http.Handler {
	handler := cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "DELETE"},
	}).Handler(s.router)
	return addContext(*s.ctx, handler)
}

func (s *Server) mountWebServices() {
	s.AddService("/operator", DefinitionService)
	s.AddService("/run", RunnerService)
	s.AddService("/share", SharingService)
	s.AddOperatorProxy("/instance")
}

func (s *Server) AddService(pathPrefix string, services *Service) {
	r := s.router.PathPrefix(pathPrefix).Subrouter()
	for path, endpoint := range services.Routes {
		path := path
		(func(endpoint *Endpoint) {
			r.HandleFunc(path, endpoint.Handle)
		})(endpoint)
	}
}

func (s *Server) AddStaticServer(pathPrefix string, directory http.Dir) {
	r := s.router.PathPrefix(pathPrefix)
	r.Handler(http.StripPrefix(pathPrefix, http.FileServer(directory)))
}

func (s *Server) AddOperatorProxy(pathPrefix string) {
	r := s.router.PathPrefix(pathPrefix)
	r.Handler(http.StripPrefix(pathPrefix,
		r.HandlerFunc(proxyRequestToOperator).GetHandler()))
}

func (s *Server) AddRedirect(path string, redirectTo string) {
	r := s.router.Path(path)
	r.Handler(http.RedirectHandler(redirectTo, http.StatusSeeOther))
}

func (s *Server) Run() error {
	return http.ListenAndServe(fmt.Sprintf(":%d", s.Port), s.Handler())
}
