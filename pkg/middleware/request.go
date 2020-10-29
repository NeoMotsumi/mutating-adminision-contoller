package middleware

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/utils"
	"github.com/NeoMotsumi/mutating-adminision-contoller/pkg/logger"
)

//RequestMiddelware defines a middleware object.
type RequestMiddelware struct {
	logger logger.Logger
}

//RequestMiddelwareHandler validates the application request content
func (rq *RequestMiddelware) RequestMiddelwareHandler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			err := errors.New("Malformed request, unable to read request body")
			utils.RespondErr(w, err, http.StatusBadRequest)
			return
		}

		if len(body) == 0 {
			err := errors.New("Request body is required")
			utils.RespondErr(w, err, http.StatusBadRequest)
			return
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			rq.logger.Errorf("Content-Type=%s, expect application/json", contentType)
			http.Error(w, "invalid Content-Type, expect `application/json`", http.StatusUnsupportedMediaType)
			return
		}

        next.ServeHTTP(w, r)
    })
}