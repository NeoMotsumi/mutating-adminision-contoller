package main

import (
	"os"

	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/middleware"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/server"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/handlers"
	"github.com/gorilla/mux"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/logger"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/env"
)


func main(){
	cfg := env.LoadApplicationConfig()
	lg := logger.NewLogger(os.Stderr, cfg.LogLevel, cfg.LogFormat)

	r := mux.NewRouter()
	handlers.RegisterMutatingWebhookHandlers(r, lg)

	rw := middleware.RequestMiddelware{
		logger = lg
	}

	r.Use(rw.RequestMiddelwareHandler)

	srv := setupServer(cfg, lg, r)

	lg.Infof("listening for requests on %s...", cfg.Address)
	srv.ListenAndServeTLS("cert.pem", "key.pem")
}

//SetupServer Creates a new http server with graceful timeout
func setupServer(cfg env.Config, lg logger.Logger, router *mux.Router) *server.Server {

	s := server.NewServer(router, cfg.GracefulTimeout, cfg.Address, os.Interrupt)
	s.Log = lg.Errorf
	
	return s
}