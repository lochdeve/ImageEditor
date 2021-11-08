package calculate

import (
	"image"
	"vpc/pkg/histogram"
	"vpc/pkg/operations"

	"gonum.org/v1/plot/plotter"
)

func Calculate(image *image.Gray, width, height int, format string) ([]uint64,
	plotter.Values, map[int]int, float64, int, int, float64, float64) {
	colors, values := operations.ColorsValues(image)
	numbersOfPixel := histogram.NumbersOfPixel(colors)
	entropy := operations.Entropy(numbersOfPixel, width*height)
	min, max := operations.ValueRange(numbersOfPixel)
	brightness := operations.Brightness(numbersOfPixel, width*height)
	contrast := operations.Contrast(numbersOfPixel, brightness, width*height)
	return colors, values, numbersOfPixel, entropy, min, max, brightness,
		contrast
}
