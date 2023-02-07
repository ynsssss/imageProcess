package imageprocess

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	_ "image/jpeg"
	"os"
	"sync"
)

func CreateImg(pixels [][]color.Color, name, path string) error {
	rect := image.Rect(0, 0, len(pixels), len(pixels[0]))
	newImg := image.NewRGBA(rect)
	for x := 0; x < len(pixels); x++ {
		for y := 0; y < len(pixels[0]); y++ {
			q := pixels[x]

			if q == nil {
				continue
			}

			p := pixels[x][y]
			if p == nil {
				continue

			}
			original, ok := color.RGBAModel.Convert(p).(color.RGBA)
			if ok {
				newImg.Set(x, y, original)
			}
		}
	}
	_ = os.Mkdir(path, os.ModePerm)
	fg, err := os.Create(fmt.Sprintf("%s/%s.%s", path, name, "jpeg"))
	if err != nil {
		return err
	}
	jpeg.Encode(fg, newImg, nil)

	return nil
}

func ConvertImageToSlice(img image.Image) [][]color.Color {

	size := img.Bounds().Size()

	var pixels [][]color.Color

	for i := 0; i < size.X; i++ {

		var y []color.Color

		for j := 0; j < size.Y; j++ {
			y = append(y, img.At(i, j))
		}

		pixels = append(pixels, y)

	}

	return pixels
}

func OpenImage(path string) (image.Image, string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	img, format, err := image.Decode(f)
	if err != nil {
		return nil, "", err
	}

	return img, format, nil

}

func RotateImageVertically(pixels [][]color.Color) {
	for i := 0; i < len(pixels); i++ {
		tr := pixels[i]
		for j := 0; j < len(pixels[i])/2; j++ {
			tr[j], tr[len(pixels[i])-j-1] = tr[len(pixels[i])-j-1], tr[j]
		}
	}

}

func RotateImageHorizontally(pixels [][]color.Color) {
	for i := 0; i < len(pixels[0]); i++ {
		tr := pixels
		for j := 0; j < len(pixels)/2; j++ {
			tr[j][i], tr[len(pixels)-j-1][i] = tr[len(pixels)-j-1][i], tr[j][i]
		}
	}
}

func GrayscaleConversion(pixels *[][]color.Color) {
	ppixels := *pixels
	xLen := len(ppixels)
	yLen := len(ppixels[0])

	newImg := make([][]color.Color, xLen)

	for i := 0; i < xLen; i++ {
		newImg[i] = make([]color.Color, yLen)
	}

	wg := sync.WaitGroup{}
	for x := 0; x < xLen; x++ {
		for y := 0; y < yLen; y++ {

			wg.Add(1)

			go func(x, y int) {
				pixel := ppixels[x][y]
				originalColor, ok := color.RGBAModel.Convert(pixel).(color.RGBA)
				if !ok {
					return
				}

				grey := uint8(float64(originalColor.R)*0.299 + float64(originalColor.B)*0.114 + float64(originalColor.G)*0.587)

				newCol := color.RGBA{
					grey,
					grey,
					grey,
					originalColor.A,
				}

				newImg[x][y] = newCol
				wg.Done()
			}(x, y)
		}
	}
	wg.Wait()

	*pixels = newImg
}
