package face

import (
	"fmt"
	"gocv.io/x/gocv"
	"image"
	// "io"
	"errors"
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

// Find max face rectangle
func (fd *FaceDetector) DetectFace(img gocv.Mat) (image.Rectangle, error) {
	if img.Empty() {
		log.Println("true == img.Empty()")
		return image.Rectangle{}, errors.New("Image is empty")
	}

	rects := fd.Classifier.DetectMultiScale(img)
	fmt.Printf("found %d faces\n", len(rects))
	if len(rects) == 0 {
		return image.Rectangle{}, errors.New("no faces was found")
	}

	var maxRect image.Rectangle
	for _, r := range rects {
		if RectArea(r) > RectArea(maxRect) {
			maxRect = r
		}
	}

	fmt.Println("max rect: ", maxRect)
	return maxRect, nil
}

func (fd *FaceDetector) DetectFaceFromUrl(url, suffix string) (image.Rectangle, error) {
	prefix := "output/images/"
	err := os.MkdirAll(prefix, os.ModePerm)
	if err != nil {
		log.Println("error", err)
		return image.Rectangle{}, err
	}

	path := prefix + suffix + ".jpg"

	resp, err := http.Get(url)
	if err != nil {
		log.Println("error", err)
		return image.Rectangle{}, err
	}
	defer resp.Body.Close()

	bd, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("error", err)
		return image.Rectangle{}, err
	}

	if !IsJpg(bd, path) {
		return image.Rectangle{}, errors.New("not jpeg")
	}

	err = ioutil.WriteFile(path, bd, os.ModePerm)
	if err != nil {
		log.Println(err)
		return image.Rectangle{}, err
	}
	// file.Close()

	mat := gocv.IMRead(path, gocv.IMReadUnchanged)

	return fd.DetectFace(mat)
}

func (fd *FaceDetector) Drop() {
	fd.Classifier.Close()
}

func IsJpg(bd []byte, path string) bool {
	kind, unknown := filetype.Match(bd)
	if unknown != nil {
		fmt.Printf("\tUnknown: %s", unknown)
	}

	fmt.Printf("File type: %s. MIME: %s\n", kind.Extension, kind.MIME.Value)
	if kind.Extension == "png" {
		log.Println("Skip image", path)
		return false
	}

	return true
}
