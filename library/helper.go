package library

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
)

type AppHandler func(w http.ResponseWriter, r *http.Request) (int, error)

func CreateHandler(f AppHandler) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {

		allowOrigin := os.Getenv("ALLOW_ORIGIN")
		if allowOrigin == "" {
			log.Println("ALLOW_ORIGIN env key is missing, fallback to *")
			log.Println("Dont use * on productions")

			allowOrigin = "*"
		}

		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		instance := os.Getenv("instance")
		log.Println("instance", instance)

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
