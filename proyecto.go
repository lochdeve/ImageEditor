package main

import (
	"image"
	"image/color"
	"strconv"
	"strings"

	histogram "vpc/pkg/histogram"
	loadandsave "vpc/pkg/loadandsave"
	operations "vpc/pkg/operations"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
	"gonum.org/v1/plot/plotter"
)

type MyWidget struct {
	widget.Icon
	image image.Image
	text  *canvas.Text
}

func new(image1 image.Image, text1 *canvas.Text) *MyWidget {
	m := &MyWidget{image: image1, text: text1}
	m.ExtendBaseWidget(m)
	return m
}

func (t *MyWidget) Tapped(_ *fyne.PointEvent) {
	//fmt.Println("I have been tapped")
}

func (w *MyWidget) FocusGained() {
	//fmt.Println("FocusGained")
}

// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (w *MyWidget) FocusLost() {
	//fmt.Println("Lost focus")
}

// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (w *MyWidget) TypedRune(_ rune) {

}

// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (w *MyWidget) TypedKey(_ *fyne.KeyEvent) {

}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (w *MyWidget) MouseIn(*desktop.MouseEvent) {
	//	fmt.Println("Inside")
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (w *MyWidget) MouseMoved(hola *desktop.MouseEvent) {
	//fmt.Println("It moves")
	r, g, b, _ := w.image.At(int(hola.AbsolutePosition.X), int(hola.AbsolutePosition.Y)).RGBA()
	r, g, b = r>>8, g>>8, b>>8
	textaux := "R:" + strconv.Itoa(int(r)) + " G:" + strconv.Itoa(int(g)) + " B:" + strconv.Itoa(int(b))
	w.text.Text = textaux
	w.text.Refresh()
	//fmt.Println(textaux)
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (w *MyWidget) MouseOut() {
	//fmt.Println("get out")
}

func (mouse *MyWidget) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewImageFromImage(mouse.image))
}

func main() {
	interfaz()
}

func interfaz() {
	application := app.New()
	mainWindow := application.NewWindow("Hello")
	// mainWindow.Resize(fyne.NewSize(500, 500))
	window := screenshot.GetDisplayBounds(0)
	mainWindow.Resize(fyne.NewSize(float32(window.Bounds().Dx()), float32(window.Bounds().Dy())))
	fileItem := buttonOpen(application, mainWindow)

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
					grayImage := scaleGray(colorImage, width, height)
					_, _, _, min, max, brightness, contrast := calculate(grayImage, width, height, format)
					informationTape := information(format, width, height, min, max, brightness, contrast)
					lutGray := operations.LutGray()

					windowName := strings.Split(fileName, "/")
					imageWindow := newWindow(application, colorImage.Bounds().Dx(), colorImage.Bounds().Dy(), windowName[len(windowName)-1])
					canvasImage := canvas.NewImageFromImage(colorImage)
					text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
					canvasText := canvas.NewText(text, color.Opaque)
					imageWindow.SetContent(container.NewBorder(nil, canvasText, nil, nil, canvasImage, new(colorImage, canvasText)))
					//imageWindow.SetContent(new())
					newItem := fyne.NewMenuItem("Image Information", func() {
						dialog.ShowInformation("Information", informationTape, imageWindow)
					})

					newItem2 := fyne.NewMenuItem("Scale gray", func() {
						GrayButton(application, grayImage, lutGray, windowName[len(windowName)-1], format, informationTape)
					})

					newItem3 := fyne.NewMenuItem("Quit", func() {
						imageWindow.Close()
					})

					newItemSeparator := fyne.NewMenuItemSeparator()

					menuItem := fyne.NewMenu("Operations", newItem, newItemSeparator, newItem2, newItemSeparator, newItem3)
					menu := fyne.NewMainMenu(menuItem)
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
		numbersofpixel, values := operations.ColorsValues(grayImage)
		histogram.Plote(histogram.NumbersOfPixel(numbersofpixel), values)
		histogramImage, _, _ := loadandsave.LoadImage("hist.png")
		width := histogramImage.Bounds().Dx()
		height := histogramImage.Bounds().Dy()
		windowImage := newWindow(application, width, height, "histogram")
		text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
		canvasText := canvas.NewText(text, color.Opaque)
		image := canvas.NewImageFromImage(histogramImage)
		windowImage.SetContent(container.NewBorder(nil, canvasText, nil, nil, image))
		windowImage.Show()
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
