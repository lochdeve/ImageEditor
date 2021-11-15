package menu

import (
	"errors"
	"image"
	"image/color"
	"strconv"
	"strings"
	histogram "vpc/pkg/histogram"
	imagecontent "vpc/pkg/imageContent"
	"vpc/pkg/information"
	loadandsave "vpc/pkg/loadandsave"
	mouse "vpc/pkg/mouse"
	newwindow "vpc/pkg/newWindow"
	operations "vpc/pkg/operations"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func generalMenu(application fyne.App, fullImage imagecontent.InformationImage,
	input string) {
	width := fullImage.Image().Bounds().Dx()
	height := fullImage.Image().Bounds().Dy()
	window := newwindow.NewWindow(application, width, height, input)
	image := canvas.NewImageFromImage(fullImage.Image())
	text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
	canvasText := canvas.NewText(text, color.Opaque)
	window.SetContent(container.NewBorder(nil, canvasText, nil, nil, image,
		mouse.New(fullImage.Image(), canvasText, text)))

	informationText := information.Information(fullImage)

	informationItem := fyne.NewMenuItem("Image Information", func() {
		dialog.ShowInformation("Information", informationText, window)
	})

	histogramItem := histogramButton(application, window, fullImage,
		"Histogram", false)

	cumulativeHistogramItem := histogramButton(application, window, fullImage,
		"Cumulative Histogram", true)

	negativeItem := negativeButton(application, fullImage)

	gammaButton := gammaButton(application, fullImage, input)

	brightnessAndContrastItem := brightnessAndContrastButton(application, fullImage)

	equalizationItem := equalizationButton(application, fullImage, input)

	imageDifferenceItem := differenceButton(application, window, fullImage)

	changeMapItem := changeMapButton(application, window, fullImage)

	sectionItem := sectionsButton(application, fullImage)

	histogramSpecificationItem := histogramSpecificationButton(application, window, fullImage)

	roiItem := roiButton(application, fullImage)

	quitItem := quitButton(window)

	separatorItem := fyne.NewMenuItemSeparator()

	menuItem := fyne.NewMenu("File", saveButton(application, fullImage.Image()),
		separatorItem, quitItem)
	menuItem2 := fyne.NewMenu("Image Information", informationItem)
	menuItem3 := fyne.NewMenu("Operations", histogramItem, separatorItem,
		cumulativeHistogramItem, separatorItem, negativeItem, separatorItem,
		brightnessAndContrastItem, separatorItem, gammaButton, separatorItem,
		equalizationItem, separatorItem, imageDifferenceItem, separatorItem,
		changeMapItem, separatorItem, sectionItem, separatorItem,
		histogramSpecificationItem, separatorItem, roiItem)
	menu := fyne.NewMainMenu(menuItem, menuItem2, menuItem3)
	window.SetMainMenu(menu)
	window.Show()
	window.Close()
}

func gammaButton(application fyne.App, fullImage imagecontent.InformationImage,
	input string) *fyne.MenuItem {
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
				newFullImage := imagecontent.New(operations.Gamma(fullImage, number),
					fullImage.LutGray(), fullImage.Format())
				generalMenu(application, newFullImage, "Gamma Image")
				windowGamma.Close()
			}
		}))
		windowGamma.SetContent(content)
		windowGamma.Show()
	})
}

func brightnessAndContrastButton(application fyne.App,
	fullImage imagecontent.InformationImage) *fyne.MenuItem {
	return fyne.NewMenuItem("Brightness and Contrast", func() {
		newWindows := newwindow.NewWindow(application, 500, 500,
			"Brightness and Contrast")
		data, data2 := binding.NewFloat(), binding.NewFloat()
		data.Set(float64(int(fullImage.Brigthness())))
		slide := widget.NewSliderWithData(0, 255, data)
		formatted := binding.FloatToStringWithFormat(data, "Float value: %0.2f")
		label := widget.NewLabelWithData(formatted)

		data2.Set(float64(int(fullImage.Contrast())))
		slide2 := widget.NewSliderWithData(0, 127, data2)
		formatted2 := binding.FloatToStringWithFormat(data2, "Float value: %0.2f")
		label2 := widget.NewLabelWithData(formatted2)

		content := widget.NewButton("Calculate", func() {
			bright, _ := data.Get()
			conts, _ := data2.Get()
			newFullImage :=
				imagecontent.New(operations.AdjustBrightnessAndContrast(fullImage, bright, conts),
					fullImage.LutGray(), fullImage.Format())
			generalMenu(application, newFullImage, "Modified Image")
		})

		brightnessText, contrastText := canvas.NewText("Brightness", color.White),
			canvas.NewText("Contrast", color.White)
		brightnessText.TextStyle, contrastText.TextStyle =
			fyne.TextStyle{Bold: true}, fyne.TextStyle{Bold: true}
		menuAndImageContainer := container.NewVBox(brightnessText, label, slide,
			contrastText, label2, slide2, content)

		newWindows.SetContent(menuAndImageContainer)
		newWindows.Show()
	})
}

func negativeButton(application fyne.App,
	content imagecontent.InformationImage) *fyne.MenuItem {
	return fyne.NewMenuItem("Negative", func() {
		fullImage := imagecontent.New(operations.Negative(content, content.LutGray()),
			content.LutGray(), content.Format())
		generalMenu(application, fullImage, "Negative")
	})
}

func differenceButton(application fyne.App, window fyne.Window,
	content imagecontent.InformationImage) *fyne.MenuItem {
	return fyne.NewMenuItem("Image difference", func() {
		newDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				fileName := reader.URI().String()[7:]
				image, _, err := loadandsave.LoadImage(fileName)
				width := image.Bounds().Dx()
				height := image.Bounds().Dy()
				if err != nil {
					dialog.ShowError(err, window)
				} else if content.Image().Bounds().Dx() != width ||
					content.Image().Bounds().Dy() != height {
					dialog.ShowError(errors.New("the image must have the same dimensions"), window)
				} else {
					newWindow := newwindow.NewWindow(application, width, height, "Used Image")
					canvasImage := canvas.NewImageFromImage(image)
					newWindow.SetContent(canvasImage)
					newWindow.Show()
					difference, err := operations.ImageDifference(content, image)
					if err != nil {
						dialog.ShowError(err, window)
					} else {
						generalMenu(application, imagecontent.New(difference, content.LutGray(),
							content.Format()), "Difference")
					}
				}
			}
		}, window)
		newDialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".png",
			".jpeg", ".tiff"}))
		newDialog.Show()
	})
}

func changeMapButton(application fyne.App, window fyne.Window,
	content imagecontent.InformationImage) *fyne.MenuItem {
	return fyne.NewMenuItem("Change Map", func() {
		newDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				fileName := reader.URI().String()[7:]
				image, _, err := loadandsave.LoadImage(fileName)
				width := image.Bounds().Dx()
				height := image.Bounds().Dy()
				if err != nil {
					dialog.ShowError(err, window)
				} else if content.Image().Bounds().Dx() != width ||
					content.Image().Bounds().Dy() != height {
					dialog.ShowError(errors.New("the image must have the same dimensions"), window)
				} else {
					newWindow := newwindow.NewWindow(application, width, height, "Used Image")
					canvasImageUsed := canvas.NewImageFromImage(image)
					newWindow.SetContent(canvasImageUsed)
					newWindow.Show()

					difference, err := operations.ImageDifference(content, image)
					if err != nil {
						dialog.ShowError(err, window)
					} else {

						differenceWindow := newwindow.NewWindow(application,
							difference.Bounds().Dx(), difference.Bounds().Dy(), "Difference")
						canvasImageDifference := canvas.NewImageFromImage(difference)
						differenceWindow.SetContent(canvasImageDifference)

						fullImage := imagecontent.New(difference, content.LutGray(),
							content.Format())

						histogramItem := histogramButton(application, window, fullImage,
							"Histogram", false)

						cumulativeHistogramItem := histogramButton(application, window,
							fullImage, "Cumulative Histogram", true)

						thresHoldItem := thresHoldButton(application, difference, image)

						quitItem := quitButton(differenceWindow)

						separatorItem := fyne.NewMenuItemSeparator()

						menuItem := fyne.NewMenu("File", saveButton(application,
							difference), separatorItem, quitItem)
						menuItem2 := fyne.NewMenu("User value", thresHoldItem)
						menuItem3 := fyne.NewMenu("Histograms", histogramItem, separatorItem,
							cumulativeHistogramItem)
						menu := fyne.NewMainMenu(menuItem, menuItem2, menuItem3)
						differenceWindow.SetMainMenu(menu)
						differenceWindow.Show()
					}
				}
			}
		}, window)
		newDialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".png",
			".jpeg", ".tiff"}))
		newDialog.Show()
	})
}

func sectionsButton(application fyne.App,
	fullImage imagecontent.InformationImage) *fyne.MenuItem {
	return fyne.NewMenuItem("Lineal Adjustment in Sections", func() {
		windowSections := newwindow.NewWindow(application, 500, 200, "Sections Number")
		input := widget.NewEntry()
		input.SetPlaceHolder("0")
		content := container.NewVBox(input, widget.NewButton("Enter", func() {
			number, err := strconv.Atoi(input.Text)
			if err != nil {
				dialog.ShowError(err, windowSections)
			} else {
				windowValues := newwindow.NewWindow(application, 500, 500, "Sections Values")

				label1, label2 := widget.NewLabel("Coordinates X: "),
					widget.NewLabel("Coordinates Y: ")
				coordinatesX, coordinatesY := container.NewVBox(label1),
					container.NewVBox(label2)

				var point, point2 *widget.Entry
				var entries []pairEntry

				for i := 0; i < number+1; i++ {
					point, point2 = widget.NewEntry(), widget.NewEntry()
					point.SetPlaceHolder("x:")
					point2.SetPlaceHolder("y:")
					coordinatesX.Add(point)
					coordinatesY.Add(point2)
					entries = append(entries, pairEntry{x: point, y: point2})
				}
				var defaultGraph map[int]int
				histogram.Plotesections(defaultGraph)
				canvasImage := canvas.NewImageFromFile(".tmp/sectHist.png")

				content := container.NewVBox(container.NewHBox(coordinatesX, coordinatesY))

				button := func(window fyne.Window) {
					var coordinates []operations.Pair
					plott := make(map[int]int)
					for i := 0; i < len(entries); i++ {
						pointX, _ := strconv.Atoi(entries[i].x.Text)
						pointY, _ := strconv.Atoi(entries[i].y.Text)
						if i != len(entries)-1 {
							pointXplus, _ := strconv.Atoi(entries[i+1].x.Text)
							if pointX == pointXplus {
								dialog.ShowError(errors.New("the values of the X axis of the points can not be repeated"),
									windowValues)
								return
							}
						}
						if pointX < 0 || pointY < 0 || pointX > 255 || pointY > 255 {
							dialog.ShowError(errors.New("the points must be between 0 and 255"),
								windowValues)
						}
						coordinates = append(coordinates, operations.Pair{X: pointX, Y: pointY})
						plott[pointX] = pointY
					}
					histogram.Plotesections(plott)
					window.Content().Refresh()
					newFullImage := imagecontent.New(operations.LinealAdjustmentInSections(fullImage,
						coordinates, number), fullImage.LutGray(),
						fullImage.Format())
					generalMenu(application, newFullImage, "Sections Result")
				}

				windowSections.Close()
				windowValues.SetContent(container.NewBorder(content,
					widget.NewButton("Enter", func() { button(windowValues) }), nil, nil, canvasImage))
				windowValues.Show()
			}
		}))
		windowSections.SetContent(content)
		windowSections.Show()
	})
}

func histogramButton(application fyne.App, window fyne.Window,
	content imagecontent.InformationImage, name string,
	cumulative bool) *fyne.MenuItem {
	return fyne.NewMenuItem(name, func() {
		histogram.Plote(content.NumbersOfPixel(), content.Values(), cumulative)
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

func roiButton(application fyne.App,
	fullImage imagecontent.InformationImage) *fyne.MenuItem {
	return fyne.NewMenuItem("Region of interest", func() {
		width := fullImage.Image().Bounds().Dx()
		height := fullImage.Image().Bounds().Dy()

		roiWindow := newwindow.NewWindow(application, 300, 300, "Values of points")

		label1, label2 := widget.NewLabel("Start Point: "),
			widget.NewLabel("Final Point: ")
		point1I, point1J, point2I, point2J := widget.NewEntry(), widget.NewEntry(),
			widget.NewEntry(), widget.NewEntry()
		point1I.SetPlaceHolder("i1:")
		point1J.SetPlaceHolder("j1:")
		point2I.SetPlaceHolder("i2:")
		point2J.SetPlaceHolder("j2:")

		initialPoint := container.NewVBox(label1, point1I, point1J)
		finalPoint := container.NewVBox(label2, point2I, point2J)
		content := container.NewVBox(container.NewHBox(initialPoint, finalPoint),
			widget.NewButton("Save", func() {
				point1IInt, _ := strconv.Atoi(point1I.Text)
				point1JInt, _ := strconv.Atoi(point1J.Text)
				point2IInt, _ := strconv.Atoi(point2I.Text)
				point2JInt, _ := strconv.Atoi(point2J.Text)
				if point1IInt < 0 || point1JInt < 0 || point2IInt < 0 || point2JInt < 0 {
					dialog.ShowError(errors.New("the i and j values must be positive"),
						roiWindow)
				} else if point1IInt > width || point1JInt > height ||
					point2IInt > width || point2JInt > height {
					dialog.ShowError(errors.New("The i value must be lower than "+
						strconv.Itoa(width)+" and j value must be lower than "+
						strconv.Itoa(height)),
						roiWindow)
				} else {
					newFullImage := imagecontent.New(operations.ROI(fullImage, point1IInt, point1JInt,
						point2IInt, point2JInt), fullImage.LutGray(), fullImage.Format())
					generalMenu(application, newFullImage, "ROI")
				}
			}))
		roiWindow.SetContent(content)
		roiWindow.Show()
	})
}

func saveButton(application fyne.App, image image.Image) *fyne.MenuItem {
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

func equalizationButton(application fyne.App, content imagecontent.InformationImage,
	input string) *fyne.MenuItem {
	return fyne.NewMenuItem("Equalization", func() {
		fullImage := imagecontent.New(operations.EqualizeAnImage(content),
			content.LutGray(), content.Format())
		generalMenu(application, fullImage, input)
	})
}

func histogramSpecificationButton(application fyne.App, window fyne.Window,
	content imagecontent.InformationImage) *fyne.MenuItem {
	return fyne.NewMenuItem("Histogram Specification", func() {
		newDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				fileName := reader.URI().String()[7:]
				refImage, _, err := loadandsave.LoadImage(fileName)
				if err != nil {
					dialog.ShowError(err, window)
				}
				refImage2 := operations.ScaleGray(refImage)
				fullImage := imagecontent.New(refImage2, content.LutGray(), content.Format())
				generalMenu(application, fullImage, "Used Image")
				generalMenu(application,
					operations.HistogramSpecification(refImage2, content), "Result Specication")
			}
		}, window)
		newDialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".png",
			".jpeg", ".tiff"}))
		newDialog.Show()
	})
}

func thresHoldButton(application fyne.App, grayImage *image.Gray,
	img image.Image) *fyne.MenuItem {
	return fyne.NewMenuItem("Threshold", func() {
		windowThreshold := newwindow.NewWindow(application, 500, 200, "Threshold Value")
		data := binding.NewFloat()
		data.Set(0)
		slide := widget.NewSliderWithData(0, 255, data)
		formatted := binding.FloatToStringWithFormat(data, "Float value: %0.2f")
		label := widget.NewLabelWithData(formatted)

		content := widget.NewButton("Calculate", func() {
			threshold, _ := data.Get()
			newImage := operations.ChangeMap(grayImage, img, threshold)
			windowResult := newwindow.NewWindow(application,
				newImage.Bounds().Dx(), newImage.Bounds().Dy(), "Result")
			canvasResult := canvas.NewImageFromImage(newImage)

			quitItem := quitButton(windowResult)
			separatorItem := fyne.NewMenuItemSeparator()
			menuItem := fyne.NewMenu("File", saveButton(application,
				newImage), separatorItem, quitItem)
			menu := fyne.NewMainMenu(menuItem)
			windowResult.SetMainMenu(menu)

			windowResult.SetContent(canvasResult)
			windowResult.Show()
		})

		threshold := canvas.NewText("Threshold", color.White)
		threshold.TextStyle = fyne.TextStyle{Bold: true}
		menuAndImageContainer := container.NewVBox(threshold, label, slide,
			content)

		windowThreshold.SetContent(menuAndImageContainer)
		windowThreshold.Show()
	})
}

func quitButton(window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem("Quit", func() {
		window.Close()
	})
}

func ButtonOpen(application fyne.App, window fyne.Window) *fyne.MenuItem {
	return fyne.NewMenuItem("Open image", func() {
		newDialog := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if reader != nil {
				fileName := reader.URI().String()[7:]
				colorImage, format, err := loadandsave.LoadImage(fileName)
				if err != nil {
					dialog.ShowError(err, window)
				} else {
					width := colorImage.Bounds().Dx()
					height := colorImage.Bounds().Dy()
					grayImage := operations.ScaleGray(colorImage)

					lutGray := operations.LutGray()

					fullImage := imagecontent.New(grayImage, lutGray, format)

					informationTape := information.Information(fullImage)

					windowName := strings.Split(fileName, "/")
					imageWindow := newwindow.NewWindow(application,
						colorImage.Bounds().Dx(), colorImage.Bounds().Dy(),
						windowName[len(windowName)-1])
					canvasImage := canvas.NewImageFromImage(colorImage)
					text := strconv.Itoa(height) + " x " + strconv.Itoa(width)
					canvasText := canvas.NewText(text, color.Opaque)
					imageWindow.SetContent(container.NewBorder(nil, canvasText, nil, nil,
						canvasImage, mouse.New(colorImage, canvasText, text)))

					imageInformationItem := fyne.NewMenuItem("Image Information", func() {
						dialog.ShowInformation("Information", informationTape, imageWindow)
					})

					scaleGrayItem := fyne.NewMenuItem("Scale gray", func() {
						generalMenu(application, fullImage, windowName[len(windowName)-1])
					})

					quitItem := quitButton(imageWindow)

					separatorItem := fyne.NewMenuItemSeparator()

					saveItem := fyne.NewMenu("File", saveButton(application,
						grayImage), separatorItem, quitItem)

					imageInformationMenu := fyne.NewMenu("Image Information", imageInformationItem)
					operationItem := fyne.NewMenu("Operations", scaleGrayItem)
					menu := fyne.NewMainMenu(saveItem, imageInformationMenu, operationItem)
					imageWindow.SetMainMenu(menu)
					imageWindow.Show()
				}
			}
		}, window)
		newDialog.SetFilter(storage.NewExtensionFileFilter([]string{".jpg", ".png",
			".jpeg", ".tiff"}))
		newDialog.Show()
	})
}

type pairEntry struct {
	x, y *widget.Entry
}
