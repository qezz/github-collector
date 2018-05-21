package face

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	// "io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"

	"gopkg.in/h2non/filetype.v1"
)

type FaceDetector struct {
	Classifier gocv.CascadeClassifier
}

func NewFaceDetector(xmlFile string) FaceDetector {
	classifier := gocv.NewCascadeClassifier()

	if !classifier.Load(xmlFile) {
		fmt.Printf("Error reading cascade file: %v\n", xmlFile)
	}

	return FaceDetector{
		Classifier: classifier,
	}
}

func RectArea(r image.Rectangle) float64 {
	sideX := math.Abs(float64(r.Max.X) - float64(r.Min.X))
	sideY := math.Abs(float64(r.Max.Y) - float64(r.Min.Y))
	return sideX * sideY
}

func (fd *FaceDetector) DetectFace(img gocv.Mat) image.Rectangle {
	// detect faces

	if img.Empty() {
		log.Fatalln("true == img.Empty()")
	}

	rects := fd.Classifier.DetectMultiScale(img)
	fmt.Printf("found %d faces\n", len(rects))

	var maxRect image.Rectangle
	for _, r := range rects {
		if RectArea(r) > RectArea(maxRect) {
			maxRect = r
		}
	}

	fmt.Println("max rect: ", maxRect)
	return maxRect
}

func (fd *FaceDetector) DetectFaceFromUrl(url, suffix string) image.Rectangle {
	path := "output/images/" + suffix

	resp, err := http.Get(url)
	if err != nil {
		log.Println("error", err)
	}
	defer resp.Body.Close()

	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error", err)
	}

	kind, unknown := filetype.Match(bd)
	if unknown != nil {
		fmt.Printf("\tUnknown: %s", unknown)
	}

	fmt.Printf("File type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
	if kind.Extension == "png" {
		log.Println("Skip image", path)
		return image.Rectangle{}
	}

	// im, _, err := image.Decode(resp.Body())
	// file, err := os.Create(path)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// _, err = io.Copy(file, bd)
	err = ioutil.WriteFile(path, bd, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	// file.Close()

	mat := gocv.IMRead(path, gocv.IMReadUnchanged)

	return fd.DetectFace(mat)
}

func (fd *FaceDetector) Drop() {
	fd.Classifier.Close()
}
