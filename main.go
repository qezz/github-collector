package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/qezz/github-collector/face"
	// "gocv.io/x/gocv"
	"golang.org/x/oauth2"
	"log"
	"os"
	"runtime"
)

const (
	githubTokenEnvVar = "GITHUB_TOKEN"
)

func main() {
	runtime.LockOSThread()
	log.Println("Service started")

	ghToken := os.Getenv(githubTokenEnvVar)
	if ghToken == "" {
		log.Fatal("Auth error: You should provide ", githubTokenEnvVar, " as environment variable.")
		return
	}

	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: ghToken},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	// repos, _, err := client.Repositories.List(ctx, "", nil)
	// if err != nil {
	// 	log.Fatal("Cannot get user repos")
	// }

	// for i, r := range repos {
	// 	log.Printf("%v: %v", i, r)
	// }

	/// ---

	xmlFile := "/Users/sergey-mishin/projects/university/project/resources/cv-rs/assets/haarcascade_frontalface_default.xml"

	fd := face.NewFaceDetector(xmlFile)
	defer fd.Drop()

	// fd.DetectFace(img)

	// saveFile := "output.jpg"
	// gocv.IMWrite(saveFile, img)

	/// --- search users test

	opts := &github.SearchOptions{
		Sort:        "followers",
		Order:       "desc",
		ListOptions: github.ListOptions{PerPage: 1},
	}
	// opts.Page

	for {
		res, resp, err := client.Search.Users(ctx, "followers:>6000", opts)

		if err != nil {
			log.Fatal("Cannot search for users:", err)
		}

		fmt.Println("---", opts.Page, "/", resp.LastPage)
		fmt.Println("Total found:", res.GetTotal())
		for i, user := range res.Users {
			// fmt.Println(i, "\t", *user.Login,
			// 	"\t[", user.GetLocation(),
			// 	"]\t", user.GetAvatarURL())
			fmt.Println(i, user)
			fd.DetectFaceFromUrl(user.GetAvatarURL())
		}

		break

		// if resp.NextPage == 0 {
		// 	break
		// }
		// opts.Page = resp.NextPage
	}

	// ---

}
