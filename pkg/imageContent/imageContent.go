package imagecontent

import (
	"image"

	"gonum.org/v1/plot/plotter"
)

type InformationImage struct {
	image                         *image.Gray
	format                        string
	min, max                      int
	brightness, contrast, entropy float64
	colors                        []uint64
	values                        plotter.Values
	numbersOfPixel                map[int]int
}

func New(image *image.Gray, format string, min, max int, brightness, contrast,
	entropy float64, colors []uint64, values plotter.Values,
	numbersOfPixel map[int]int) InformationImage {
	return InformationImage{image: image, format: format, min: min, max: max,
		brightness: brightness, contrast: contrast, entropy: entropy, colors: colors,
		values: values, numbersOfPixel: numbersOfPixel}
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
