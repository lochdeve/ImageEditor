package graybutton

import (
	"image"
	"image/color"
	"strconv"
	"strings"
	"vpc/pkg/calculate"
	histogrambutton "vpc/pkg/histogramButton"
	"vpc/pkg/information"
	"vpc/pkg/loadandsave"
	"vpc/pkg/mouse"
	negative "vpc/pkg/negativeButton"
	newwindow "vpc/pkg/newWindow"
	"vpc/pkg/operations"
	saveitem "vpc/pkg/saveItem"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func GrayButton(application fyne.App, grayImage *image.Gray, lutGray map[int]int,
	input, format, informationText string) {
	width := grayImage.Bounds().Dx()
	height := grayImage.Bounds().Dy()
	window := newwindow.NewWindow(application, width, height, input)
	image := canvas.NewImageFromImage(grayImage)
	text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
	canvasText := canvas.NewText(text, color.Opaque)
	window.SetContent(container.NewBorder(nil, canvasText, nil, nil, image,
		mouse.New(grayImage, canvasText, text)))

	_, values, numbersOfPixel, _, _, _, brightness, contrast :=
		calculate.Calculate(grayImage, width, height, format)

	imageInformationItem := fyne.NewMenuItem("Image Information", func() {
		dialog.ShowInformation("Information", informationText, window)
	})

	histogramItem := histogrambutton.HistogramButton(application, window, values, numbersOfPixel, "Histogram", false)

	cumulativeHistogramItem := histogrambutton.HistogramButton(application, window, values, numbersOfPixel, "Cumulative Histogram", true)

	negativeItem := fyne.NewMenuItem("Negative", func() {
		negativeImage := operations.Negative(grayImage, lutGray, width, height)
		negative.NegativeButton(application, negativeImage, lutGray, format, "Negative")
	})

	quitItem := fyne.NewMenuItem("Quit", func() {
		window.Close()
	})

	brightnessAndContrastItem := fyne.NewMenuItem("Brightness and Contrast", func() {
		newWindows := newwindow.NewWindow(application, 500, 500, "Brightness and Contrast")
		data := binding.NewFloat()
		data.Set(float64(int(brightness)))
		slide := widget.NewSliderWithData(0, 255, data)
		formatted := binding.FloatToStringWithFormat(data, "Float value: %0.2f")
		label := widget.NewLabelWithData(formatted)

		data2 := binding.NewFloat()
		data2.Set(float64(int(contrast)))
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
			_, _, _, entropy, newMin, newMax, newBrightness, newContrast :=
				calculate.Calculate(newImage, width, height, format)
			newInformationText := information.Information(format, width, height, newMin, newMax,
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

	imageDifferenceItem := differenceDialogItem(application, window, grayImage)

	/*changeMap := fyne.NewMenuItem("Change Map", func() {
		differenceDialogItem(application, width, height, grayImage)
	})*/

	separatorItem := fyne.NewMenuItemSeparator()

	menuItem := fyne.NewMenu("File", saveitem.SaveItem(application, grayImage), separatorItem, quitItem)
	menuItem2 := fyne.NewMenu("Image Information", imageInformationItem)
	menuItem3 := fyne.NewMenu("Operations", histogramItem, separatorItem,
		cumulativeHistogramItem, separatorItem, negativeItem, separatorItem,
		brightnessAndContrastItem, separatorItem, separatorItem, imageDifferenceItem) // , separatorItem, changeMap)
	menu := fyne.NewMainMenu(menuItem, menuItem2, menuItem3)
	window.SetMainMenu(menu)
	window.Show()
	window.Close()
}

func differenceDialogItem(application fyne.App, window fyne.Window,
	grayImage *image.Gray) *fyne.MenuItem {
	dialogItem := fyne.NewMenuItem("Image difference", func() {
		// windowImage := newwindow.NewWindow(application, width, height, "Difference")
		hola(application, window, grayImage)
		/*item := buttonOpen(application, windowImage)
		menuItem := fyne.NewMenu("Operations", item)
		menu := fyne.NewMainMenu(menuItem)
		windowImage.SetMainMenu(menu)*/
		// windowImage.Show()
	})
	return dialogItem
}

func hola(application fyne.App, window fyne.Window, grayImage *image.Gray) {
	dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
		if reader != nil {
			fileName := reader.URI().String()[7:]
			image, format, err := loadandsave.LoadImage(fileName)
			if err != nil {
				dialog.ShowError(err, window)
			}
			difference, err := operations.ImageDifference(grayImage, image)
			if err != nil {
				dialog.ShowError(err, window)
			}
			windowName := strings.Split(fileName, "/")
			newWindow := newwindow.NewWindow(application, image.Bounds().Dx(),
				image.Bounds().Dy(), "Used image for difference: "+windowName[len(windowName)-1])
			canvasImage1 := canvas.NewImageFromImage(image)
			newWindow.SetContent(canvasImage1)
			newWindow.Show()
			width := difference.Bounds().Dx()
			height := difference.Bounds().Dy()
			_, _, _, entropy, min, max, brightness, contrast :=
				calculate.Calculate(difference, width, height, format)
			informationDifference :=
				information.Information(format, width, height, min, max, brightness, contrast, entropy)
			GrayButton(application, difference, operations.LutGray(), "Difference", format, informationDifference)
			/*windowDifference := newwindow.NewWindow(application, image.Bounds().Dx(),
				image.Bounds().Dy(), "Difference")
			canvasImage := canvas.NewImageFromImage(difference)
			windowDifference.SetContent(canvasImage)
			windowDifference.Show()*/
			/*if opcion == 1 {
				canvasImage := canvas.NewImageFromImage(image)
				newWindow := newWindow(application, image.Bounds().Dx(), image.Bounds().Dy(), fileName)
				newWindow.SetContent(canvasImage)
				newWindow.Show()
				GrayButton(application, difference, nil, "", "", "")
			} else {
			}*/
		}
	}, window)
	dialog.Show()
}
