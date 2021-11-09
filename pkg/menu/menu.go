package menu

import (
	"image"
	"image/color"
	"strconv"
	calculate "vpc/pkg/calculate"
	histogrambutton "vpc/pkg/histogramButton"
	information "vpc/pkg/information"
	mouse "vpc/pkg/mouse"
	newwindow "vpc/pkg/newWindow"
	operations "vpc/pkg/operations"
	saveitem "vpc/pkg/saveItem"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func GeneralMenu(application fyne.App, grayImage *image.Gray, lutGray map[int]int, input, format string) {
	width := grayImage.Bounds().Dx()
	height := grayImage.Bounds().Dy()
	window := newwindow.NewWindow(application, width, height, input)
	image := canvas.NewImageFromImage(grayImage)
	text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
	canvasText := canvas.NewText(text, color.Opaque)
	window.SetContent(container.NewBorder(nil, canvasText, nil, nil, image,
		mouse.New(grayImage, canvasText, text)))

	_, values, numbersOfPixel, entropy, min, max, brightness, contrast :=
		calculate.Calculate(grayImage, width, height, format)

	informationText := information.Information(format, width, height, min, max,
		brightness, contrast, entropy)

	informationItem := fyne.NewMenuItem("Image Information", func() {
		dialog.ShowInformation("Information", informationText, window)
	})

	histogramItem := histogrambutton.HistogramButton(application, window, values, numbersOfPixel, "Histogram", false)

	cumulativeHistogramItem := histogrambutton.HistogramButton(application, window, values, numbersOfPixel, "Cumulative Histogram", true)

	negativeItem := negativeButton(application, grayImage, lutGray, width, height, format)

	gammaButton := gammaButton(application, grayImage, lutGray, input, format)

	brightnessAndContrastItem := brightnessAndContrastButton(application,
		grayImage, brightness, contrast, numbersOfPixel, lutGray, format)

	quitItem := fyne.NewMenuItem("Quit", func() {
		window.Close()
	})

	separatorItem := fyne.NewMenuItemSeparator()

	menuItem := fyne.NewMenu("File", saveitem.SaveItem(application, grayImage), separatorItem, quitItem)
	menuItem2 := fyne.NewMenu("Image Information", informationItem)
	menuItem3 := fyne.NewMenu("Operations", histogramItem, separatorItem,
		cumulativeHistogramItem, separatorItem, negativeItem, separatorItem,
		brightnessAndContrastItem, separatorItem, gammaButton) // , separatorItem, imageDifferenceItem)
	menu := fyne.NewMainMenu(menuItem, menuItem2, menuItem3)
	window.SetMainMenu(menu)
	window.Show()
	window.Close()
}

func gammaButton(application fyne.App, grayImage *image.Gray, lutGray map[int]int, input,
	format string) *fyne.MenuItem {
	newItem := fyne.NewMenuItem("Gamma", func() {
		newWindow := newwindow.NewWindow(application, 500, 150, "Gamma Value")
		data := binding.NewFloat()
		data.Set(0.05)
		slide := widget.NewSliderWithData(0.05, 20, data)
		formatted := binding.FloatToStringWithFormat(data, "Float value: %0.2f")
		label := widget.NewLabelWithData(formatted)

		content := widget.NewButton("Calculate", func() {
			width := grayImage.Bounds().Dx()
			height := grayImage.Bounds().Dy()
			gammaValue, _ := data.Get()
			img := operations.Gamma(grayImage, width, height, gammaValue)
			GeneralMenu(application, img, lutGray, "Gamma Image", format)
		})

		gammaText := canvas.NewText("Gamma", color.White)
		gammaText.TextStyle = fyne.TextStyle{Bold: true}
		menuAndImageContainer := container.NewVBox(gammaText, label, slide, content)

		newWindow.SetContent(menuAndImageContainer)
		newWindow.Show()
	})
	return newItem
}

func brightnessAndContrastButton(application fyne.App, grayImage *image.Gray,
	brightness, contrast float64, numbersOfPixel map[int]int, lutGray map[int]int,
	format string) *fyne.MenuItem {
	brightnessAndContrastItem := fyne.NewMenuItem("Brightness and Contrast", func() {
		newWindows := newwindow.NewWindow(application, 500, 500, "Brightness and Contrast")
		data := binding.NewFloat()
		data.Set(brightness)
		slide := widget.NewSliderWithData(0, 255, data)
		formatted := binding.FloatToStringWithFormat(data, "Float value: %0.2f")
		label := widget.NewLabelWithData(formatted)

		data2 := binding.NewFloat()
		data2.Set(contrast)
		slide2 := widget.NewSliderWithData(0, 127, data2)
		formatted2 := binding.FloatToStringWithFormat(data2, "Float value: %0.2f")
		label2 := widget.NewLabelWithData(formatted2)

		content := widget.NewButton("Calculate", func() {
			width := grayImage.Bounds().Dx()
			height := grayImage.Bounds().Dy()
			bright, _ := data.Get()
			conts, _ := data2.Get()
			newImage := operations.AdjustBrightnessAndContrast(bright, conts,
				numbersOfPixel, grayImage, width*height)
			GeneralMenu(application, newImage, lutGray, "Modified Image", format)
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
	return brightnessAndContrastItem
}

func negativeButton(application fyne.App, grayImage *image.Gray, lutGray map[int]int,
	width, height int, format string) *fyne.MenuItem {
	return fyne.NewMenuItem("Negative", func() {
		negativeImage := operations.Negative(grayImage, lutGray, width, height)
		GeneralMenu(application, negativeImage, lutGray, "Negative", format)
		// negative.NegativeButton(application, negativeImage, lutGray, format, "Negative")
	})
}
