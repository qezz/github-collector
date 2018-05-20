package main

import (
	"context"
	"github.com/google/go-github/github"
	"golang.org/x/oauth2"
	"log"
	"os"
)

func main() {
	log.Println("Start")

	ghToken := os.Getenv("GITHUB_TOKEN")
	if ghToken == "" {
		log.Fatal("You should privide GITHUB_TOKEN as environment variable.")
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	_ = github.NewClient(tc)

}
