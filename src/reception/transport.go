package main

import (
	"net/http"
	"reception/api"
)

type coreHandler func(http.ResponseWriter, *http.Request) (int, error)

func (fn coreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := fn(w, r); err != nil {
		http.Error(w, err.Error(), status)
	}
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	b, err := api.Fire(r.URL.String(), "token")
	if err != nil {
		return http.StatusInternalServerError, err
	}

	h := w.Header()
	w.WriteHeader(200)
	h.Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)

	return http.StatusOK, nil
}
