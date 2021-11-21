package main

import (
	"vpc/pkg/menu"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/kbinani/screenshot"
)

func main() {
	application := app.New()
	mainWindow := application.NewWindow("Hello")
	window := screenshot.GetDisplayBounds(0)
	mainWindow.Resize(fyne.NewSize(float32(window.Bounds().Dx()),
		float32(window.Bounds().Dy())))

	canvasImage := canvas.NewImageFromFile(".tmp/ULL.jpg")
	canvasImage.FillMode = canvas.ImageFillOriginal
	text := widget.NewLabel("This prototype is for Computer vision a subject\nof 4º grade of Computer Engineering. \n\nWas develop by:\n\t- Carlos García Lezcano\n\t- Eduardo Expósito Barrera")

	openFileItem := menu.ButtonOpen(application, mainWindow)

	quitItem := fyne.NewMenuItem("Quit", func() {
		mainWindow.Close()
	})

	newItemSeparator := fyne.NewMenuItemSeparator()

	fileItem := fyne.NewMenu("File", openFileItem, newItemSeparator, quitItem)
	menu := fyne.NewMainMenu(fileItem)
	mainWindow.SetMainMenu(menu)
	mainWindow.SetContent(container.NewVBox(container.NewCenter(text),
		container.NewCenter(canvasImage)))
	mainWindow.ShowAndRun()
}