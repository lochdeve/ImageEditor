package imagecontent

import (
	"image"
	"math"
	"vpc/pkg/histogram"

	"gonum.org/v1/plot/plotter"
)

type InformationImage struct {
	image                         *image.Gray
	format                        string
	min, max                      int
	brightness, contrast, entropy float64
	colors                        []uint64
	values                        plotter.Values
	numbersOfPixel, lutGray       map[int]int
}

func New(newImage *image.Gray, newLutGray map[int]int, newFormat string) InformationImage {
	newColors, newValues, newNumbersOfPixel, newEntropy,
		newMin, newMax, newBrightness, newContrast := calculate(newImage)
	return InformationImage{image: newImage, format: newFormat, min: newMin,
		max: newMax, brightness: newBrightness, contrast: newContrast,
		entropy: newEntropy, colors: newColors, values: newValues,
		numbersOfPixel: newNumbersOfPixel, lutGray: newLutGray}
}

func (content InformationImage) Image() *image.Gray {
	return content.image
}

func (content InformationImage) Format() string {
	return content.format
}

func (content InformationImage) Min() int {
	return content.min
}

func (content InformationImage) Max() int {
	return content.max
}

func (content InformationImage) Brigthness() float64 {
	return content.brightness
}

func (content InformationImage) Contrast() float64 {
	return content.contrast
}

func (content InformationImage) Entropy() float64 {
	return content.entropy
}

func (content InformationImage) Colors() []uint64 {
	return content.colors
}

func (content InformationImage) Values() plotter.Values {
	return content.values
}

func (content InformationImage) NumbersOfPixel() map[int]int {
	return content.numbersOfPixel
}

func (content InformationImage) LutGray() map[int]int {
	return content.lutGray
}

func Brightness(numbersOfPixels map[int]int, size int) float64 {
	sumValues := 0
	for i := 0; i < len(numbersOfPixels); i++ {
		sumValues += i * numbersOfPixels[i]
	}
	return float64(float64(sumValues) / float64(size))
}

func Contrast(numbersOfPixels map[int]int, average float64, size int) float64 {
	calculations := 0.0
	for i := 0; i < len(numbersOfPixels); i++ {
		calculations += float64(numbersOfPixels[i]) * math.Pow(float64(float64(i)-average), 2)
	}
	contrast := math.Sqrt(calculations / float64(size))
	return contrast
}

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

func entropy(numbersOfPixel map[int]int, size int) float64 {
	entropy := 0.0
	for i := 0; i < len(numbersOfPixel); i++ {
		if numbersOfPixel[i] > 0 {
			p := float64(numbersOfPixel[i]) / float64(size)
			entropy += p * math.Log2(p)
		}
	}
	entropy *= -1.0
	return entropy
}

func calculate(image *image.Gray) ([]uint64,
	plotter.Values, map[int]int, float64, int, int, float64, float64) {
	size := image.Bounds().Dx() * image.Bounds().Dy()
	colors, values := ColorsValues(image)
	numbersOfPixel := histogram.NumbersOfPixel(colors)
	entropy := entropy(numbersOfPixel, size)
	min, max := valueRange(numbersOfPixel)
	brightness := Brightness(numbersOfPixel, size)
	contrast := Contrast(numbersOfPixel, brightness, size)
	return colors, values, numbersOfPixel, entropy, min, max, brightness,
		contrast
}
