package api

import (
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
)

var baseURL = "https://api.twitch.tv/kraken"

// Fire runs the request to the Twitch API
func Fire(url string, accessToken string) ([]byte, error) {
	var err error

	split := strings.Split(url, "/api")

	if len(split) != 2 {
		return nil, errStringSplit
	}

	path := split[1]

	err = validate(path)
	if err != nil {
		return nil, err
	}

	resp, err := http.Get(baseURL + path)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func extractPath(url string) string {
	return ""
}

func validate(path string) error {
	return nil
}

var (
	errStringSplit = errors.New("Invalid string passed")
)
