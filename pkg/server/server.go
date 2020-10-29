package server

import (
	"path/filepath"
	"log"
	"context"
	"os/signal"
	"time"
	"os"
	"net/http"
)

//LogFunc ...
type LogFunc func(msg string, args ...interface{})

const (
	tlsDir      = `/etc/webhook/certs/`
)

// NewServer creates a wrapper around the given handler.
func NewServer(handler http.Handler, timeout time.Duration, addr string, signals ...os.Signal) *Server {
	s := &Server{}
	s.server = &http.Server{
		Handler: handler,
		Addr: addr,
	}
	s.signals = signals
	s.Log = log.Printf

	return s
}


func (s *Server) ListenAndServeTLS(certFile, keyFile string) {
	
	certPath := filepath.Join(tlsDir, certFile)
	keyPath := filepath.Join(tlsDir, keyFile)

	log.Fatal(s.server.ListenAndServeTLS(certPath, keyPath))
}


// Server application server wrapper.
type Server struct {
	Addr string
	Log  LogFunc
	server  *http.Server
	signals []os.Signal
	timeout time.Duration
	err     error
}

//ListenAndServe Create creates the http lister using the specified address
func (s *Server) ListenAndServe() error {
	go func() {
		s.server.Addr = s.Addr
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			s.err = err
		}
	}()
	return s.waitForShutdown()
}

func (s *Server) waitForShutdown() error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, s.signals...)
	_ = <-sig

	if s.Log != nil {
		s.Log("received interrupt. shutting down..")
	}

	ctx, cancel := context.WithTimeout(context.Background(), s.timeout)
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}


	return s.err
}