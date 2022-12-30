package spotify

import (
	"encoding/base64"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

// https://developer.spotify.com/documentation/general/guides/authorization/code-flow/#request-access-token
func GetTokenRequest(
	code, redirectUrl, clientId, clientSecret string,
) (
	*http.Request, error,
) {
	// URL-encode the body that is required by the api/token spotify endpoint
	tokenRequest := url.Values{}
	tokenRequest.Set("grant_type", "authorization_code")
	tokenRequest.Set("code", code)
	tokenRequest.Set("redirect_uri", redirectUrl)
	encodedRequest := tokenRequest.Encode()

	tokenPost, err := http.NewRequest(
		"POST",
		"https://accounts.spotify.com/api/token",
		strings.NewReader(encodedRequest),
	)
	if err != nil {
		return nil, err
	}

	auth := base64.StdEncoding.EncodeToString(
		[]byte(clientId + ":" + clientSecret))

	tokenPost.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	tokenPost.Header.Add("Authorization", "Basic "+auth)
	tokenPost.Header.Add("Content-Length", strconv.Itoa(len(encodedRequest)))

	return tokenPost, nil
}
