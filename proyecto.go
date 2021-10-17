package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/canvas"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	"github.com/kbinani/screenshot"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
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
	return window
}

func buttonOpen(application fyne.App) *fyne.MenuItem {
	fileItem := fyne.NewMenuItem("Open image", func() {
		fileWindow := application.NewWindow("OpenFile")
		fileWindow.Resize(fyne.NewSize(500, 500))
		// window := screenshot.GetDisplayBounds(0)
		// fileWindow.Resize(fyne.NewSize(window.Bounds().Dx()/2, window.Bounds().Dy()/2))

		input := widget.NewEntry()
		input.SetPlaceHolder("ejemplo.png")

		content := container.NewVBox(input, widget.NewButton("Open", func() {
			img, err := loadImage(input.Text)
			check(err)
			width := img.Bounds().Dx()
			height := img.Bounds().Dy()
			fmt.Printf("Width: %d\n", width)
			fmt.Printf("Height: %d\n", height)

			imageWindow := newWindow(application, img.Bounds().Dx(), img.Bounds().Dy(), input.Text)
			image := canvas.NewImageFromFile(input.Text)
			imageWindow.SetContent(image)

			newItem := fyne.NewMenuItem("Histogram", func() { fmt.Println("Falta por hacer") })
			newItem2 := fyne.NewMenuItem("Cumulative histogram", func() { fmt.Println("Falta por hacer") })
			newItem3 := fyne.NewMenuItem("Scale gray", func() {
				img := scaleGray(img, width, height)
				GrayButton(application, img, input.Text)
			})

			menuItem := fyne.NewMenu("Operations", newItem, newItem2, newItem3)
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

func GrayButton(application fyne.App, img *image.Gray, input string) {
	width := img.Bounds().Dx()
	height := img.Bounds().Dy()
	window := newWindow(application, width, height, input)
	image := canvas.NewImageFromImage(img)
	window.SetContent(image)
	colors, values := colorsValues(img)

	newItem := fyne.NewMenuItem("Histogram", func() {
		plote(histogram(colors), values)
	})
	newItem2 := fyne.NewMenuItem("Cumulative histogram", func() {
		plote(lutGray(), cumulativeHistogram(values))
	})
	newItem3 := fyne.NewMenuItem("Negative", func() {
		img := negative(img, lutGray(), width, height)
		GrayButton(application, img, input)
	})

	menuItem := fyne.NewMenu("Operations", newItem, newItem2, newItem3)
	menu := fyne.NewMainMenu(menuItem)
	window.SetMainMenu(menu)
	window.Show()
	window.Close()
}

func loadImage(fileName string) (image.Image, error) {
	fmt.Println("Load the image:", fileName)

	fimg, err := os.Open(fileName)
	check(err)

	// fmt.Println("Dirección de memoria de la imagen: ", fimg)
	defer fimg.Close()

	img, _, err := image.Decode(fimg)
	check(err)

	return img, err
}

func saveImage(fileName string, img image.Image) error {
	fmt.Println("Saving the image:", fileName)

	var err error
	var fimg *os.File
	extension := fileName[len(fileName)-3:]
	if extension == "jpg" || extension == "jpeg" || extension == "png" {
		fimg, err = os.Create(fileName)
		check(err)
	}
	defer fimg.Close()

	err = checkImgFormat(extension, fimg, img)
	return err
}

func scaleGray(img image.Image, width, height int) *image.Gray {
	img2 := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			r, g, b, _ := img.At(i, j).RGBA()
			r, g, b = r>>8, g>>8, b>>8
			//y := 0.375*float64(r) + 0.5*float64(g) + 0.125*float64(b)
			y := 0.299*float64(r) + 0.587*float64(g) + 0.114*float64(b) // our
			grayColor := color.Gray{uint8(y)}
			img2.Set(i, j, grayColor)
		}
	}
	return img2
}

func colorsValues(image *image.Gray) ([]uint64, plotter.Values) {
	var colors []uint64
	var values plotter.Values

	for i := 0; i < image.Bounds().Dx(); i++ {
		for j := 0; j < image.Bounds().Dy(); j++ {
			y := image.GrayAt(i, j).Y
			colors = append(colors, uint64(y))
			values = append(values, float64(y))
		}
	}
	return colors, values
}

func negative(img *image.Gray, lutGray map[int]int, width, height int) *image.Gray {
	img2 := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{width, height}})
	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			newColor := color.Gray{uint8(float32(lutGray[int(img.GrayAt(i, j).Y)]))}
			img2.Set(i, j, newColor)
		}
	}
	return img2
}

func lutGray() map[int]int {
	table := make(map[int]int)
	for i := 0; i <= 255; i++ {
		table[i] = 255 - i
	}
	return table
}

func histogram(colors []uint64) map[int]int {
	histogram := make(map[int]int)
	for i := 0; i <= 255; i++ {
		cont := 0
		for j := 0; j < len(colors); j++ {
			if i == int(colors[j]) {
				cont++
			}
		}
		histogram[i] = cont
	}
	// fmt.Println(histogram)
	return histogram
}

func cumulativeHistogram(values plotter.Values) plotter.Values {
	var newValues plotter.Values
	var cont float64
	for i := 0; i <= len(values); i++ {
		cont += values[i]
		values = append(values, cont)
	}
	// fmt.Println(histogram)
	return newValues
}

func valueRange(histogram map[int]int) (int, int) {
	// 0 Negro
	// 255 Blanco
	min := 300 // Negro
	max := 0   // Blanco
	for i := 0; i < len(histogram); i++ {
		if i >= max && histogram[i] != 0 {
			max = i
		}
		if i <= min && histogram[i] != 0 {
			min = i
		}
	}
	return min, max
}

func plote(histogram map[int]int, values plotter.Values) {

	p := plot.New()

	p.Title.Text = "Histogram plot"

	hist, err2 := plotter.NewHist(values, len(histogram))
	if err2 != nil {
		panic(err2)
	}
	// hist.Normalize(1)
	p.Add(hist)

	if err := p.Save(3*vg.Inch, 3*vg.Inch, "hist.png"); err != nil {
		panic(err)
	}
}

func checkImgFormat(extension string, fimg *os.File, img image.Image) error {
	var err error
	switch extension {
	case "jpg", "jpeg":
		err = jpeg.Encode(fimg, img, nil)
	case "png":
		err = png.Encode(fimg, img)
	default:
		err = errors.New("unsupported format")
	}
	return err
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
