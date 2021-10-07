package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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

	// outputImageName := "guardo.png"
	// outputImageName := "guardo.jpg"
	/*outputImageNameNegative := "guardo-negativo.jpg"
	img2, _ := scaleGray(img, width, height, outputImageName) // colorsGray
	err = saveImage(outputImageName, img2)
	check(err)
	tableGray := lutGray()
	img3 := negative(img2, tableGray, width, height)
	err = saveImage(outputImageNameNegative, img3)
	check(err)*/
	// fmt.Println(tableGray)
	_, colorsGray := scaleGray(img, width, height) // colorsGray
	histogram := histogram(colorsGray)
	// fmt.Println(histogram)
	min, max := valueRange(histogram)
	fmt.Printf("Rango de valores: [%d,%d]\n", min, max)
	plote(histogram)

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

func scaleGray(img image.Image, width, height int) (*image.Gray, []uint32) {
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
	return img2, colors
}

func negative(img *image.Gray, lutGray map[int]int, width, height int) *image.Gray {
	img2 := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			newColor := color.Gray{uint8(float32(lutGray[int(img.GrayAt(i, j).Y)]))}
			img2.Set(i, j, newColor)
		}
	}
	return img2
}

func lutGray() map[int]int {
	table := make(map[int]int)
	for i := 0; i <= 255; i++ {
		table[i] = 255 - i
	}
	return table
}

func histogram(colors []uint32) map[int]int {
	histogram := make(map[int]int)
	for i := 0; i <= 255; i++ {
		cont := 0
		for j := 0; j < len(colors); j++ {
			if i == int(colors[j]) {
				cont++
			}
		}
		histogram[i] = cont
	}
	// fmt.Println(histogram)
	return histogram
}

func valueRange(histogram map[int]int) (int, int) {
	// 0 Negro
	// 255 Blanco
	min := 300 // Negro
	max := 0   // Blanco
	for i := 0; i < len(histogram); i++ {
		if i >= max && histogram[i] != 0 {
			max = i
		}
		if i <= min && histogram[i] != 0 {
			min = i
		}
	}
	return min, max
}

func plote(histogram map[int]int) {

	p := plot.New()

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"
	pts := make(plotter.XYs, len(histogram))
	for i := 0; i < len(histogram); i++ {
		pts[i].X = float64(i)
		pts[i].Y = float64(histogram[i])
	}
	err := plotutil.AddLinePoints(p, "First", pts)
	check(err)

	// Save the plot to a PNG file.
	err = p.Save(4*vg.Inch, 4*vg.Inch, "points.png")
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
