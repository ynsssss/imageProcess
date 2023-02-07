package main

import (
	"fmt"

	"github.com/ynsssss/imageProcess/cmd/imageprocess"
)

func main() {
	img, _, err := imageprocess.OpenImage("images/download.jpg")
	if err != nil {
		fmt.Println(err)
		return
	}
	imgSlice := imageprocess.ConvertImageToSlice(img)
	imageprocess.RotateImageHorizontally(imgSlice)
	imageprocess.RotateImageVertically(imgSlice)
	imageprocess.GrayscaleConversion(&imgSlice)
	imageprocess.CreateImg(imgSlice, "newImage", "images/")
}
