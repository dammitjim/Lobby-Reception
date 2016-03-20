package auth

import (
	"encoding/json"
	"io/ioutil"
)

var auth authJSON

// Setup loads the auth json and validates it
func Setup(filepath string) error {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, &auth)
	return err
}

// RedirectURL returns the loaded redirect_url from the auth.json file
func RedirectURL() string {
	return auth.RedirectURL
}

// ClientID returns the loaded client_id from the auth.json file
func ClientID() string {
	return auth.ClientID
}

// ClientSecret returns the loaded client_secret from the auth.json file
func ClientSecret() string {
	return auth.ClientSecret
}
