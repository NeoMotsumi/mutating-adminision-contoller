package  handlers

import (
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/mutators"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/logger"
	"github.com/gorilla/mux"
)


//RegisterMutatingWebhookHandlers registers all the webhook handlers.
func RegisterMutatingWebhookHandlers(r *mux.Router, lg logger.Logger)  {
	r.Handle("/mutate", mutators.MutatePod(lg))
}


