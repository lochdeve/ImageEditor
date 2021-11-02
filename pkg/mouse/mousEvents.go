package mouse

import (
	"image"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
)

type MouseEvents struct {
	widget.Icon
	image image.Image
	text  *canvas.Text
	size  string
}

func New(image1 image.Image, text1 *canvas.Text, text2 string) *MouseEvents {
	m := &MouseEvents{image: image1, text: text1, size: text2}
	m.ExtendBaseWidget(m)
	return m
}

func (t *MouseEvents) Tapped(_ *fyne.PointEvent) {
	//fmt.Println("I have been tapped")
}

func (w *MouseEvents) FocusGained() {
	//fmt.Println("FocusGained")
}

// FocusLost is a hook called by the focus handling logic after this object lost the focus.
func (w *MouseEvents) FocusLost() {
	//fmt.Println("Lost focus")
}

// TypedRune is a hook called by the input handling logic on text input events if this object is focused.
func (w *MouseEvents) TypedRune(_ rune) {

}

// TypedKey is a hook called by the input handling logic on key events if this object is focused.
func (w *MouseEvents) TypedKey(_ *fyne.KeyEvent) {

}

// MouseIn is a hook that is called if the mouse pointer enters the element.
func (w *MouseEvents) MouseIn(*desktop.MouseEvent) {
	//	fmt.Println("Inside")
}

// MouseMoved is a hook that is called if the mouse pointer moved over the element.
func (w *MouseEvents) MouseMoved(mousePosition *desktop.MouseEvent) {
	r, g, b, _ := w.image.At(int(mousePosition.AbsolutePosition.X),
		int(mousePosition.AbsolutePosition.Y)).RGBA()
	r, g, b = r>>8, g>>8, b>>8
	textaux := w.size + " R:" + strconv.Itoa(int(r)) + " G:" + strconv.Itoa(int(g)) +
		" B:" + strconv.Itoa(int(b))
	w.text.Text = textaux
	w.text.Refresh()
}

// MouseOut is a hook that is called if the mouse pointer leaves the element.
func (w *MouseEvents) MouseOut() {
	//fmt.Println("get out")
}

func (mouse *MouseEvents) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(canvas.NewImageFromImage(mouse.image))
}
