package loadandsave

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func LoadImage(fileName string) (image.Image, string, error) {
	fmt.Println("Load the image:", fileName)

	fimg, err := os.Open(fileName)
	check(err)

	// fmt.Println("Dirección de memoria de la imagen: ", fimg)
	defer fimg.Close()

	img, format, err := image.Decode(fimg)
	check(err)

	return img, format, err
}

func SaveImage(fileName string, img image.Image) error {
	// fmt.Println("Saving the image:", fileName)
	var err error
	var fimg *os.File
	parts := strings.Split(fileName, ".")
	extension := parts[1]
	if extension == "jpg" || extension == "jpeg" || extension == "png" {
		fimg, err = os.Create(fileName)
		check(err)
	}
	defer fimg.Close()

	err = checkImgFormat(extension, fimg, img)
	return err
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
