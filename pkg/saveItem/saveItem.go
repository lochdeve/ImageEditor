package saveitem

import (
	"image"
	"vpc/pkg/loadandsave"
	newwindow "vpc/pkg/newWindow"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func SaveItem(application fyne.App, image image.Image) *fyne.MenuItem {
	return fyne.NewMenuItem("Save Image", func() {
		fileWindow := newwindow.NewWindow(application, 500, 500, "SaveFile")
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
}
