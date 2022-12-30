package main

import (
	"context"
	"encoding/json"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/schwiet/slack-spotify/spotify"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"time"
)

var cbURL url.URL = url.URL{
	Scheme: "https",
	// TODO: get Host from request?
	Host: "localhost:8008",
	Path: "authorize-callback",
}

var AUTH_CB_URI string = cbURL.String()

func main() {
	lambda.Start(AuthorizeCallbackHandler)
}

func AuthorizeCallbackHandler(ctx context.Context, req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	clientId, cid_ok := os.LookupEnv("SPOTIFY_CLIENT_ID")
	clientSecret, cs_ok := os.LookupEnv("SPOTIFY_CLIENT_SECRET")

	if !cid_ok || !cs_ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       `"Spotify Client ID and Secret must exist in env"`,
		}, nil
	}

	authErr, ok := req.QueryStringParameters["error"]
	if ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 407,
			Body:       authErr,
		}, nil
	}

	code, ok := req.QueryStringParameters["code"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `"Must provide code and state Query Parameters"`,
		}, nil
	}

	_, ok = req.QueryStringParameters["state"]
	if !ok {
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       `"Must provide code and state Query Parameters"`,
		}, nil
	}

	// TODO: validate state

	tokenPost, err := spotify.GetTokenRequest(
		code,
		AUTH_CB_URI,
		clientId,
		clientSecret,
	)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	client := &http.Client{
		Timeout: time.Second * 10,
	}
	response, err := client.Do(tokenPost)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 502,
			Body:       err.Error(),
		}, nil
	}
	defer response.Body.Close()

	respBody, err := ioutil.ReadAll(response.Body)
	var tokenResp tokenResponse
	err = json.Unmarshal(respBody, &tokenResp)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 502,
			Body:       err.Error(),
		}, nil
	}

	// TODO: store AccessToken, RefreshToken and expiration date

	return events.APIGatewayProxyResponse{StatusCode: 200}, nil
}

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
