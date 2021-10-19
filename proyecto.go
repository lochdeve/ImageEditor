package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
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
	fileItem2 := fyne.NewMenuItem("Save image", func() { fmt.Println("Saving the image") })

	menuItem := fyne.NewMenu("File", fileItem, fileItem2)
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
			loadandsave.Check(err)
			width := colorImage.Bounds().Dx()
			height := colorImage.Bounds().Dy()
			// fmt.Printf("Width: %d\n", width)
			// fmt.Printf("Height: %d\n", height)

			imageWindow := newWindow(application, colorImage.Bounds().Dx(), colorImage.Bounds().Dy(), input.Text)
			image := canvas.NewImageFromFile(input.Text)
			imageWindow.SetContent(image)

			// newItem := fyne.NewMenuItem("Histogram", func() { fmt.Println("Falta por hacer") })
			// newItem2 := fyne.NewMenuItem("Cumulative histogram", func() { fmt.Println("Falta por hacer") })
			newItem := fyne.NewMenuItem("Image Information", func() {
				information := "Format: " + format + "\nSize:\n\tWidth: " +
					strconv.Itoa(width) + "\n\tHeight: " + strconv.Itoa(height)
				dialog.ShowInformation("Information", information, imageWindow)
			})
			newItem2 := fyne.NewMenuItem("Scale gray", func() {
				grayImage := scaleGray(colorImage, width, height)
				GrayButton(application, grayImage, colorImage, input.Text)
			})

			// menuItem := fyne.NewMenu("Operations", newItem, newItem2, newItem3)
			menuItem := fyne.NewMenu("Operations", newItem, newItem2)
			menu := fyne.NewMainMenu(menuItem)
			imageWindow.SetMainMenu(menu)
			imageWindow.Show()
			fileWindow.Close()
		}))
		fileWindow.SetContent(content)
		fileWindow.Show()
	})
	return fileItem
}

func GrayButton(application fyne.App, grayImage *image.Gray, colorImage image.Image, input string) {
	width := grayImage.Bounds().Dx()
	height := grayImage.Bounds().Dy()
	window := newWindow(application, width, height, input)
	image := canvas.NewImageFromImage(grayImage)
	window.SetContent(image)
	colors, values := operations.ColorsValues(grayImage)
	histogramData := histogram.Histogram(colors)
	newItem := fyne.NewMenuItem("Histogram", func() {
		histogram.Plote(histogramData, values)
	})
	newItem2 := fyne.NewMenuItem("Cumulative histogram", func() {
		// plote(lutGray(), cumulativeHistogram(values))
	})
	newItem3 := fyne.NewMenuItem("Negative", func() {
		grayImage := operations.Negative(grayImage, operations.LutGray(), width, height)
		GrayButton(application, grayImage, colorImage, input)
	})

	newItem4 := fyne.NewMenuItem("Convert to RGB", func() {
		loadandsave.SaveImage("RGB.jpg", colorImage)
	})

	newItem5 := fyne.NewMenuItem("Save image", func() {
		var fimg *os.File
		fimg, _ = os.Create("hola.png")
		defer fimg.Close()
		_ = png.Encode(fimg, grayImage)
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
			//y := 0.375*float64(r) + 0.5*float64(g) + 0.125*float64(b)
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
