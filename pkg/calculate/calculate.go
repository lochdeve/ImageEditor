package calculate

import (
	"image"
	"vpc/pkg/histogram"
	"vpc/pkg/operations"

	"gonum.org/v1/plot/plotter"
)

func Calculate(image *image.Gray) ([]uint64,
	plotter.Values, map[int]int, float64, int, int, float64, float64) {
	size := image.Bounds().Dx() * image.Bounds().Dy()
	colors, values := operations.ColorsValues(image)
	numbersOfPixel := histogram.NumbersOfPixel(colors)
	entropy := operations.Entropy(numbersOfPixel, size)
	min, max := operations.ValueRange(numbersOfPixel)
	brightness := operations.Brightness(numbersOfPixel, size)
	contrast := operations.Contrast(numbersOfPixel, brightness, size)
	return colors, values, numbersOfPixel, entropy, min, max, brightness,
		contrast
}
