package api

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"reception/auth"
	"reception/cache"
	"strings"
)

var baseURL = "https://api.twitch.tv/kraken"
var authPaths = []string{
	"/streams/followed",
}

var httpClient http.Client

// Fire runs the request to the Twitch API
func Fire(r *http.Request, accessToken string) ([]byte, error) {
	var err error

	split := strings.Split(r.URL.String(), "/api")

	if len(split) != 2 {
		return nil, errStringSplit
	}

	path := split[1]

	err = validate(path)
	if err != nil {
		return nil, err
	}

	defer func() {
		ip := strings.Split(r.RemoteAddr, ":")
		log.Printf("%s: %s %s", ip[0], r.Method, path)
	}()

	switch r.Method {
	case "GET":
		return cache.Process(path, func() ([]byte, error) {
			req, err := http.NewRequest("GET", baseURL+path, nil)
			if err != nil {
				return nil, err
			}

			authcode := r.Header.Get("Authorization")
			if authcode != "" {
				req.Header.Set("Authorization", authcode)
			}

			req.Header.Set("client_id", auth.ClientID())

			resp, err := httpClient.Do(req)
			if err != nil {
				return nil, err
			}

			defer resp.Body.Close()
			return ioutil.ReadAll(resp.Body)
		})
	}

	return nil, errors.New("Invalid method supplied, found " + r.Method)

}

func validate(path string) error {
	return nil
}

var (
	errStringSplit = errors.New("Invalid string passed")
)
