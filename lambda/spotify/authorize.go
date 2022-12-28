package main

import (
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/schwiet/slack-spotify/util"
	"net/url"
	"os"
)

var cbURL url.URL = url.URL{
	Scheme: "https",
	// TODO: get Host from request?
	Host: "localhost:8008",
	Path: "authorize-callback",
}

var AUTH_CB_URI string = cbURL.String()

func main() {
	lambda.Start(AuthorizeHandler)
}

func AuthorizeHandler() (events.APIGatewayProxyResponse, error) {
	clientId, ok := os.LookupEnv("SPOTIFY_CLIENT_ID")

	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `"No Spotify Client ID in env"`,
		}, nil
	}

	// get a random state string. This will be returned by the Spotify API
	// we'll use it to store the redirect_uri, since it is needed by the callback
	state := util.RandStringBytesMaskImprSrcUnsafe(64)

	// TODO: write state:expiration to database for retrieval in callback

	spURL := AuthorizeURL(clientId, AUTH_CB_URI, state)
	return events.APIGatewayProxyResponse{
		StatusCode: 303,
		Headers:    map[string]string{"Location": spURL.String()},
	}, nil
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
