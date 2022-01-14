package operations

import (
	"image"
	"image/color"
	"math"
	"vpc/pkg/histogram"
	imagecontent "vpc/pkg/imageContent"
)

type Pair struct {
	X, Y int
}

func Negative(content imagecontent.InformationImage,
	lutGray map[int]int) *image.Gray {
	width := content.Image().Bounds().Dx()
	height := content.Image().Bounds().Dy()
	img2 := image.NewGray(image.Rectangle{image.Point{0, 0},
		image.Point{width, height}})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			newColor := color.Gray{uint8(float32(lutGray[int(content.Image().GrayAt(i, j).Y)]))}
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

func AdjustBrightnessAndContrast(content imagecontent.InformationImage,
	newBrightness, newContrast float64) *image.Gray {
	width := content.Image().Bounds().Dx()
	height := content.Image().Bounds().Dy()
	brightness := content.Brigthness()
	contrast := content.Contrast()
	img2 := image.NewGray(image.Rectangle{image.Point{0, 0},
		image.Point{width, height}})

	A := newContrast / contrast
	B := newBrightness - (A * brightness)
	newValue := 0.0
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			newValue = A*float64(content.Image().GrayAt(i, j).Y) + B
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

func ScaleGray(img image.Image) *image.Gray {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	img2 := image.NewGray(image.Rectangle{image.Point{0, 0},
		image.Point{width, height}})
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

func Gamma(content imagecontent.InformationImage, gammaValue float64) *image.Gray {
	width := content.Image().Bounds().Dx()
	height := content.Image().Bounds().Dy()
	img := image.NewGray(image.Rectangle{image.Point{0, 0},
		image.Point{width, height}})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			currentColor := float64(content.Image().GrayAt(i, j).Y)
			a := currentColor / 255.0
			b := math.Pow(a, gammaValue)
			colorOut := b * 255.0
			newColor := color.Gray{uint8(colorOut)}
			img.Set(i, j, newColor)
		}
	}
	return img
}

func ImageDifference(content imagecontent.InformationImage,
	image2 image.Image) *image.Gray {
	widthImage1 := content.Image().Bounds().Dx()
	heightImage1 := content.Image().Bounds().Dy()
	differenceImage := image.NewGray(image.Rectangle{image.Point{0, 0},
		image.Point{widthImage1, heightImage1}})

	img2 := ScaleGray(image2)

	for i := 0; i < widthImage1; i++ {
		for j := 0; j < heightImage1; j++ {
			newValue :=
				math.Abs(float64(content.Image().GrayAt(i, j).Y) - float64(img2.GrayAt(i, j).Y))
			newColor := color.Gray{uint8(newValue)}
			differenceImage.Set(i, j, newColor)
		}
	}
	return differenceImage
}

func EqualizeAnImage(content imagecontent.InformationImage) *image.Gray {
	width := content.Image().Bounds().Dx()
	height := content.Image().Bounds().Dy()
	img := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	equalizeLut := histogram.Equalization(content.HistogramMap(), width*height)
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			color := color.Gray{uint8(equalizeLut[int(content.Image().GrayAt(i, j).Y)])}
			img.SetGray(i, j, color)
		}
	}
	return img
}

func ChangeMap(difference *image.Gray, img image.Image, threshold float64) image.Image {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	newImage := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			var newColor color.RGBA
			if difference.GrayAt(i, j).Y >= uint8(threshold) {
				newColor = color.RGBA{R: uint8(255), G: uint8(0), B: uint8(0), A: uint8(255)}
			} else {
				r, g, b, a := img.At(i, j).RGBA()
				r, g, b, a = r>>8, g>>8, b>>8, a>>8
				newColor = color.RGBA{R: uint8(r), G: uint8(g), B: uint8(b), A: uint8(a)}
			}
			newImage.SetRGBA(i, j, newColor)
		}
	}
	return newImage
}

func LinearAdjustmentInSections(content imagecontent.InformationImage,
	coordinates []Pair, number int) *image.Gray {
	width := content.Image().Bounds().Dx()
	height := content.Image().Bounds().Dy()
	var m, n float64
	lut := make(map[int]int)
	for i := 0; i < number; i++ {
		m = (float64(coordinates[i+1].Y) - float64(coordinates[i].Y)) /
			(float64(coordinates[i+1].X) - float64(coordinates[i].X))
		n = float64(coordinates[i].Y) - m*float64(coordinates[i].X)
		for j := coordinates[i].X; j <= coordinates[i+1].X; j++ {
			lut[j] = int(math.Round(m*float64(j) + n))
		}
	}
	img := image.NewGray(image.Rectangle{image.Point{0, 0},
		image.Point{width, height}})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			newColor := color.Gray{uint8(lut[int(content.Image().GrayAt(i, j).Y)])}
			img.Set(i, j, newColor)
		}
	}
	return img
}

func ROI(content imagecontent.InformationImage, i1, j1, i2, j2 int) *image.Gray {
	width := j2 - j1
	height := i2 - i1
	newImage := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	k := 0
	z := 0
	for i := j1; i <= j2; i++ {
		z = 0
		for j := i1; j <= i2; j++ {
			newImage.Set(k, z, content.Image().GrayAt(i, j))
			z++
		}
		k++
	}
	return newImage
}

func HistogramSpecification(referenceImage imagecontent.InformationImage,
	originalImage imagecontent.InformationImage) imagecontent.InformationImage {
	widthOriginal := originalImage.Image().Bounds().Dx()
	heightOriginal := originalImage.Image().Bounds().Dy()
	newImage := image.NewGray(image.Rectangle{image.Point{0, 0},
		image.Point{widthOriginal, heightOriginal}})
	referenceHistogram := histogram.CumulativeHistogram(referenceImage.HistogramMap())
	originalHistogram := histogram.CumulativeHistogram(originalImage.HistogramMap())
	referenceProbability, originalProbability := make(map[int]float64), make(map[int]float64)
	lutMap := make(map[int]int)

	for i := 0; i < 256; i++ {
		originalProbability[i] =
			float64(float64(originalHistogram[i]) / float64(len(originalImage.AllImageColors())))
		referenceProbability[i] =
			float64(float64(referenceHistogram[i]) / float64(len(referenceImage.AllImageColors())))
	}

	for i := 0; i < len(originalProbability)-1; i++ {
		j := 255
		for j >= 0 && originalProbability[i] <= referenceProbability[j] {
			lutMap[i] = j
			j = j - 1
		}
	}

	for i := 0; i < widthOriginal; i++ {
		for j := 0; j < heightOriginal; j++ {
			newImage.Set(i, j,
				color.Gray{uint8(lutMap[int(originalImage.Image().GrayAt(i, j).Y)])})
		}
	}
	return imagecontent.New(newImage, originalImage.LutGray(), originalImage.Format())
}

func RotateImg(img imagecontent.InformationImage, option int) imagecontent.InformationImage {
	if option != 0 {
		witdh := img.Image().Bounds().Dx()
		height := img.Image().Bounds().Dy()
		option--
		newImage := image.NewGray(image.Rectangle{image.Point{0, 0},
			image.Point{height, witdh}})

		l := newImage.Bounds().Dy()

		for i := witdh; i >= 0; i-- {
			k := 0
			for j := height; j >= 0; j-- {
				newImage.SetGray(k, l, img.Image().GrayAt(i, j))
				k++
			}
			l--
		}
		return RotateImg(imagecontent.New(newImage, img.LutGray(), img.Format()), option)
	} else {
		return img
	}
}

func Scale(fullimage imagecontent.InformationImage, option string,
	percentageWidth, percentageHeight float64) imagecontent.InformationImage {
	widthOriginal := fullimage.Image().Bounds().Dx()
	heightOriginal := fullimage.Image().Bounds().Dy()
	newWidth := int(percentageWidth * float64(widthOriginal))
	newHeight := int(percentageHeight * float64(heightOriginal))
	newImage := image.NewGray(image.Rectangle{image.Point{0, 0},
		image.Point{newWidth, newHeight}})

	for i := 0; i < newWidth; i++ {
		for j := 0; j < newHeight; j++ {
			x := float64(i) / percentageWidth
			y := float64(j) / percentageHeight
			xTrunc := int(x)
			yTrunc := int(y)
			A := Pair{X: xTrunc, Y: yTrunc + 1}
			B := Pair{X: xTrunc + 1, Y: yTrunc + 1}
			C := Pair{X: xTrunc, Y: yTrunc}
			D := Pair{X: xTrunc + 1, Y: yTrunc}
			points := [4]Pair{A, B, C, D}
			if option == "VMP" {
				var distances []float64
				for k := 0; k < len(points); k++ {
					d := math.Sqrt(math.Pow(float64(points[k].X)-x, 2) +
						math.Pow(float64(points[k].Y)-y, 2))
					distances = append(distances, d)
				}
				smallestNumber := distances[0]
				position := 0
				for i, element := range distances {
					if element < smallestNumber {
						smallestNumber = element
						position = i
					}
				}
				newImage.Set(i, j, fullimage.Image().GrayAt(points[position].X, points[position].Y))
			} else {
				grayA := float64(fullimage.Image().GrayAt(A.X, A.Y).Y)
				grayB := float64(fullimage.Image().GrayAt(B.X, B.Y).Y)
				grayC := float64(fullimage.Image().GrayAt(C.X, C.Y).Y)
				grayD := float64(fullimage.Image().GrayAt(D.X, D.Y).Y)
				p := x - float64(C.X)
				q := y - float64(C.Y)
				newValue := grayC + (grayD-grayC)*p + (grayA-grayC)*q + (grayB+grayC-grayA-grayD)*p*q
				newColor := color.Gray{uint8(newValue)}
				newImage.Set(i, j, newColor)
			}
		}
	}
	return imagecontent.New(newImage, fullimage.LutGray(), fullimage.Format())
}

func Rotate(fullimage imagecontent.InformationImage,
	angle float64, option string) imagecontent.InformationImage {
	widthOriginal := fullimage.Width()
	heightOriginal := fullimage.Height()
	angleRad := (angle * math.Pi) / float64(180)
	E := Pair{X: cal(0, 0, angleRad, 0), Y: cal(0, 0, angleRad, 1)}
	F := Pair{X: cal(0, widthOriginal, angleRad, 0), Y: cal(0, widthOriginal, angleRad, 1)}
	H := Pair{X: cal(widthOriginal, 0, angleRad, 0), Y: cal(widthOriginal, 0, angleRad, 1)}
	G := Pair{X: cal(widthOriginal, heightOriginal, angleRad, 0), Y: cal(widthOriginal,
		heightOriginal, angleRad, 1)}
	points := [4]Pair{E, F, H, G}
	maxX, minX := points[0].X, points[0].X
	maxY, minY := points[0].Y, points[0].Y
	for k := 0; k < len(points); k++ {
		if points[k].X > maxX {
			maxX = points[k].X
		}
		if points[k].Y > maxY {
			maxY = points[k].Y
		}
		if points[k].X < minX {
			minX = points[k].X
		}
		if points[k].Y < minY {
			minY = points[k].Y
		}
	}
	newWidth := int(math.Abs(float64(maxX) - float64(minX)))
	newHeight := int(math.Abs(float64(maxY) - float64(minY)))
	newImage := image.NewGray(image.Rectangle{image.Point{0, 0},
		image.Point{newWidth, newHeight}})

	cont := 0
	for i := 0; i < newWidth; i++ {
		for j := 0; j < newHeight; j++ {
			xPrima := i + minX
			yPrima := j + minY
			angleRad := (angle * math.Pi) / float64(180)
			x := float64(xPrima)*math.Cos(angleRad) + float64(yPrima)*math.Sin(angleRad)
			y := float64(-1)*float64(xPrima)*math.Sin(angleRad) + float64(yPrima)*math.Cos(angleRad)
			xTrunc := int(x)
			yTrunc := int(y)
			if xTrunc > widthOriginal || yTrunc > heightOriginal || xTrunc < 0 || yTrunc < 0 {
				cont++
			}
			A := Pair{X: xTrunc, Y: yTrunc + 1}
			B := Pair{X: xTrunc + 1, Y: yTrunc + 1}
			C := Pair{X: xTrunc, Y: yTrunc}
			D := Pair{X: xTrunc + 1, Y: yTrunc}
			points := [4]Pair{A, B, C, D}
			if option == "VMP" {
				var distances []float64
				for k := 0; k < len(points); k++ {
					d := math.Sqrt(math.Pow(float64(points[k].X)-x, 2) +
						math.Pow(float64(points[k].Y)-y, 2))
					distances = append(distances, d)
				}
				smallestNumber := distances[0]
				position := 0
				for i, element := range distances {
					if element < smallestNumber {
						smallestNumber = element
						position = i
					}
				}
				newImage.Set(i, j, fullimage.Image().GrayAt(points[position].X, points[position].Y))
			} else if option == "Bilineal" {
				grayA := float64(fullimage.Image().GrayAt(A.X, A.Y).Y)
				grayB := float64(fullimage.Image().GrayAt(B.X, B.Y).Y)
				grayC := float64(fullimage.Image().GrayAt(C.X, C.Y).Y)
				grayD := float64(fullimage.Image().GrayAt(D.X, D.Y).Y)
				p := x - float64(C.X)
				q := y - float64(C.Y)
				newValue := grayC + (grayD-grayC)*p + (grayA-grayC)*q + (grayB+grayC-grayA-grayD)*p*q
				newColor := color.Gray{uint8(newValue)}
				newImage.Set(i, j, newColor)
			} else {
				newImage.Set(i, j, fullimage.Image().GrayAt(xTrunc, yTrunc))
			}
		}
	}
	original := imagecontent.New(newImage, fullimage.LutGray(), fullimage.Format())
	original.SetCont(cont)
	return original
}

func cal(x, y int, angle float64, option int) int {
	if option == 0 {
		return int(float64(x)*math.Cos(angle) - float64(y)*math.Sin(angle))
	}
	return int(float64(x)*math.Sin(angle) + float64(y)*math.Cos(angle))
}
