package server

import (
	"log"
	"context"
	"os/signal"
	"net"
	"time"
	"os"
	"net/http"
)

//LogFunc ...
type LogFunc func(msg string, args ...interface{})

// NewServer creates a wrapper around the given handler.
func NewServer(handler http.Handler, timeout time.Duration, signals ...os.Signal) *Server {
	s := &Server{}
	s.server = &http.Server{
		Handler: handler,
	}
	s.signals = signals
	s.Log = log.Printf
	return s
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


// ServeTLS starts a secure http server
func (s *Server) ServeTLS(l net.Listener, cf, kf string) error {
	go func() {
		if err := s.server.ServeTLS(l, cf, kf); err != nil {
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