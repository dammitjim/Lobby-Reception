package main

import (
	"errors"
	"net/http"
	"os"

	"reception/auth"
	"reception/cache"
	"reception/logging"

	"github.com/gorilla/mux"
)

func main() {

	logging.Initialise(logging.Opts{
		ServiceName:  "lobby",
		ServiceGroup: "meliora",
		Level:        os.Getenv("LOGGING_LEVEL"),
	})

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
		logging.Log().Fatal(err.Error())
	}

	err := auth.Setup(filepath)
	if err != nil {
		logging.Log().Fatal(errAuthFileDoesNotExist)
	}

	addr := os.Getenv("REDIS_ADDR")
	auth := os.Getenv("REDIS_AUTH")

	if addr == "" {
		logging.Log().Fatal("REDIS_ADDR not supplied")
		os.Exit(0)
	}

	cache.Setup(addr, auth)

	r := mux.NewRouter()
	r.PathPrefix("/api").Handler(coreHandler(catchAllHandler)).Methods("GET", "POST")

	logging.Log().Info("Reception listening on " + hostname + ":" + port)
	http.ListenAndServe(hostname+":"+port, r)

}

var (
	errAuthFileDoesNotExist = errors.New("auth.json not found in root")
)
