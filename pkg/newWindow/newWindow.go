package newwindow

import "fyne.io/fyne/v2"

func NewWindow(application fyne.App, width, height int, name string) fyne.Window {
	window := application.NewWindow(name)
	window.Resize(fyne.NewSize(float32(width), float32(height)))
	return window
}
