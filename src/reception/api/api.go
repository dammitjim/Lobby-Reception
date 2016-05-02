package api

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"reception/auth"
	"reception/cache"

	log "github.com/Sirupsen/logrus"
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

	begin := time.Now()

	defer func() {
		ip := strings.Split(r.RemoteAddr, ":")
		log.WithFields(log.Fields{
			"ip":     ip[0],
			"method": r.Method,
			"path":   path,
			"took":   time.Since(begin),
		}).Info("serve")
	}()

	switch r.Method {
	case "GET":
		return cache.Process(path, func() ([]byte, error) {
			log.WithFields(log.Fields{
				"method": r.Method,
				"path":   path,
			}).Info("regenerate")
			req, err := http.NewRequest("GET", baseURL+path, nil)
			if err != nil {
				return nil, err
			}

			authcode := r.Header.Get("Authorization")
			if authcode != "" {
				// check to see the requested url is allowed
				match := false
				for _, url := range authPaths {
					if url == path {
						match = true
						break
					}
				}

				if !match {
					return nil, fmt.Errorf("Requested url %s not allowed", path)
				}

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
	case "POST":
		if path != "/oauth2/token" {
			return nil, fmt.Errorf("Requested url %s not allowed", path)
		}

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}

		c := strings.Split(string(body), "=")
		if len(c) != 2 {
			return nil, errors.New("Strings split failed")
		}

		data := url.Values{}
		data.Set("client_id", auth.ClientID())
		data.Set("client_secret", auth.ClientSecret())
		data.Set("grant_type", "authorization_code")
		data.Set("redirect_uri", auth.RedirectURL())
		data.Set("code", c[1])

		req, err := http.NewRequest("POST", baseURL+path, bytes.NewBufferString(data.Encode()))
		if err != nil {
			return nil, err
		}

		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("Content-Length", strconv.Itoa(len(data.Encode())))
		req.Header.Set("accept", "*/*")

		resp, err := httpClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return body, nil
	}

	return nil, errors.New("Invalid method supplied, found " + r.Method)

}

func validate(path string) error {
	return nil
}

var (
	errStringSplit = errors.New("Invalid string passed")
)
