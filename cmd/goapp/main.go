package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"vincehpicton/click/internal/goapp"
)

func main() {
	router := mux.NewRouter()
	goapp.RESTRun(router)

	port := 3030
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), router))

}
