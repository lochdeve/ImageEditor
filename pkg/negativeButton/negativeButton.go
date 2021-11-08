package negative

import (
	"image"
	"image/color"
	"strconv"
	"vpc/pkg/mouse"
	"vpc/pkg/operations"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"

	calculate "vpc/pkg/calculate"
	histogrambutton "vpc/pkg/histogramButton"
	information "vpc/pkg/information"
	newwindow "vpc/pkg/newWindow"
	saveitem "vpc/pkg/saveItem"
)

func NegativeButton(application fyne.App, negativeImage *image.Gray,
	lutGray map[int]int, format, input string) {
	width := negativeImage.Bounds().Dx()
	height := negativeImage.Bounds().Dy()
	window := newwindow.NewWindow(application, width, height, input)
	image := canvas.NewImageFromImage(negativeImage)
	text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
	canvasText := canvas.NewText(text, color.Opaque)
	window.SetContent(container.NewBorder(nil, canvasText, nil, nil, image, mouse.New(negativeImage, canvasText, text)))

	_, values, numbersOfPixel, entropy, min, max, brightness, contrast :=
		calculate.Calculate(negativeImage, width, height, format)

	informationTape := information.Information(format, width, height, min, max, brightness, contrast, entropy)

	informationItem := fyne.NewMenuItem("Image Information", func() {
		dialog.ShowInformation("Information", informationTape, window)
	})

	histogramItem := histogrambutton.HistogramButton(application, window, values, numbersOfPixel, "Histogram", false)

	cumulativeHistogramItem := histogrambutton.HistogramButton(application, window, values, numbersOfPixel, "Cumulative Histogram", true)

	negativeImageItem := fyne.NewMenuItem("Negative", func() {
		negativeImage := operations.Negative(negativeImage, lutGray, width, height)
		NegativeButton(application, negativeImage, lutGray, format, "Negative")
	})

	newItem5 := fyne.NewMenuItem("Quit", func() {
		window.Close()
	})

	newItemSeparator := fyne.NewMenuItemSeparator()

	menuItem := fyne.NewMenu("File", saveitem.SaveItem(application, negativeImage), newItemSeparator, newItem5)
	menuItem2 := fyne.NewMenu("Image Information", informationItem)
	menuItem3 := fyne.NewMenu("Operations", histogramItem, newItemSeparator, cumulativeHistogramItem, newItemSeparator, negativeImageItem)
	menu := fyne.NewMainMenu(menuItem, menuItem2, menuItem3)
	window.SetMainMenu(menu)
	window.Show()
	window.Close()
}
