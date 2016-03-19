package main

import (
	"fmt"
	"net/http"
)

type coreHandler func(http.ResponseWriter, *http.Request) (int, error)

func (fn coreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := fn(w, r); err != nil {
		http.Error(w, err.Error(), status)
	}
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	fmt.Println("Hey")
	return http.StatusOK, nil
}
