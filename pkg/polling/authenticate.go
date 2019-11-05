package polling

import (
	"context"

	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
)

// Authenticate and return a Github client
func Authenticate() (*github.Client, context.Context) {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		// &oauth2.Token{AccessToken: "bb8ef3ba60abf17d10002562b3bf889e3a9a2c2a"},
		&oauth2.Token{AccessToken: "1a44fd58ca485074e8c16ba303def5f718feef45"},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client, ctx
}
