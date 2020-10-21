package utils

import (
	"fmt"
	"encoding/json"
	"net/http"
)

//ReadRequest reads the request body and validates, the json format.
func ReadRequest(req *http.Request, v interface{}) error {
	if err := json.NewDecoder(req.Body).Decode(v); err != nil {
		return err
	}

	return nil
}

//Respond formats the http reponse to json format
func Respond(wr http.ResponseWriter, status int, v interface{}) {
	wr.Header().Set("Content-type", "application/json")
	wr.WriteHeader(status)
	json.NewEncoder(wr).Encode(v)
}


//RespondErr formats application errors to internal server error.
func RespondErr(wr http.ResponseWriter, err error) {
	http.Error(wr, fmt.Sprintf("Error writing response: %v", err), http.StatusInternalServerError)
}