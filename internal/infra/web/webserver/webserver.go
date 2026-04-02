package webserver

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type WebServer struct {
	Router        chi.Router
	Handlers      map[string]http.HandlerFunc
	WebServerPort string
}

func NewWebServer(serverPort string) *WebServer {
	return &WebServer{
		Router:        chi.NewRouter(),
		Handlers:      make(map[string]http.HandlerFunc),
		WebServerPort: serverPort,
	}
}

func (s *WebServer) AddHandler(method, path string, handler http.HandlerFunc) {
	s.Handlers[method+":"+path] = handler
	switch method {
	case "GET":
		s.Router.Get(path, handler)
	case "POST":
		s.Router.Post(path, handler)
	case "PUT":
		s.Router.Put(path, handler)
	case "DELETE":
		s.Router.Delete(path, handler)
	}
}

func (s *WebServer) Start() {
	s.Router.Use(middleware.Logger)
	http.ListenAndServe(s.WebServerPort, s.Router)
}
