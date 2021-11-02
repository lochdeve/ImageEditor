package histogram

import (
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func NumbersOfPixel(colors []uint64) map[int]int {
	numbersOfPixel := make(map[int]int)
	for i := 0; i <= 255; i++ {
		cont := 0
		for j := 0; j < len(colors); j++ {
			if i == int(colors[j]) {
				cont++
			}
		}
		numbersOfPixel[i] = cont
	}
	// fmt.Println(numbersofpixel)
	return numbersOfPixel
}

func Plote(numbersOfPixel map[int]int, values plotter.Values) {

	p := plot.New()

	p.Title.Text = "Histogram plot"

	hist, err2 := plotter.NewHist(values, len(numbersOfPixel))
	if err2 != nil {
		panic(err2)
	}
	hist.Normalize(1)
	p.Add(hist)
	if err := p.Save(3*vg.Inch, 3*vg.Inch, ".tmp/hist.png"); err != nil {
		panic(err)
	}
}