package main

import (
	"vpc/pkg/menu"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"github.com/kbinani/screenshot"
)

func main() {
	application := app.New()
	mainWindow := application.NewWindow("Hello")
	window := screenshot.GetDisplayBounds(0)
	mainWindow.Resize(fyne.NewSize(float32(window.Bounds().Dx()),
		float32(window.Bounds().Dy())))
	openFileItem := menu.ButtonOpen(application, mainWindow)

	quitItem := fyne.NewMenuItem("Quit", func() {
		mainWindow.Close()
	})

	newItemSeparator := fyne.NewMenuItemSeparator()

	fileItem := fyne.NewMenu("File", openFileItem, newItemSeparator, quitItem)
	menu := fyne.NewMainMenu(fileItem)
	mainWindow.SetMainMenu(menu)
	mainWindow.ShowAndRun()
}
