package server

import (
	"fmt"
	"log"
	"net/http"

	promapi "github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/yobro/introvert/web"
)

// Server server
type Server struct {
	httpServer *http.Server
	promclient *promapi.Client
}

// New returns new server
func New(port int, promaddress string) (*Server, error) {

	promclient, err := promapi.NewClient(promapi.Config{Address: promaddress})
	if err != nil {
		return nil, err
	}

	s := Server{promclient: &promclient}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/web/", http.StripPrefix("/web/", http.FileServer(web.Assets)))

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	return &s, nil
}

// Start starts server
func (s *Server) Start() error {
	log.Printf("starting server on %s\n", s.httpServer.Addr)
	return s.httpServer.ListenAndServe()
}

// Stop stops server
func (s *Server) Stop() {
	log.Printf("stopping server\n")
	s.httpServer.Close()
}
