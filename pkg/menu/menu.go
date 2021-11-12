package menu

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"strconv"
	calculate "vpc/pkg/calculate"
	histogrambutton "vpc/pkg/histogramButton"
	information "vpc/pkg/information"
	loadandsave "vpc/pkg/loadandsave"
	mouse "vpc/pkg/mouse"
	newwindow "vpc/pkg/newWindow"
	operations "vpc/pkg/operations"
	saveitem "vpc/pkg/saveItem"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func GeneralMenu(application fyne.App, grayImage *image.Gray,
	lutGray map[int]int, input, format string) {
	width := grayImage.Bounds().Dx()
	height := grayImage.Bounds().Dy()
	window := newwindow.NewWindow(application, width, height, input)
	image := canvas.NewImageFromImage(grayImage)
	text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
	canvasText := canvas.NewText(text, color.Opaque)
	window.SetContent(container.NewBorder(nil, canvasText, nil, nil, image,
		mouse.New(grayImage, canvasText, text)))

	_, values, numbersOfPixel, entropy, min, max, brightness, contrast :=
		calculate.Calculate(grayImage, width, height, format)

	informationText := information.Information(format, width, height, min, max,
		brightness, contrast, entropy)

	informationItem := fyne.NewMenuItem("Image Information", func() {
		dialog.ShowInformation("Information", informationText, window)
	})

	histogramItem := histogrambutton.HistogramButton(application, window, values,
		numbersOfPixel, "Histogram", false)

	cumulativeHistogramItem := histogrambutton.HistogramButton(application,
		window, values, numbersOfPixel, "Cumulative Histogram", true)

	negativeItem := negativeButton(application, grayImage, lutGray, width, height,
		format)

	gammaButton := gammaButton(application, grayImage, lutGray, input, format)

	brightnessAndContrastItem := brightnessAndContrastButton(application,
		grayImage, brightness, contrast, numbersOfPixel, lutGray, format)

	equalizationItem := fyne.NewMenuItem("Equalization", func() {
		GeneralMenu(application, operations.EqualizeAnImage(numbersOfPixel,
			grayImage), lutGray, input, format)
	})

	imageDifferenceItem := differenceButton(application, width, height, grayImage,
		lutGray, format, window)

	changeMapItem :=
		changeMapButton(application, width, height, grayImage, lutGray, format, window)

	sectionItem := sectionsButton(application)

	quitItem := quitButton(window)

	separatorItem := fyne.NewMenuItemSeparator()

	menuItem := fyne.NewMenu("File", saveitem.SaveItem(application, grayImage),
		separatorItem, quitItem)
	menuItem2 := fyne.NewMenu("Image Information", informationItem)
	menuItem3 := fyne.NewMenu("Operations", histogramItem, separatorItem,
		cumulativeHistogramItem, separatorItem, negativeItem, separatorItem,
		brightnessAndContrastItem, separatorItem, gammaButton, separatorItem,
		equalizationItem, separatorItem, imageDifferenceItem, separatorItem,
		changeMapItem, separatorItem, sectionItem)
	menu := fyne.NewMainMenu(menuItem, menuItem2, menuItem3)
	window.SetMainMenu(menu)
	window.Show()
	window.Close()
}

func gammaButton(application fyne.App, grayImage *image.Gray,
	lutGray map[int]int, input, format string) *fyne.MenuItem {
	return fyne.NewMenuItem("Gamma", func() {
		windowGamma := newwindow.NewWindow(application, 500, 250, "Gamma Value")
		input := widget.NewEntry()
		input.SetPlaceHolder("15.0")
		content := container.NewVBox(input, widget.NewButton("Enter", func() {
			number, err := strconv.ParseFloat(input.Text, 64)
			if err != nil {
				dialog.ShowError(err, windowGamma)
			} else if number > 20.0 || number < 0.05 {
				dialog.ShowError(errors.New("the value must be between 0.05 and 20.0"),
					windowGamma)
			} else {
				img := operations.Gamma(grayImage, grayImage.Bounds().Dx(),
					grayImage.Bounds().Dy(), number)
				GeneralMenu(application, img, lutGray, "Gamma Image", format)
				windowGamma.Close()
			}
		}))
		windowGamma.SetContent(content)
		windowGamma.Show()
	})
}

func brightnessAndContrastButton(application fyne.App, grayImage *image.Gray,
	brightness, contrast float64, numbersOfPixel map[int]int, lutGray map[int]int,
	format string) *fyne.MenuItem {
	return fyne.NewMenuItem("Brightness and Contrast", func() {
		newWindows := newwindow.NewWindow(application, 500, 500,
			"Brightness and Contrast")
		data := binding.NewFloat()
		data.Set(float64(int(brightness)))
		slide := widget.NewSliderWithData(0, 255, data)
		formatted := binding.FloatToStringWithFormat(data, "Float value: %0.2f")
		label := widget.NewLabelWithData(formatted)

		data2 := binding.NewFloat()
		data2.Set(float64(int(contrast)))
		slide2 := widget.NewSliderWithData(0, 127, data2)
		formatted2 := binding.FloatToStringWithFormat(data2, "Float value: %0.2f")
		label2 := widget.NewLabelWithData(formatted2)

		content := widget.NewButton("Calculate", func() {
			width := grayImage.Bounds().Dx()
			height := grayImage.Bounds().Dy()
			bright, _ := data.Get()
			conts, _ := data2.Get()
			newImage := operations.AdjustBrightnessAndContrast(bright, conts,
				numbersOfPixel, grayImage, width*height)
			GeneralMenu(application, newImage, lutGray, "Modified Image", format)
		})

		brightnessText := canvas.NewText("Brightness", color.White)
		brightnessText.TextStyle = fyne.TextStyle{Bold: true}
		contrastText := canvas.NewText("Contrast", color.White)
		contrastText.TextStyle = fyne.TextStyle{Bold: true}
		menuAndImageContainer := container.NewVBox(brightnessText, label, slide,
			contrastText, label2, slide2, content)

		newWindows.SetContent(menuAndImageContainer)
		newWindows.Show()
	})
}

func negativeButton(application fyne.App, grayImage *image.Gray,
	lutGray map[int]int, width, height int, format string) *fyne.MenuItem {
	return fyne.NewMenuItem("Negative", func() {
		negativeImage := operations.Negative(grayImage, lutGray, width, height)
		GeneralMenu(application, negativeImage, lutGray, "Negative", format)
	})
}

func differenceButton(application fyne.App, width, height int,
	grayImage *image.Gray, lutGray map[int]int, format string,
	window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem("Image difference", func() {
		newDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				fileName := reader.URI().String()[7:]
				image, _, err := loadandsave.LoadImage(fileName)
				if err != nil {
					dialog.ShowError(err, window)
				} else if grayImage.Bounds().Dx() != image.Bounds().Dx() ||
					grayImage.Bounds().Dy() != image.Bounds().Dy() {
					dialog.ShowError(errors.New("the image must have the same dimensions"), window)
				} else {
					newWindow := newwindow.NewWindow(application, image.Bounds().Dx(),
						image.Bounds().Dy(), "Used Image")
					canvasImage := canvas.NewImageFromImage(image)
					newWindow.SetContent(canvasImage)
					newWindow.Show()
					difference, err := operations.ImageDifference(grayImage, image)
					if err != nil {
						dialog.ShowError(err, window)
					} else {
						GeneralMenu(application, difference, lutGray, "Difference", format)
					}
				}
			}
		}, window)
		newDialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".png",
			".jpeg", ".tiff"}))
		newDialog.Show()
	})
}

func changeMapButton(application fyne.App, width, height int,
	grayImage *image.Gray, lutGray map[int]int, format string,
	window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem("Change Map", func() {
		newDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				fileName := reader.URI().String()[7:]
				image, _, err := loadandsave.LoadImage(fileName)
				if err != nil {
					dialog.ShowError(err, window)
				} else if grayImage.Bounds().Dx() != image.Bounds().Dx() ||
					grayImage.Bounds().Dy() != image.Bounds().Dy() {
					dialog.ShowError(errors.New("the image must have the same dimensions"), window)
				} else {
					newWindow := newwindow.NewWindow(application, image.Bounds().Dx(),
						image.Bounds().Dy(), "Used Image")
					canvasImageUsed := canvas.NewImageFromImage(image)
					newWindow.SetContent(canvasImageUsed)
					newWindow.Show()

					difference, _ := operations.ImageDifference(grayImage, image)
					differenceWindow := newwindow.NewWindow(application,
						difference.Bounds().Dx(), difference.Bounds().Dy(), "Difference")
					canvasImageDifference := canvas.NewImageFromImage(difference)
					differenceWindow.SetContent(canvasImageDifference)

					_, values, numbersOfPixel, _, _, _, _, _ :=
						calculate.Calculate(difference, width, height, format)

					histogramItem := histogrambutton.HistogramButton(application, window,
						values, numbersOfPixel, "Histogram", false)

					cumulativeHistogramItem := histogrambutton.HistogramButton(application,
						window, values, numbersOfPixel, "Cumulative Histogram", true)

					thresHoldItem := fyne.NewMenuItem("Threshold", func() {
						windowThreshold := newwindow.NewWindow(application, 500, 200, "Threshold Value")
						data := binding.NewFloat()
						data.Set(0)
						slide := widget.NewSliderWithData(0, 255, data)
						formatted := binding.FloatToStringWithFormat(data, "Float value: %0.2f")
						label := widget.NewLabelWithData(formatted)

						content := widget.NewButton("Calculate", func() {
							threshold, _ := data.Get()
							newImage := operations.ChangeMap(difference, image, threshold)
							windowResult := newwindow.NewWindow(application,
								newImage.Bounds().Dx(), newImage.Bounds().Dy(), "Result")
							canvasR := canvas.NewImageFromImage(newImage)

							quitItem := quitButton(windowResult)
							separatorItem := fyne.NewMenuItemSeparator()
							menuItem := fyne.NewMenu("File", saveitem.SaveItem(application,
								newImage), separatorItem, quitItem)
							menu := fyne.NewMainMenu(menuItem)
							windowResult.SetMainMenu(menu)

							windowResult.SetContent(canvasR)
							windowResult.Show()
						})

						threshold := canvas.NewText("Threshold", color.White)
						threshold.TextStyle = fyne.TextStyle{Bold: true}
						menuAndImageContainer := container.NewVBox(threshold, label, slide,
							content)

						windowThreshold.SetContent(menuAndImageContainer)
						windowThreshold.Show()
					})

					quitItem := quitButton(differenceWindow)

					separatorItem := fyne.NewMenuItemSeparator()

					menuItem := fyne.NewMenu("File", saveitem.SaveItem(application,
						difference), separatorItem, quitItem)
					menuItem2 := fyne.NewMenu("User value", thresHoldItem)
					menuItem3 := fyne.NewMenu("Histograms", histogramItem, separatorItem,
						cumulativeHistogramItem)
					menu := fyne.NewMainMenu(menuItem, menuItem2, menuItem3)
					differenceWindow.SetMainMenu(menu)
					differenceWindow.Show()
				}
			}
		}, window)
		newDialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".png",
			".jpeg", ".tiff"}))
		newDialog.Show()
	})
}

func sectionsButton(application fyne.App) *fyne.MenuItem {
	return fyne.NewMenuItem("Linear Adjustment in Sections", func() {
		windowSections := newwindow.NewWindow(application, 500, 200, "Sections Number")
		input := widget.NewEntry()
		input.SetPlaceHolder("0")
		content := container.NewVBox(input, widget.NewButton("Enter", func() {
			number, err := strconv.Atoi(input.Text)
			if err != nil {
				dialog.ShowError(err, windowSections)
			} else {
				windowValues := newwindow.NewWindow(application, 500, 500, "Sections Values")

				label1 := widget.NewLabel("Coordinates X: ")
				label2 := widget.NewLabel("Coordinates Y: ")
				coordinatesX := container.NewVBox(label1)
				coordinatesY := container.NewVBox(label2)
				var point *widget.Entry
				var point2 *widget.Entry
				for i := 0; i < number; i++ {
					point = widget.NewEntry()
					point.SetPlaceHolder("x:")
					point2 = widget.NewEntry()
					point2.SetPlaceHolder("y:")
					coordinatesX.Add(point)
					coordinatesY.Add(point2)
				}
				content := container.NewVBox(container.NewHBox(coordinatesX, coordinatesY),
					widget.NewButton("Enter", func() {
						var coordinates []pair
						for i := 0; i < number; i++ {
							point1, _ := strconv.Atoi(point.Text)
							point2, _ := strconv.Atoi(point2.Text)
							coordinates = append(coordinates, pair{x: point1, y: point2})
						}
						for i := 0; i < len(coordinates); i++ {
							fmt.Println(coordinates[i])
						}
					}))
				windowValues.SetContent(content)
				windowValues.Show()

				/*img := operations.Gamma(grayImage, grayImage.Bounds().Dx(),
					grayImage.Bounds().Dy(), number)
				GeneralMenu(application, img, lutGray, "Gamma Image", format)
				windowGamma.Close()*/
			}
		}))
		windowSections.SetContent(content)
		windowSections.Show()
	})
}

func quitButton(window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem("Quit", func() {
		window.Close()
	})
}

type pair struct {
	x, y int
}
