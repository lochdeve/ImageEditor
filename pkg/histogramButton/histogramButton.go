package histogrambutton

import (
	"image/color"
	"strconv"
	"vpc/pkg/histogram"
	"vpc/pkg/loadandsave"

	newwindow "vpc/pkg/newWindow"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"gonum.org/v1/plot/plotter"
)

func HistogramButton(application fyne.App, window fyne.Window,
	values plotter.Values, numbersOfPixel map[int]int, name string,
	cumulative bool) *fyne.MenuItem {
	return fyne.NewMenuItem(name, func() {
		histogram.Plote(numbersOfPixel, values, cumulative)
		histogramImage, _, err := loadandsave.LoadImage(".tmp/hist.png")
		if err != nil {
			dialog.ShowError(err, window)
		} else {
			width := histogramImage.Bounds().Dx()
			height := histogramImage.Bounds().Dy()
			windowImage := newwindow.NewWindow(application, width, height, "Histogram")
			text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
			canvasText := canvas.NewText(text, color.Opaque)
			image := canvas.NewImageFromImage(histogramImage)
			windowImage.SetContent(container.NewBorder(nil, canvasText, nil, nil, image))
			windowImage.Show()
		}
	})
}
