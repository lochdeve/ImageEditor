package main

import (
	"image"
	"image/color"
	"strconv"
	"strings"

	histogram "vpc/pkg/histogram"
	loadandsave "vpc/pkg/loadandsave"
	mouse "vpc/pkg/mouse"
	operations "vpc/pkg/operations"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
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
	mainWindow.Resize(fyne.NewSize(float32(window.Bounds().Dx()), float32(window.Bounds().Dy())))
	openFileItem := buttonOpen(application, mainWindow)

	quitItem := fyne.NewMenuItem("Quit", func() {
		mainWindow.Close()
	})

	newItemSeparator := fyne.NewMenuItemSeparator()

	fileItem := fyne.NewMenu("File", openFileItem, newItemSeparator, quitItem)
	optionItem := fyne.NewMenu("Options")
	menu := fyne.NewMainMenu(fileItem, optionItem)
	mainWindow.SetMainMenu(menu)
	mainWindow.ShowAndRun()
}

func newWindow(application fyne.App, width, height int, name string) fyne.Window {
	window := application.NewWindow(name)
	window.Resize(fyne.NewSize(float32(width), float32(height)))
	// window.Canvas().Focused().FocusGained()
	return window
}

func buttonOpen(application fyne.App, window fyne.Window) *fyne.MenuItem {
	fileItem := fyne.NewMenuItem("Open image", func() {
		fd := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				fileName := reader.URI().String()[7:]
				colorImage, format, err := loadandsave.LoadImage(fileName)
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					width := colorImage.Bounds().Dx()
					height := colorImage.Bounds().Dy()
					grayImage := operations.ScaleGray(colorImage, width, height)
					_, _, _, entropy, min, max, brightness, contrast := calculate(grayImage, width, height, format)
					informationTape := information(format, width, height, min, max, brightness, contrast, entropy)
					lutGray := operations.LutGray()

					windowName := strings.Split(fileName, "/")
					imageWindow := newWindow(application, colorImage.Bounds().Dx(), colorImage.Bounds().Dy(), windowName[len(windowName)-1])
					canvasImage := canvas.NewImageFromImage(colorImage)
					text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
					canvasText := canvas.NewText(text, color.Opaque)
					imageWindow.SetContent(container.NewBorder(nil, canvasText, nil, nil, canvasImage, mouse.New(colorImage, canvasText, text)))

					imageInformationItem := fyne.NewMenuItem("Image Information", func() {
						dialog.ShowInformation("Information", informationTape, imageWindow)
					})

					scaleGrayItem := fyne.NewMenuItem("Scale gray", func() {
						GrayButton(application, grayImage, lutGray, windowName[len(windowName)-1], format, informationTape)
					})

					quitItem := fyne.NewMenuItem("Quit", func() {
						imageWindow.Close()
					})

					separatorItem := fyne.NewMenuItemSeparator()
					saveItem := fyne.NewMenu("File", saveItem(application, grayImage), separatorItem, quitItem)

					imageInformationMenu := fyne.NewMenu("ImageInformation", imageInformationItem)
					operationItem := fyne.NewMenu("Operations", scaleGrayItem, separatorItem, quitItem)
					menu := fyne.NewMainMenu(saveItem, imageInformationMenu, operationItem)
					imageWindow.SetMainMenu(menu)
					imageWindow.Show()
				}
			}
		}, window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".png", ".jpeg"}))
		fd.Show()
	})
	return fileItem
}

func GrayButton(application fyne.App, grayImage *image.Gray, lutGray map[int]int,
	input, format, informationText string) {
	width := grayImage.Bounds().Dx()
	height := grayImage.Bounds().Dy()
	window := newWindow(application, width, height, input)
	image := canvas.NewImageFromImage(grayImage)
	text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
	canvasText := canvas.NewText(text, color.Opaque)
	window.SetContent(container.NewBorder(nil, canvasText, nil, nil, image,
		mouse.New(grayImage, canvasText, text)))

	colors, values, numbersOfPixel, _, _, _, brightness, contrast :=
		calculate(grayImage, width, height, format)

	imageInformationItem := fyne.NewMenuItem("Image Information", func() {
		dialog.ShowInformation("Information", informationText, window)
	})

	histogramItem := histogramButton(application, window, values, numbersOfPixel)

	cumulativeHistogramItem := fyne.NewMenuItem("Cumulative histogram", func() {
		histogram.Plote(numbersOfPixel, values, true)
		histogramImage, _, err := loadandsave.LoadImage(".tmp/hist.png")
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			width := histogramImage.Bounds().Dx()
			height := histogramImage.Bounds().Dy()
			windowImage := newWindow(application, width, height, "Histogram")
			text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
			canvasText := canvas.NewText(text, color.Opaque)
			image := canvas.NewImageFromImage(histogramImage)
			windowImage.SetContent(container.NewBorder(nil, canvasText, nil, nil, image))
			windowImage.Show()
		}
	})

	negativeItem := fyne.NewMenuItem("Negative", func() {
		negativeImage := operations.Negative(grayImage, lutGray, width, height)
		NegativeButton(application, negativeImage, lutGray, format, "Negative")
	})

	quitItem := fyne.NewMenuItem("Quit", func() {
		window.Close()
	})

	brightnessAndContrastItem := fyne.NewMenuItem("Brightness and Contrast", func() {
		newWindows := newWindow(application, 500, 500, "Brightness and Contrast")
		data := binding.NewFloat()
		data.Set(float64(brightness))
		slide := widget.NewSliderWithData(0, 255, data)
		formatted := binding.FloatToStringWithFormat(data, "Float value: %0.2f")
		label := widget.NewLabelWithData(formatted)

		data2 := binding.NewFloat()
		data2.Set(float64(contrast))
		slide2 := widget.NewSliderWithData(0, 127, data2)
		formatted2 := binding.FloatToStringWithFormat(data2, "Float value: %0.2f")
		label2 := widget.NewLabelWithData(formatted2)

		content := widget.NewButton("Calculate", func() {
			width := grayImage.Bounds().Dx()
			height := grayImage.Bounds().Dy()
			bright, _ := data.Get()
			conts, _ := data2.Get()
			newImage := operations.AdjustBrightnessAndContrast(int(bright), int(conts),
				numbersOfPixel, grayImage, width*height)
			_, _, _, entropy, newMin, newMax, newBrightness, newContrast :=
				calculate(newImage, width, height, format)
			newInformationText := information(format, width, height, newMin, newMax,
				newBrightness, newContrast, entropy)
			GrayButton(application, newImage, operations.LutGray(), "Modified Image",
				format, newInformationText)
		})

		brightnessText := canvas.NewText("Brightness", color.White)
		brightnessText.TextStyle = fyne.TextStyle{Bold: true}
		contrastText := canvas.NewText("Contrast", color.White)
		contrastText.TextStyle = fyne.TextStyle{Bold: true}
		menuAndImageContainer := container.NewVBox(brightnessText, label, slide,
			contrastText, label2, slide2, content)

		newWindows.SetContent(menuAndImageContainer)
		newWindows.Show()
	})

	imageDifferenceItem := differenceDialogItem(application, width, height, grayImage)

	changeMap := fyne.NewMenuItem("Change Map", func() {
		differenceDialogItem(application, width, height, grayImage)
	})

	separatorItem := fyne.NewMenuItemSeparator()

	menuItem := fyne.NewMenu("File", saveItem(application, grayImage), separatorItem, quitItem)
	menuItem2 := fyne.NewMenu("Image Information", imageInformationItem)
	menuItem3 := fyne.NewMenu("Operations", histogramItem, separatorItem, cumulativeHistogramItem, separatorItem, negativeItem, separatorItem, brightnessAndContrastItem, separatorItem, imageDifferenceItem, separatorItem, changeMap)
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
	window.SetContent(container.NewBorder(nil, canvasText, nil, nil, image, mouse.New(negativeImage, canvasText, text)))

	_, values, numbersOfPixel, entropy, min, max, brightness, contrast :=
		calculate(negativeImage, width, height, format)

	informationTape := information(format, width, height, min, max, brightness, contrast, entropy)

	informationItem := fyne.NewMenuItem("Image Information", func() {
		dialog.ShowInformation("Information", informationTape, window)
	})

	histogramItem := histogramButton(application, window, values, numbersOfPixel)

	cumulativeHistogramItem := fyne.NewMenuItem("Cumulative histogram", func() {
		histogram.Plote(numbersOfPixel, values, true)
		histogramImage, _, err := loadandsave.LoadImage(".tmp/hist.png")
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			width := histogramImage.Bounds().Dx()
			height := histogramImage.Bounds().Dy()
			windowImage := newWindow(application, width, height, "Histogram")
			text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
			canvasText := canvas.NewText(text, color.Opaque)
			image := canvas.NewImageFromImage(histogramImage)
			windowImage.SetContent(container.NewBorder(nil, canvasText, nil, nil, image))
			windowImage.Show()
		}
	})

	negativeImageItem := fyne.NewMenuItem("Negative", func() {
		negativeImage := operations.Negative(negativeImage, lutGray, width, height)
		NegativeButton(application, negativeImage, lutGray, format, "Negative")
	})

	newItem5 := fyne.NewMenuItem("Quit", func() {
		window.Close()
	})

	newItemSeparator := fyne.NewMenuItemSeparator()

	menuItem := fyne.NewMenu("File", saveItem(application, negativeImage), newItemSeparator, newItem5)
	menuItem2 := fyne.NewMenu("Image Information", informationItem)
	menuItem3 := fyne.NewMenu("Operations", histogramItem, newItemSeparator, cumulativeHistogramItem, newItemSeparator, negativeImageItem)
	menu := fyne.NewMainMenu(menuItem, menuItem2, menuItem3)
	window.SetMainMenu(menu)
	window.Show()
	window.Close()
}

func histogramButton(application fyne.App, window fyne.Window,
	values plotter.Values, numbersOfPixel map[int]int) *fyne.MenuItem {
	histogramItem := fyne.NewMenuItem("Histogram", func() {
		histogram.Plote(numbersOfPixel, values, false)
		histogramImage, _, err := loadandsave.LoadImage(".tmp/hist.png")
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			width := histogramImage.Bounds().Dx()
			height := histogramImage.Bounds().Dy()
			windowImage := newWindow(application, width, height, "Histogram")
			text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
			canvasText := canvas.NewText(text, color.Opaque)
			image := canvas.NewImageFromImage(histogramImage)
			windowImage.SetContent(container.NewBorder(nil, canvasText, nil, nil, image))
			windowImage.Show()
		}
	})
	return histogramItem
}

func saveItem(application fyne.App, image image.Image) *fyne.MenuItem {
	saveImageItem := fyne.NewMenuItem("Save Image", func() {
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
	return saveImageItem
}

func calculate(image *image.Gray, width, height int, format string) ([]uint64,
	plotter.Values, map[int]int, float64, int, int, int, int) {
	colors, values := operations.ColorsValues(image)
	numbersOfPixel := histogram.NumbersOfPixel(colors)
	entropy := operations.Entropy(numbersOfPixel, width*height)
	min, max := operations.ValueRange(numbersOfPixel)
	brightness := operations.Brightness(numbersOfPixel, width*height)
	contrast := operations.Contrast(numbersOfPixel, brightness, width*height)
	return colors, values, numbersOfPixel, entropy, min, max, brightness,
		contrast
}

func information(format string, width, height, min, max, brightness, contrast int,
	entropy float64) string {
	information := "Format: " + format + "\nSize:\n - Width: " +
		strconv.Itoa(width) + "\n - Height: " + strconv.Itoa(height) +
		"\nScale Gray Range:\n - Min: " + strconv.Itoa(min) + "\n - Max: " +
		strconv.Itoa(max) + "\nEntropy: " + strconv.Itoa(int(entropy)) +
		"\nBrightness: " + strconv.Itoa(brightness) + "\nContrast: " +
		strconv.Itoa(contrast)
	return information

}

func differenceDialogItem(application fyne.App, width, height int, grayImage *image.Gray) *fyne.MenuItem {
	dialogItem := fyne.NewMenuItem("Image difference", func() {
		windowImage := newWindow(application, width, height, "difference")
		dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {

			if reader != nil {
				fileName := reader.URI().String()[7:]
				image, _, err := loadandsave.LoadImage(fileName)
				if err != nil {
					dialog.ShowError(err, windowImage)
				}
				difference, err := operations.ImageDifference(grayImage, image)
				if err != nil {
					dialog.ShowError(err, windowImage)
				}
				canvasImage := canvas.NewImageFromImage(difference)
				windowImage.SetContent(canvasImage)
			}
		}, windowImage)

		item := buttonOpen(application, windowImage)
		menuItem := fyne.NewMenu("Operations", item)
		menu := fyne.NewMainMenu(menuItem)
		windowImage.SetMainMenu(menu)
		windowImage.Show()
		dialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".png", ".jpeg"}))

		dialog.Show()
	})
	return dialogItem
}
