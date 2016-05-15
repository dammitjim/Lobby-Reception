package main

import (
	"net/http"
	"strings"

	"reception/api"
	"reception/logging"
)

type coreHandler func(http.ResponseWriter, *http.Request) (int, error)

func (fn coreHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if status, err := fn(w, r); err != nil {
		http.Error(w, err.Error(), status)
	}
}

func catchAllHandler(w http.ResponseWriter, r *http.Request) (int, error) {
	b, err := api.Fire(r, "token")
	if err != nil {
		split := strings.Split(r.URL.String(), "/api")
		path := split[1]
		ip := strings.Split(r.RemoteAddr, ":")
		logging.WithFields(logging.Fields{
			"ip":     ip[0],
			"method": r.Method,
			"path":   path,
		}).Error(err)
		return http.StatusInternalServerError, err
	}

	h := w.Header()
	w.WriteHeader(200)
	h.Set("Content-Type", "application/json; charset=utf-8")
	w.Write(b)

	return http.StatusOK, nil
}
