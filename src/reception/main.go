package main

import (
	"errors"
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

	if os.Getenv("HOST") != "" {
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

	err := auth.Setup(filepath)
	if err != nil {
		log.Fatal(errAuthFileDoesNotExist)
	}

	r := mux.NewRouter()
	r.PathPrefix("/api").Handler(coreHandler(catchAllHandler))

	log.Println("Reception listening on " + hostname + ":" + port)
	http.ListenAndServe(hostname+":"+port, r)

}

var (
	errAuthFileDoesNotExist = errors.New("auth.json not found in root")
)
