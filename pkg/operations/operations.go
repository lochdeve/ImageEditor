package operations

import (
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

func Brightness(numbersOfPixels map[int]int, size int) int {
	sumValues := 0
	for i := 0; i < len(numbersOfPixels); i++ {
		sumValues += i * numbersOfPixels[i]
	}
	// println(average / numberOfColors)
	return sumValues / size
}

func Contrast(numbersOfPixels map[int]int, average, size int) int {
	calculations := 0.0
	for i := 0; i < len(numbersOfPixels); i++ {
		calculations += float64(numbersOfPixels[i]) * math.Pow(float64(i-average), 2)
	}
	contrast := int(math.Sqrt(float64(calculations) / float64(size)))
	// println(contrast)
	return contrast
}
