package polling

import (
	"context"
	"log"
	"os"

	"github.com/google/go-github/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

// Authenticate and return a Github client
func Authenticate() (*github.Client, context.Context) {
	// Load env variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		// &oauth2.Token{AccessToken: "bb8ef3ba60abf17d10002562b3bf889e3a9a2c2a"},
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)
	client := github.NewClient(tc)
	return client, ctx
}
