package library

import (
	"log"
	"net/http"
)

func CORSMiddleware(f AppHandler) AppHandler {

	return func(w http.ResponseWriter, r *http.Request) (int, error) {
		log.Println("CORSMiddleware")

		// call appHandler func
		return f(w, r)
	}

}
