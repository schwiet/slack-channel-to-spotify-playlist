package spotify

import (
	"encoding/base64"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
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

func GetAuthToken(
	code, redirectUrl, clientId, clientSecret string,
) (
	*TokenResponse, AuthTokenError,
) {
	tokenPost, err := GetTokenRequest(
		code,
		redirectUrl,
		clientId,
		clientSecret,
	)
	if err != nil {
		return nil, authTokenErr(500, err.Error())
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(tokenPost)
	if err != nil {
		return nil, authTokenErr(502, err.Error())
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	var tokenResp *TokenResponse
	err = json.Unmarshal(respBody, tokenResp)
	if err != nil {
		return tokenResp, authTokenErr(502, err.Error())
	}

	return tokenResp, nil
}

type authTokenResponse struct {
	StatusCode int
	Body       string
}
type AuthTokenError *authTokenResponse

func authTokenErr(statusCode int, body string) AuthTokenError {
	return &authTokenResponse{StatusCode: statusCode, Body: body}
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
