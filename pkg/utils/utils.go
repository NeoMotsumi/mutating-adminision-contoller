package utils

import (
	"fmt"
	"encoding/json"
	"net/http"
)

//Respond formats the http reponse to json format
func Respond(wr http.ResponseWriter, status int, v interface{}) {
	wr.Header().Set("Content-type", "application/json")
	wr.WriteHeader(status)
	json.NewEncoder(wr).Encode(v)
}


//RespondErr formats application errors to internal server error.
func RespondErr(wr http.ResponseWriter, err error, status int) {
	http.Error(wr, fmt.Sprintf("Error writing response: %v", err), status)
}