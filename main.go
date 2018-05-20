package main

import (
	"context"
	"fmt"
	"github.com/google/go-github/github"
	"gocv.io/x/gocv"
	"golang.org/x/oauth2"
	"image"
	"image/color"
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
		}

		break

		// if resp.NextPage == 0 {
		// 	break
		// }
		// opts.Page = resp.NextPage
	}

	// ---

	deviceID := 0
	xmlFile := "/Users/sergey-mishin/projects/university/project/resources/cv-rs/assets/haarcascade_frontalface_default.xml"
	webcam, err := gocv.VideoCaptureDevice(deviceID)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer webcam.Close()

	window := gocv.NewWindow("Face Detect")
	defer window.Close()

	img := gocv.NewMat()
	defer img.Close()

	blue := color.RGBA{0, 0, 255, 0}

	classifier := gocv.NewCascadeClassifier()
	defer classifier.Close()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
		return
	}

	fmt.Printf("start reading camera device: %v\n", deviceID)
	for {
		if ok := webcam.Read(&img); !ok {
			fmt.Printf("cannot read device %d\n", deviceID)
			return
		}
		if img.Empty() {
			continue
		}

		// detect faces
		rects := classifier.DetectMultiScale(img)
		fmt.Printf("found %d faces\n", len(rects))

		// draw a rectangle around each face on the original image,
		// along with text identifying as "Human"
		for _, r := range rects {
			gocv.Rectangle(&img, r, blue, 3)

			size := gocv.GetTextSize("Human", gocv.FontHersheyPlain, 1.2, 2)
			pt := image.Pt(r.Min.X+(r.Min.X/2)-(size.X/2), r.Min.Y-2)
			gocv.PutText(&img, "Human", pt, gocv.FontHersheyPlain, 1.2, blue, 2)
		}

		// show the image in the window, and wait 1 millisecond
		window.IMShow(img)
		if window.WaitKey(1) >= 0 {
			break
		}
	}

}
