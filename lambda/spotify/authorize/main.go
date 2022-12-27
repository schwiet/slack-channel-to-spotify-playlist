package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/schwiet/slack-spotify/util"
	"net/url"
	"os"
)

func main() {
	lambda.Start(AuthorizeHandler)
}

func AuthorizeHandler() (response events.APIGatewayProxyResponse, err error) {
	clientId, ok := os.LookupEnv("CLIENT_ID")

	if !ok {
		response = events.APIGatewayProxyResponse{StatusCode: 400}
		return
	}

	// get a random state string. This will be returned by the Spotify API
	state := util.RandStringBytesMaskImprSrcUnsafe(64)

	cbURL := url.URL{
		Scheme: "http",
		Host:   "localhost:8008",
		Path:   "auth-callback",
	}

	spURL := AuthorizeURL(clientId, cbURL.String(), state)
	response = events.APIGatewayProxyResponse{
		StatusCode: 303,
		Headers:    map[string]string{"Location": spURL.String()},
	}
	return
}

/*
 *
 */
func AuthorizeURL(clientId, redirectUrl, state string) url.URL {
	u := url.URL{
		Scheme: "https",
		Host:   "accounts.spotify.com",
		Path:   "authorize",
	}

	v := u.Query()
	v.Set("response_type", "code")
	v.Set("client_id", clientId)
	v.Set("scope", "playlist-modify-private playlist-modify-public")
	v.Set("redirect_uri", redirectUrl)
	v.Set("state", state)

	u.RawQuery = v.Encode()
	return u
}
