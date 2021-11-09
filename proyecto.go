package main

import (
	"image/color"
	"strconv"
	"strings"

	calculate "vpc/pkg/calculate"
	information "vpc/pkg/information"
	loadandsave "vpc/pkg/loadandsave"
	"vpc/pkg/menu"
	mouse "vpc/pkg/mouse"
	newwindow "vpc/pkg/newWindow"
	operations "vpc/pkg/operations"
	saveitem "vpc/pkg/saveItem"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
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
					_, _, _, entropy, min, max, brightness, contrast :=
						calculate.Calculate(grayImage, width, height, format)
					informationTape := information.Information(format, width, height, min,
						max, brightness, contrast, entropy)

					lutGray := operations.LutGray()

					windowName := strings.Split(fileName, "/")
					imageWindow := newwindow.NewWindow(application, colorImage.Bounds().Dx(), colorImage.Bounds().Dy(), windowName[len(windowName)-1])
					canvasImage := canvas.NewImageFromImage(colorImage)
					text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
					canvasText := canvas.NewText(text, color.Opaque)
					imageWindow.SetContent(container.NewBorder(nil, canvasText, nil, nil, canvasImage, mouse.New(colorImage, canvasText, text)))

					imageInformationItem := fyne.NewMenuItem("Image Information", func() {
						dialog.ShowInformation("Information", informationTape, imageWindow)
					})

					scaleGrayItem := fyne.NewMenuItem("Scale gray", func() {
						menu.GeneralMenu(application, grayImage, lutGray,
							windowName[len(windowName)-1], format)
					})

					quitItem := fyne.NewMenuItem("Quit", func() {
						imageWindow.Close()
					})

					separatorItem := fyne.NewMenuItemSeparator()
					saveItem := fyne.NewMenu("File", saveitem.SaveItem(application, grayImage), separatorItem, quitItem)

					imageInformationMenu := fyne.NewMenu("ImageInformation", imageInformationItem)
					operationItem := fyne.NewMenu("Operations", scaleGrayItem)
					menu := fyne.NewMainMenu(saveItem, imageInformationMenu, operationItem)
					imageWindow.SetMainMenu(menu)
					imageWindow.Show()
				}
			}
		}, window)
		fd.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".png", ".jpeg", ".tiff"}))
		fd.Show()
	})
	return fileItem
}

/*func differenceDialogItem(application fyne.App, width, height, opcion int,
	grayImage *image.Gray) *fyne.MenuItem {
	dialogItem := fyne.NewMenuItem("Image difference", func() {
		windowImage := newwindow.NewWindow(application, width, height, "difference")
		hola(application, windowImage, grayImage, opcion)
		item := buttonOpen(application, windowImage)
		menuItem := fyne.NewMenu("Operations", item)
		menu := fyne.NewMainMenu(menuItem)
		windowImage.SetMainMenu(menu)
		windowImage.Show()
	})
	return dialogItem
}*/

/*func hola(application fyne.App, windowImage fyne.Window, grayImage *image.Gray, opcion int) {
var Edu string
dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
	if reader != nil {
		fileName := reader.URI().String()[7:]
		image, _, err := loadandsave.LoadImage(fileName)
		if err != nil {
			dialog.ShowError(err, windowImage)
		}
		difference, err := operations.ImageDifference(grayImage, image)
		Edu = "Edu"
		if err != nil {
			dialog.ShowError(err, windowImage)
		}
		if opcion == 1 {
			canvasImage := canvas.NewImageFromImage(image)
			newWindow := newwindow.NewWindow(application, image.Bounds().Dx(), image.Bounds().Dy(), fileName)
			newWindow.SetContent(canvasImage)
			newWindow.Show()
			graybutton.GrayButton(application, difference, nil, "", "", "")
		} else {
			canvasImage := canvas.NewImageFromImage(difference)
			windowImage.SetContent(canvasImage)
		}
	}
}, windowImage)
fmt.Print(Edu)
// dialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".png", ".jpeg"}))

/*dialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
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

		if opcion == 1 {
			canvasImage := canvas.NewImageFromImage(image)
			newWindow := newWindow(application, image.Bounds().Dx(), image.Bounds().Dy(), fileName)
			newWindow.SetContent(canvasImage)
			newWindow.Show()
			GrayButton(application, difference, nil, "", "", "")
		} else {
			canvasImage := canvas.NewImageFromImage(difference)
			windowImage.SetContent(canvasImage)
		}
	}
}, windowImage)
dialog.Show()*/
// }
