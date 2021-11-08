package operations

import (
	"errors"
	"image"
	"image/color"
	"math"

	"gonum.org/v1/plot/plotter"
)

func ColorsValues(image *image.Gray) ([]uint64, plotter.Values) {
	var colors []uint64
	var values plotter.Values

	for i := 0; i < image.Bounds().Dx(); i++ {
		for j := 0; j < image.Bounds().Dy(); j++ {
			y := image.GrayAt(i, j).Y
			colors = append(colors, uint64(y))
			values = append(values, float64(y))
		}
	}
	return colors, values
}

func Negative(img *image.Gray, lutGray map[int]int, width, height int) *image.Gray {
	img2 := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			newColor := color.Gray{uint8(float32(lutGray[int(img.GrayAt(i, j).Y)]))}
			img2.Set(i, j, newColor)
		}
	}
	return img2
}

func LutGray() map[int]int {
	table := make(map[int]int)
	for i := 0; i <= 255; i++ {
		table[i] = 255 - i
	}
	return table
}

func ValueRange(histogram map[int]int) (int, int) {
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

func Brightness(numbersOfPixels map[int]int, size int) float64 {
	sumValues := 0
	for i := 0; i < len(numbersOfPixels); i++ {
		sumValues += i * numbersOfPixels[i]
	}
	// println(average / numberOfColors)
	return float64(float64(sumValues) / float64(size))
}

func Contrast(numbersOfPixels map[int]int, average float64, size int) float64 {
	calculations := 0.0
	for i := 0; i < len(numbersOfPixels); i++ {
		calculations += float64(numbersOfPixels[i]) * math.Pow(float64(float64(i)-average), 2)
	}
	contrast := math.Sqrt(calculations / float64(size))
	// println(contrast)
	return contrast
}

func AdjustBrightnessAndContrast(newBrightness, newContrast float64,
	numbersOfPixels map[int]int, images *image.Gray, size int) *image.Gray {
	brightness := Brightness(numbersOfPixels, size)
	contrast := Contrast(numbersOfPixels, brightness, size)
	img2 := image.NewGray(image.Rectangle{image.Point{0, 0},
		image.Point{images.Bounds().Dx(), images.Bounds().Dy()}})

	A := newContrast / contrast
	B := newBrightness - (A * brightness)
	newValue := 0.0
	for i := 0; i < images.Bounds().Dx(); i++ {
		for j := 0; j < images.Bounds().Dy(); j++ {
			newValue = A*float64(images.GrayAt(i, j).Y) + B
			if newValue > 255 {
				newValue = 255
			} else if newValue < 0 {
				newValue = 0
			}
			newColor := color.Gray{uint8(newValue)}
			img2.Set(i, j, newColor)
		}
	}
	return img2
}

func Entropy(numbersOfPixel map[int]int, size int) float64 {
	var entropy float64
	for i := 0; i < len(numbersOfPixel); i++ {
		if numbersOfPixel[i] != 0 {
			p := float64(float64(numbersOfPixel[i]) / float64(size))
			entropy += float64(p * math.Log2(p))
		}
	}
	entropy *= -1.0
	// fmt.Print(entropy)
	return entropy
}

func ScaleGray(img image.Image, width, height int) *image.Gray {
	img2 := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			r, g, b = r>>8, g>>8, b>>8
			y := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
			grayColor := color.Gray{uint8(y)}
			img2.Set(i, j, grayColor)
		}
	}
	return img2
}

func ImageDifference(image1 *image.Gray, image2 image.Image) (*image.Gray, error) {
	differenceImage := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{image1.Bounds().Dx(), image1.Bounds().Dy()}})
	if image1.Bounds().Dx() != image2.Bounds().Dx() || image1.Bounds().Dy() != image2.Bounds().Dy() {
		return differenceImage, errors.New("the image must contain extension")
	}

	img2 := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{image2.Bounds().Dx(), image2.Bounds().Dy()}})

	for i := 0; i < image1.Bounds().Dx(); i++ {
		for j := 0; j < image1.Bounds().Dy(); j++ {
			var grayColor color.Gray
			r, g, b, _ := image2.At(i, j).RGBA()
			if r > 255 || g > 255 || b > 255 {
				r, g, b = r>>8, g>>8, b>>8
				y := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b)
				grayColor = color.Gray{uint8(y)}
				img2.Set(i, j, grayColor)
			} else {
				if r == 0 && g == 0 && b == 0 {
					grayColor = color.Gray{uint8(0)}
					img2.Set(i, j, grayColor)
				} else {
					img2.Set(i, j, image2.At(i, j))
				}
			}
			newValue := int(math.Abs(float64(uint32(image1.GrayAt(i, j).Y) - uint32(img2.GrayAt(i, j).Y))))

			newColor := color.Gray{uint8(newValue)}

			differenceImage.Set(i, j, newColor)
		}
	}
	return differenceImage, nil
}
