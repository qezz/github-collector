package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"github.com/qezz/github-collector/face"
	// "gocv.io/x/gocv"
	"github.com/globalsign/mgo"
	"golang.org/x/oauth2"
	"log"
	"os"
	"time"
	// "runtime"

	"github.com/qezz/github-collector/models"
)

const (
	githubTokenEnvVar = "GITHUB_TOKEN"
)

func main() {
	// runtime.LockOSThread()
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

	xmlFile := "/Users/sergey-mishin/projects/university/project/resources/cv-rs/assets/haarcascade_frontalface_default.xml"

	fd := face.NewFaceDetector(xmlFile)
	defer fd.Drop()

	// --- mongo

	session, err := mgo.Dial("127.0.0.1")
	if err != nil {
		panic(err)
	}

	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	c := session.DB("deepllp-debug").C("people")

	// ---

	fd.DetectFaceFromUrl("https://avatars2.githubusercontent.com/u/3346272?s=460&v=4", "meow")

	// return

	opts := &github.SearchOptions{
		Sort:        "followers",
		Order:       "desc",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	for {
		res, resp, err := client.Search.Users(ctx, "followers:>6000", opts)

		if err != nil {
			log.Fatal("Cannot search for users:", err)
		}

		fmt.Println("---", opts.Page, "/", resp.LastPage)
		fmt.Println("Total found:", res.GetTotal())
		for i, user := range res.Users {
			// fmt.Println(i, user)
			fmt.Printf("%v ", i)
			u, _, err := client.Users.Get(ctx, user.GetLogin())
			if err != nil {
				log.Fatalln("Can't get user")
			}
			uuu := models.NewUser(u.GetID(), u.GetLogin(), u.GetName(), u.GetLocation())
			fmt.Println(uuu)
			err = c.Insert(&uuu)
			if err != nil {
				log.Println("DB Wirte error:", err)
			}

			fd.DetectFaceFromUrl(user.GetAvatarURL(), u.GetLogin())
		}

		time.Sleep(1 * time.Second)

		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
}
