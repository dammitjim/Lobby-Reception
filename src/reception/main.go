package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type coreHandler func(http.ResponseWriter, *http.Request) (int, error)

func (fn coreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := fn(w, r); err != nil {
		fmt.Println(status)
	}
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Println("Hey")
	return http.StatusOK, nil
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/api").Handler(coreHandler(catchAllHandler))
	http.ListenAndServe("localhost:4000", r)
}
