package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
	"path/filepath"
	"strings"

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
	origImageList := getJpegFileList(".")
	resizedImageList := getJpegFileList(resizeDir)
	unresizedImageList := difference(origImageList, resizedImageList)

	if len(unresizedImageList) != 0 {
		fmt.Printf("start resize image\n")
		for i, name := range(unresizedImageList) {
			f, err := os.Open(name)
			if err != nil {
				panic(err.Error())
			}
			defer f.Close()

			img, err := jpeg.Decode(f)
			if err != nil {
				panic(err.Error())
			}

			fmt.Printf("resize[%d/%d]: %s\n", i + 1, len(unresizedImageList), name)
			resizeImg := ResizeImageKeepAspect(img)
			err = SaveImage(filepath.Join(pwd, resizeDir, name), resizeImg)
			if err != nil {
				panic(err.Error())
			}
		}
	}
	fmt.Printf("resize completed\n")
}

func getJpegFileList(path string) []string {
	files, err := os.ReadDir(path)
	if err != nil {
		panic(err.Error())
	}

	var fileList []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(strings.ToLower(file.Name()), "jpg") {
			fileList = append(fileList, file.Name())
		}
	}
	return fileList
}

func difference(a, b []string) []string {
	mb := make(map[string]struct{}, len(b))
	for _, x := range b {
		mb[x] = struct{}{}
	}
	var diff []string
	for _, x := range a {
		if _, found := mb[x]; !found {
			diff = append(diff, x)
		}
	}
	return diff
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
