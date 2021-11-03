package histogram

import (
	"fmt"
	"os"

	"github.com/wcharczuk/go-chart/v2"
	"gonum.org/v1/plot/plotter"
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

func CumulativeHistogram(histogram map[int]int) map[int]int {
	cumulativeHistogram := make(map[int]int)
	cumulative := 0
	for i := 0; i < len(histogram); i++ {
		cumulative += histogram[i]
		cumulativeHistogram[i] = cumulative
	}
	return cumulativeHistogram
}

func Plote(numbersOfPixel map[int]int, values plotter.Values, cumulative bool) {

	/*p := plot.New()

	p.Title.Text = "Histogram plot"

	hist, err2 := plotter.NewHist(values, len(numbersOfPixel))
	if err2 != nil {
		panic(err2)
	}
	hist.Normalize(1)
	p.Add(hist)
	if err := p.Save(3*vg.Inch, 3*vg.Inch, ".tmp/hist2.png"); err != nil {
		panic(err)
	}*/
	number := numbersOfPixel
	if cumulative {
		number = CumulativeHistogram(numbersOfPixel)
	}

	barr := make([]chart.Value, 0)
	fmt.Println(len(values))
	for i := 0; i < len(number); i++ {
		value := chart.Value{
			Value: float64(float64(number[i]) / float64(len(values))),
		}
		barr = append(barr, value)
	}
	graph := chart.BarChart{
		Title: "Histogram",
		Background: chart.Style{
			Padding: chart.Box{
				Top: 40,
			},
		},
		Height:   600,
		Width:    900,
		BarWidth: 1,
		Bars:     barr,
	}

	f, _ := os.Create(".tmp/hist.png")
	defer f.Close()
	graph.Render(chart.PNG, f)
	/*
		graph2 := chart.BarChart{
			Title: "Test Bar Chart",
			Background: chart.Style{
				Padding: chart.Box{
					Top: 20,
				},
			},
			BarSpacing: 0,
			Height:     600,
			Width:      1000,
			BarWidth:   1,
			Bars: []chart.Value{
				{Value: 5.25, Label: "Blue", Style: chart.Style{StrokeColor: chart.ColorAlternateGray}},
				{Value: 4.88, Label: "Green", Style: chart.Style{StrokeColor: chart.ColorAlternateGray}},
				{Value: 4.74, Label: "Gray", Style: chart.Style{StrokeColor: chart.ColorAlternateGray}},
				{Value: 3.22, Label: "Orange", Style: chart.Style{StrokeColor: chart.ColorAlternateGray}},
				{Value: 3, Label: "Test", Style: chart.Style{StrokeColor: chart.ColorAlternateGray}},
				{Value: 2.27, Label: "??", Style: chart.Style{StrokeColor: chart.ColorAlternateGray}},
				{Value: 1, Label: "!!"},
			},
		}

		f1, _ := os.Create("output.png")
		defer f.Close()
		graph2.Render(chart.PNG, f1)*/
}
