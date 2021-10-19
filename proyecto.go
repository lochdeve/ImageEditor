package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"

	histogram "vpc/pkg/histogram"
	loadandsave "vpc/pkg/loadandsave"
	operations "vpc/pkg/operations"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/widget"
	"github.com/kbinani/screenshot"
	"gonum.org/v1/plot/plotter"
)

func main() {
	interfaz()
}

func interfaz() {
	application := app.New()
	mainWindow := application.NewWindow("Hello")
	// mainWindow.Resize(fyne.NewSize(500, 500))
	window := screenshot.GetDisplayBounds(0)
	mainWindow.Resize(fyne.NewSize(window.Bounds().Dx(), window.Bounds().Dy()))

	fileItem := buttonOpen(application)

	menuItem := fyne.NewMenu("File", fileItem)
	menuItem2 := fyne.NewMenu("Options")
	menu := fyne.NewMainMenu(menuItem, menuItem2)
	mainWindow.SetMainMenu(menu)
	mainWindow.ShowAndRun()
}

func newWindow(application fyne.App, width, height int, name string) fyne.Window {
	window := application.NewWindow(name)
	window.Resize(fyne.NewSize(width, height))
	// window.Canvas().Focused().FocusGained()
	return window
}

func buttonOpen(application fyne.App) *fyne.MenuItem {
	fileItem := fyne.NewMenuItem("Open image", func() {
		fileWindow := application.NewWindow("OpenFile")
		fileWindow.Resize(fyne.NewSize(500, 500))
		window := screenshot.GetDisplayBounds(0)
		fileWindow.Resize(fyne.NewSize(window.Bounds().Dx()/2, window.Bounds().Dy()/2))

		input := widget.NewEntry()
		input.SetPlaceHolder("ejemplo.png")

		content := container.NewVBox(input, widget.NewButton("Open", func() {
			colorImage, format, err := loadandsave.LoadImage(input.Text)
			if err != nil {
				err = errors.New("No se encuentra la imagen que desea cargar.")
				dialog.ShowError(err, fileWindow)
			} else {
				width := colorImage.Bounds().Dx()
				height := colorImage.Bounds().Dy()
				grayImage := scaleGray(colorImage, width, height)
				colors, values := operations.ColorsValues(grayImage)
				numbersofpixel := histogram.NumbersOfPixel(colors)
				min, max := operations.ValueRange(numbersofpixel)
				brightness := operations.Brightness(numbersofpixel)

				imageWindow := newWindow(application, colorImage.Bounds().Dx(), colorImage.Bounds().Dy(), input.Text)
				image := canvas.NewImageFromFile(input.Text)
				imageWindow.SetContent(image)

				newItem := fyne.NewMenuItem("Image Information", func() {
					information := "Format: " + format + "\nSize:\n - Width: " +
						strconv.Itoa(width) + "\n - Height: " + strconv.Itoa(height) +
						"\nScale Gray Range:\n - Min: " + strconv.Itoa(min) + "\n - Max: " +
						strconv.Itoa(max) + "\nBrightness: " + strconv.Itoa(brightness)
					dialog.ShowInformation("Information", information, imageWindow)
				})

				newItem2 := fyne.NewMenuItem("Scale gray", func() {
					GrayButton(application, grayImage, colors, values, numbersofpixel, input.Text)
				})

				// menuItem := fyne.NewMenu("Operations", newItem, newItem2, newItem3)
				menuItem := fyne.NewMenu("Operations", newItem, newItem2)
				menu := fyne.NewMainMenu(menuItem)
				imageWindow.SetMainMenu(menu)
				imageWindow.Show()
				fileWindow.Close()
			}
		}))
		fileWindow.SetContent(content)
		fileWindow.Show()
	})
	return fileItem
}

func GrayButton(application fyne.App, grayImage *image.Gray, colors []uint64,
	values plotter.Values, numbersofpixel map[int]int, input string) {
	width := grayImage.Bounds().Dx()
	height := grayImage.Bounds().Dy()
	window := newWindow(application, width, height, input)
	image := canvas.NewImageFromImage(grayImage)
	window.SetContent(image)

	newItem := fyne.NewMenuItem("Histogram", func() {
		histogram.Plote(numbersofpixel, values)
	})

	newItem2 := fyne.NewMenuItem("Cumulative histogram", func() {
		// plote(lutGray(), cumulativeHistogram(values))
	})

	newItem3 := fyne.NewMenuItem("Negative", func() {
		// grayImage := operations.Negative(grayImage, operations.LutGray(), width, height)
		// GrayButton(application, grayImage, colorImage, input)
	})

	newItem4 := fyne.NewMenuItem("Convert to RGB", func() {
		// loadandsave.SaveImage("RGB.png", colorImage)
	})

	newItem5 := fyne.NewMenuItem("Save image", func() {
		loadandsave.SaveImage("Prueba.png", grayImage)
	})

	menuItem := fyne.NewMenu("Operations", newItem, newItem2, newItem3, newItem4)
	menuItem2 := fyne.NewMenu("Save Image", newItem5)
	menu := fyne.NewMainMenu(menuItem, menuItem2)
	window.SetMainMenu(menu)
	window.Show()
	window.Close()
}

func scaleGray(img image.Image, width, height int) *image.Gray {
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

func RGB(img image.Gray, width, height int) *image.RGBA {
	img2 := image.NewRGBA(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			r, g, b, a := img.At(i, j).RGBA()
			r, g, b = r>>8, g>>8, b>>8
			rNuevo := float64(r) / 0.299
			gNuevo := float64(g) / 0.587
			bNuevo := float64(b) / 0.114
			fmt.Println(uint(rNuevo), uint(gNuevo), uint(bNuevo))
			rgbColor := color.RGBA{uint8(uint(rNuevo) << 8), uint8(uint(gNuevo) << 8), uint8(uint(bNuevo) << 8), uint8(a)}
			img2.Set(i, j, rgbColor)
		}
	}
	return img2
}
