package myapp

import (
	"net/http"

	"github.com/gorilla/mux"
)

func helloWorldHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body := []byte("Hello World")
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func RESTRun(r *mux.Router) {
	r.HandleFunc("/", helloWorldHandler)
}
