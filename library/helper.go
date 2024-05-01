package library

import (
	"encoding/json"
	"log"
	"net/http"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) (int, error)

func CreateHandler(f AppHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		status, err := f(w, r)
		if err != nil {
			log.Println("Error:", err.Error())
			resp := NewResp(err.Error(), nil)
			WriteJson(w, status, resp)
		}
	}
}

func WriteJson(w http.ResponseWriter, status int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}
