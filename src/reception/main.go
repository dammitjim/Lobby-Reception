package main

import (
	"log"
	"net/http"
	"os"
	"reception/auth"

	"github.com/gorilla/mux"
)

func main() {

	hostname := "localhost"
	port := "4000"
	filepath := "./auth.json"

	if os.Getenv("HOSTNAME") != "" {
		hostname = os.Getenv("HOSTNAME")
	}

	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	if os.Getenv("FILEPATH") != "" {
		filepath = os.Getenv("FILEPATH")
	}

	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Fatal(err.Error())
	}

	auth.Setup(filepath)

	r := mux.NewRouter()
	r.PathPrefix("/api").Handler(coreHandler(catchAllHandler))
	http.ListenAndServe(hostname+":"+port, r)
}
