package main

import (
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

	newItem := fyne.NewMenuItem("Quit", func() {
		mainWindow.Close()
	})

	newItemSeparator := fyne.NewMenuItemSeparator()

	menuItem := fyne.NewMenu("File", fileItem, newItemSeparator, newItem)
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
		// fileWindow.CenterOnScreen()
		// window := screenshot.GetDisplayBounds(0)
		// fileWindow.Resize(fyne.NewSize(window.Bounds().Dx()/2, window.Bounds().Dy()/2))

		input := widget.NewEntry()
		input.SetPlaceHolder("example.png")

		content := container.NewVBox(input, widget.NewButton("Open", func() {
			colorImage, format, err := loadandsave.LoadImage(input.Text)
			if err != nil {
				dialog.ShowError(err, fileWindow)
			} else {
				width := colorImage.Bounds().Dx()
				height := colorImage.Bounds().Dy()
				grayImage := scaleGray(colorImage, width, height)
				_, _, _, min, max, brightness, contrast := calculate(grayImage, width, height, format)
				informationTape := information(format, width, height, min, max, brightness, contrast)
				lutGray := operations.LutGray()

				imageWindow := newWindow(application, colorImage.Bounds().Dx(), colorImage.Bounds().Dy(), input.Text)
				image := canvas.NewImageFromFile(input.Text)
				text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
				canvasText := canvas.NewText(text, color.Opaque)
				imageWindow.SetContent(container.NewBorder(nil, canvasText, nil, nil, image))

				newItem := fyne.NewMenuItem("Image Information", func() {
					dialog.ShowInformation("Information", informationTape, imageWindow)
				})

				newItem2 := fyne.NewMenuItem("Scale gray", func() {
					GrayButton(application, grayImage, lutGray, input.Text, format, informationTape)
				})

				newItem3 := fyne.NewMenuItem("Quit", func() {
					imageWindow.Close()
				})

				newItemSeparator := fyne.NewMenuItemSeparator()

				menuItem := fyne.NewMenu("Operations", newItem, newItemSeparator, newItem2, newItemSeparator, newItem3)
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

func GrayButton(application fyne.App, grayImage *image.Gray, lutGray map[int]int,
	input, format, information string) {
	width := grayImage.Bounds().Dx()
	height := grayImage.Bounds().Dy()
	window := newWindow(application, width, height, input)
	image := canvas.NewImageFromImage(grayImage)
	text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
	canvasText := canvas.NewText(text, color.Opaque)
	window.SetContent(container.NewBorder(nil, canvasText, nil, nil, image))

	newItem := fyne.NewMenuItem("Image Information", func() {
		dialog.ShowInformation("Information", information, window)
	})

	newItem2 := fyne.NewMenuItem("Histogram", func() {
		// histogram.Plote(numbersofpixel, values)
	})

	newItem3 := fyne.NewMenuItem("Cumulative histogram", func() {
		// plote(lutGray(), cumulativeHistogram(values))
	})

	newItem4 := fyne.NewMenuItem("Negative", func() {
		negativeImage := operations.Negative(grayImage, lutGray, width, height)
		NegativeButton(application, negativeImage, lutGray, format, "Negative")
	})

	newItem5 := fyne.NewMenuItem("Convert to RGB", func() {
		// loadandsave.SaveImage("RGB.png", colorImage)
	})

	newItem6 := fyne.NewMenuItem("Quit", func() {
		window.Close()
	})

	newItemSeparator := fyne.NewMenuItemSeparator()

	menuItem := fyne.NewMenu("File", saveItem(application, grayImage), newItemSeparator, newItem6)
	menuItem2 := fyne.NewMenu("Image Information", newItem)
	menuItem3 := fyne.NewMenu("Operations", newItem2, newItemSeparator, newItem3, newItemSeparator, newItem4, newItemSeparator, newItem5)
	menu := fyne.NewMainMenu(menuItem, menuItem2, menuItem3)
	window.SetMainMenu(menu)
	window.Show()
	window.Close()
}

func NegativeButton(application fyne.App, negativeImage *image.Gray,
	lutGray map[int]int, format, input string) {
	width := negativeImage.Bounds().Dx()
	height := negativeImage.Bounds().Dy()
	window := newWindow(application, width, height, input)
	image := canvas.NewImageFromImage(negativeImage)
	text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
	canvasText := canvas.NewText(text, color.Opaque)
	window.SetContent(container.NewBorder(nil, canvasText, nil, nil, image))

	_, _, _, min, max, brightness, contrast := calculate(negativeImage, width, height, format)

	informationTape := information(format, width, height, min, max, brightness, contrast)

	newItem := fyne.NewMenuItem("Image Information", func() {
		dialog.ShowInformation("Information", informationTape, window)
	})

	newItem2 := fyne.NewMenuItem("Histogram", func() {
		// histogram.Plote(numbersOfPixel, values)
	})

	newItem3 := fyne.NewMenuItem("Cumulative histogram", func() {
		// plote(lutGray(), cumulativeHistogram(values))
	})

	newItem4 := fyne.NewMenuItem("Negative", func() {
		negativeImage := operations.Negative(negativeImage, lutGray, width, height)
		NegativeButton(application, negativeImage, lutGray, format, "Negative")
	})

	newItem5 := fyne.NewMenuItem("Convert to RGB", func() {
		// loadandsave.SaveImage("RGB.png", colorImage)
	})

	newItem6 := fyne.NewMenuItem("Quit", func() {
		window.Close()
	})

	newItemSeparator := fyne.NewMenuItemSeparator()

	menuItem := fyne.NewMenu("File", saveItem(application, negativeImage), newItemSeparator, newItem6)
	menuItem2 := fyne.NewMenu("Image Information", newItem)
	menuItem3 := fyne.NewMenu("Operations", newItem2, newItemSeparator, newItem3, newItemSeparator, newItem4, newItemSeparator, newItem5)
	menu := fyne.NewMainMenu(menuItem, menuItem2, menuItem3)
	window.SetMainMenu(menu)
	window.Show()
	window.Close()
}

func saveItem(application fyne.App, image image.Image) *fyne.MenuItem {
	newItem := fyne.NewMenuItem("Save Image", func() {
		fileWindow := application.NewWindow("SaveFile")
		fileWindow.Resize(fyne.NewSize(500, 500))
		input := widget.NewEntry()
		input.SetPlaceHolder("example.png")
		content := container.NewVBox(input, widget.NewButton("Save", func() {
			err := loadandsave.SaveImage(input.Text, image)
			if err != nil {
				dialog.ShowError(err, fileWindow)
			} else {
				fileWindow.Close()
			}
		}))
		fileWindow.SetContent(content)
		fileWindow.Show()
	})
	return newItem
}

func calculate(image *image.Gray, width, height int, format string) ([]uint64,
	plotter.Values, map[int]int, int, int, int, int) {
	colors, values := operations.ColorsValues(image)
	numbersOfPixel := histogram.NumbersOfPixel(colors)
	min, max := operations.ValueRange(numbersOfPixel)
	brightness := operations.Brightness(numbersOfPixel, width*height)
	contrast := operations.Contrast(numbersOfPixel, brightness, width*height)
	return colors, values, numbersOfPixel, min, max, brightness,
		contrast
}

func information(format string, width, height, min, max, brightness, contrast int) string {
	information := "Format: " + format + "\nSize:\n - Width: " +
		strconv.Itoa(width) + "\n - Height: " + strconv.Itoa(height) +
		"\nScale Gray Range:\n - Min: " + strconv.Itoa(min) + "\n - Max: " +
		strconv.Itoa(max) + "\nBrightness: " + strconv.Itoa(brightness) +
		"\nContrast: " + strconv.Itoa(contrast)
	return information

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
