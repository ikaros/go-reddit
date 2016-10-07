package main

import (
	"fmt"
	"net/http"
	"os"
	"runtime"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"

	"github.com/ikaros/go-reddit/reddit"
)

func main() {
	cfg := &clientcredentials.Config{
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		TokenURL:     "https://www.reddit.com/api/v1/access_token",
		Scopes:       []string{"read"},
	}
	rc := reddit.NewClient(cfg.Client(context.Background()))
	rc.UserAgent = reddit.UserAgent(runtime.GOOS,
		"com.github.ikaros.go-reddit", "v0.0.1", "anon")

	// This is the only way to set the User-Agent header for the
	// way how the oauth2 package aquires the access_token.
	// Aaaand reddit is VERY PICKY about the useragent. Even for a single request.
	http.DefaultClient.Transport = reddit.WrapHTTPTransport(rc.UserAgent,
		http.DefaultTransport)

	links, err := rc.Listings.ByID("t3_ffff0s")
	if err != nil {
		fmt.Println(err)
	}

	links, err = rc.Listings.ByID("t3_5572gp")
	if err != nil {
		fmt.Println(err)
	}

	for _, l := range links {
		fmt.Printf("%#v\n", l.Title)
	}
}
