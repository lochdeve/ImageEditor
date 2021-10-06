package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/nfnt/resize"
)

func main() {
	// imgName := "among.png"
	imgName := "paisaje.jpg"
	img, err := loadImage(imgName)
	check(err)

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	fmt.Printf("Width: %d\n", width)
	fmt.Printf("Height: %d\n", height)

	// Guardando imagen
	// outputImageName := "guardo.png"
	// outputImageName := "guardo.jpg"
	// err = saveImage(outputImageName, img)
	// check(err)
	//zoom(500, 200, img, "pepe.jpg")
	salida := "pep.jpg"
	scaleGray(img, width, height, salida)

}

func loadImage(fileName string) (image.Image, error) {
	fmt.Println("Load the image:", fileName)

	fimg, err := os.Open(fileName)
	check(err)

	// fmt.Println("Dirección de memoria de la imagen: ", fimg)
	defer fimg.Close()

	img, _, err := image.Decode(fimg)
	check(err)

	return img, err
}

func saveImage(fileName string, img image.Image) error {
	fmt.Println("Saving the image:", fileName)

	var err error
	var fimg *os.File
	extension := fileName[len(fileName)-3:]
	if extension == "jpg" || extension == "jpeg" || extension == "png" {
		fimg, err = os.Create(fileName)
		check(err)
	}
	defer fimg.Close()

	err = checkImgFormat(extension, fimg, img)
	return err
}

func zoom(width, height int, img image.Image, fileName string) {
	// Llame a la biblioteca de cambio de tamaño para ampliar la img
	newsize := resize.Resize(uint(width), uint(height), img, resize.Lanczos3)
	err := saveImage(fileName, newsize)
	check(err)
}

func scaleGray(img image.Image, width, height int, fileName string) {
	var colors []uint32

	img2 := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			r, g, b = r>>8, g>>8, b>>8
			//y := 0.375*float64(r) + 0.5*float64(g) + 0.125*float64(b)
			y := 0.299*float32(r) + 0.587*float32(g) + 0.114*float32(b) // our
			grayColor := color.Gray{uint8(y)}
			img2.Set(i, j, grayColor)
			colors = append(colors, uint32(y))
		}
	}
	// fmt.Println(colors)
	// fmt.Println(len(colors))

	err := saveImage(fileName, img2)
	check(err)
}

func checkImgFormat(extension string, fimg *os.File, img image.Image) error {
	var err error
	switch extension {
	case "jpg", "jpeg":
		err = jpeg.Encode(fimg, img, nil)
	case "png":
		err = png.Encode(fimg, img)
	default:
		err = errors.New("unsupported format")
	}
	return err
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
