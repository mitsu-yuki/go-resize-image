package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"

	"golang.org/x/image/draw"
)

const (
	resizeDir = "resized"
)

var pwd string

func init() {
	pwd, err := os.Getwd()
	if err != nil {
		panic("error: " + err.Error())
	}
	filepath.Join(pwd, resizeDir)
	_, err = os.Stat(filepath.Join(pwd, resizeDir))
	if err != nil {
		err := os.Mkdir(filepath.Join(pwd, resizeDir), 0755)
		if err != nil {
			panic(err.Error())
		}
	}

}

func main() {
	fmt.Printf("start resize image\n")

	fileList, err := filepath.Glob("*.jpg")
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(fileList)
	for _, name := range(fileList) {
		f, err := os.Open(name)
		if err != nil {
			panic(err.Error())
		}
		defer f.Close()

		img, err := jpeg.Decode(f)
		if err != nil {
			panic(err.Error())
		}

		resizeImg := ResizeImageKeepAspect(img)
		err = SaveImage(filepath.Join(pwd, resizeDir, name), resizeImg)
		if err != nil {
			panic(err.Error())
		}
	}
}

func ResizeImageKeepAspect(img image.Image) image.Image {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	if width > height {
		return ResizeImage(img, 2100, 1400)
	} else {
		return ResizeImage(img, 1400, 2100)
	}
}

func ResizeImage(img image.Image, width, height int) image.Image {
	newImage := image.NewRGBA(image.Rect(0, 0, width, height))
	draw.CatmullRom.Scale(newImage, newImage.Bounds(), img, img.Bounds(), draw.Over, nil)

	return newImage
}

func SaveImage(path string, img image.Image) error {
	f, err := os.Create(path)
	if err != nil {
		panic(err.Error())
	}
	defer f.Close()

	err = jpeg.Encode(f, img, &jpeg.Options{
		Quality: 100,
	})
	return err
}
