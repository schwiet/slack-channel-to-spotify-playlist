package spotify

import (
	"github.com/schwiet/slack-spotify/lambda/util"
	"net/http"
	"net/url"
	"os"
)

func AuthorizeHandler(w http.ResponseWriter, r *http.Request) {
	clientId, ok := os.LookupEnv("CLIENT_ID")

	if !ok {
		http.Error(w, "no CLIENT_ID in env", 400)
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
	// http.Error(w, spURL.String(), 501)
	http.Redirect(w, r, spURL.String(), 303)
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
