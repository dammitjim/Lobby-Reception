package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"reception/cache"
	"strings"
)

var baseURL = "https://api.twitch.tv/kraken"
var acceptedURLS = []string{
	"/games/top",
}

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

	switch r.Method {
	case "GET":
		return cache.Process(path, func() ([]byte, error) {
			resp, err := http.Get(baseURL + path)
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
