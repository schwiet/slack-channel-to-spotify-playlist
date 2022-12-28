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

func AuthorizeHandler(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	clientId, ok := os.LookupEnv("SPOTIFY_CLIENT_ID")

	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `"No Spotify Client ID in env"`,
		}, nil
	}

	// get a random state string. This will be returned by the Spotify API
	state := util.RandStringBytesMaskImprSrcUnsafe(64)
	// TODO: may want to utilize this state by storing and retrieving in callback
	//       to verify that the callback request had a legitimate origin

	cbURL := url.URL{
		Scheme: "http",
		Host:   "localhost:8008",
		Path:   "auth-callback",
	}

	spURL := AuthorizeURL(clientId, cbURL.String(), state)
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
