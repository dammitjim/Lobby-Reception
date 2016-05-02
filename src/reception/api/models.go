package api

// authenticationCallback represents the data returned by
// the authentication request
type authenticationCallback struct {
	AccessToken string   `json:"access_token"`
	Scope       []string `json:"scope"`
}
