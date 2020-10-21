package main

import (
	"os"

	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/server"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/handlers"
	"github.com/gorilla/mux"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/logger"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/env"
)


func main(){
	cfg := env.LoadApplicationConfig()
	lg := logger.NewLogger(os.Stderr, cfg.LogLevel, cfg.LogFormat)

	router := mux.NewRouter()
	handlers.RegisterMutatingWebhookHandlers(router, lg)

	srv := setupServer(cfg, lg, router)
	
	lg.Infof("listening for requests on :%s...", cfg.Address)
	if err := srv.ListenAndServe(); err != nil {
		lg.Fatalf("http server exited: %s", err)
	}
}

//SetupServer Creates a new http server with graceful timeout
func setupServer(cfg env.Config, lg logger.Logger, router *mux.Router) *server.Server {

	s := server.NewServer(router, cfg.GracefulTimeout, os.Interrupt)
	s.Log = lg.Errorf
	s.Addr = cfg.Address
	return s
}