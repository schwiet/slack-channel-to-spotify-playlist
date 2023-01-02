package spotify

import (
	"net/url"
)

func GetRedirectURI(host string) string {
	cbURL := url.URL{
		Scheme: "https",
		Host:   host,
		Path:   "authorize-callback",
	}

	return cbURL.String()
}

type Authorization struct {
	Id           string `json:"id" dynamodbav:"id"`
	AccessToken  string `json:"access_token" dynamodbav:"access_token"`
	Expiration   int    `json:"expiration" dynamodbav:"expiration"`
	RefreshToken string `json:"refresh_token" dynamodbav:"refresh_token"`
}
