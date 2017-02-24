// server/server.go
// Handles web server to listen for GitHub hooks

package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/guregu/kami"
)

// Server holds all information for starting a gits web server
type Server struct {
	server *http.Server
	router *kami.Mux
}

// NewServer creates a new Server
func NewServer(r *kami.Mux) *Server {
	return &Server{
		router: r,
	}
}

// Serve starts a web server at given port
func (s *Server) Serve(p int) error {
	s.server = &http.Server{
		Addr:    ":" + strconv.Itoa(p),
		Handler: s.router,

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return s.server.ListenAndServe()
}
