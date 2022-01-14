package imagecontent

import (
	"image"
	"math"
	histogram "vpc/pkg/histogram"
)

type InformationImage struct {
	image                         *image.Gray
	format                        string
	min, max, width, height       int
	brightness, contrast, entropy float64
	allImageColors                []uint64
	histogramMap, lutGray         map[int]int
}

func New(newImage *image.Gray, newLutGray map[int]int, newFormat string) InformationImage {
	newColors, newHistogramMap, newEntropy,
		newMin, newMax, newBrightness, newContrast := calculate(newImage, 0)
	return InformationImage{image: newImage, format: newFormat, min: newMin,
		max: newMax, brightness: newBrightness, contrast: newContrast,
		entropy: newEntropy, allImageColors: newColors,
		histogramMap: newHistogramMap, lutGray: newLutGray, width: newImage.Bounds().Dx(),
		height: newImage.Bounds().Dy()}
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

func (content InformationImage) AllImageColors() []uint64 {
	return content.allImageColors
}

func (content InformationImage) HistogramMap() map[int]int {
	return content.histogramMap
}

func (content InformationImage) LutGray() map[int]int {
	return content.lutGray
}

func (content InformationImage) Width() int {
	return content.width
}

func (content InformationImage) Height() int {
	return content.height
}

func (content *InformationImage) SetCont(value int) {
	content.allImageColors, content.histogramMap, content.entropy, content.min, content.max,
		content.brightness, content.contrast = calculate(content.image, value)
}

func Brightness(histogram map[int]int, size int) float64 {
	sumValues := 0
	for i := 0; i < len(histogram); i++ {
		sumValues += i * histogram[i]
	}
	return float64(float64(sumValues) / float64(size))
}

func Contrast(histogram map[int]int, average float64, size int) float64 {
	calculations := 0.0
	for i := 0; i < len(histogram); i++ {
		calculations += float64(histogram[i]) * math.Pow(float64(float64(i)-average), 2)
	}
	contrast := math.Sqrt(calculations / float64(size))
	return contrast
}

func GetAllImageColors(image *image.Gray, value, option int) []uint64 {
	var colors []uint64
	cont := value
	for i := 0; i < image.Bounds().Dx(); i++ {
		for j := 0; j < image.Bounds().Dy(); j++ {
			y := image.GrayAt(i, j).Y
			if option == 0 {
				colors = append(colors, uint64(y))
			} else if option == 1 {
				if y != 0 {
					colors = append(colors, uint64(y))
				} else if y == 0 && cont == 0 {
					colors = append(colors, uint64(y))
				} else {
					cont--
				}
			}
		}
	}
	return colors
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

func entropy(HistogramMap map[int]int, size int) float64 {
	entropy := 0.0
	for i := 0; i < len(HistogramMap); i++ {
		if HistogramMap[i] > 0 {
			p := float64(HistogramMap[i]) / float64(size)
			entropy += p * math.Log2(p)
		}
	}
	entropy *= -1.0
	return entropy
}

func calculate(image *image.Gray, value int) ([]uint64, map[int]int, float64,
	int, int, float64, float64) {
	size := image.Bounds().Dx() * image.Bounds().Dy()
	var colors []uint64
	if value == 0 {
		colors = GetAllImageColors(image, value, 0)
	} else {
		colors = GetAllImageColors(image, value, 1)
		size -= value
	}
	HistogramMap := histogram.HistogramMap(colors)
	entropy := entropy(HistogramMap, size)
	min, max := valueRange(HistogramMap)
	brightness := Brightness(HistogramMap, size)
	contrast := Contrast(HistogramMap, brightness, size)
	return colors, HistogramMap, entropy, min, max, brightness,
		contrast
}
