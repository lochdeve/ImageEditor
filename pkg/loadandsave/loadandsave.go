package loadandsave

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
)

func LoadImage(fileName string) (image.Image, string, error) {
	fmt.Println("Load the image:", fileName)

	fimg, err := os.Open(fileName)
	Check(err)

	// fmt.Println("Dirección de memoria de la imagen: ", fimg)
	defer fimg.Close()

	img, format, err := image.Decode(fimg)
	Check(err)

	return img, format, err
}

func SaveImage(fileName string, img image.Image) error {
	fmt.Println("Saving the image:", fileName)

	var err error
	var fimg *os.File
	extension := fileName[len(fileName)-3:] // Mirar esto
	if extension == "jpg" || extension == "jpeg" || extension == "png" {
		fimg, err = os.Create(fileName)
		Check(err)
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

func Check(err error) {
	if err != nil {
		panic(err)
	}
}
