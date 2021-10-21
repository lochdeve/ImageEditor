package loadandsave

import (
	"errors"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"strings"
)

func LoadImage(fileName string) (image.Image, string, error) {
	// fmt.Println("Load the image:", fileName)

	var err error
	var fimg *os.File
	var img image.Image
	format := ""

	if len(fileName) != 0 {
		fimg, err = os.Open(fileName)
		if err != nil {
			return img, format, err
		}
		img, format, err = image.Decode(fimg)
		fimg.Close()
	} else {
		err = errors.New("The image must contain extension.")
	}
	// fmt.Println("Dirección de memoria de la imagen: ", fimg)
	return img, format, err
}

func SaveImage(fileName string, img image.Image) error {
	// fmt.Println("Saving the image:", fileName)
	var err error
	var fimg *os.File
	parts := strings.Split(fileName, ".")
	if len(parts) == 2 {
		extension := parts[1]
		if extension == "jpg" || extension == "jpeg" || extension == "png" {
			fimg, err = os.Create(fileName)
			if err != nil {
				return err
			}
			err = encodeImage(extension, fimg, img)
			fimg.Close()
		} else {
			err = errors.New("Unsupported format. Supported formats: png, jpg and jpeg.")
		}
	} else {
		err = errors.New("The file has not any extension.")
	}
	return err
}

func encodeImage(extension string, fimg *os.File, img image.Image) error {
	var err error
	switch extension {
	case "jpg", "jpeg":
		err = jpeg.Encode(fimg, img, nil)
	case "png":
		err = png.Encode(fimg, img)
	}
	return err
}
