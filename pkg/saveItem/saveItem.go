package saveitem

import (
	"image"
	"vpc/pkg/loadandsave"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func SaveItem(application fyne.App, image image.Image) *fyne.MenuItem {
	saveImageItem := fyne.NewMenuItem("Save Image", func() {
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
	return saveImageItem
}
