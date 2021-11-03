package histogram

import (
	"fmt"
	"os"
	"strconv"

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

	number := numbersOfPixel
	if cumulative {
		number = CumulativeHistogram(numbersOfPixel)
	}

	barr := make([]chart.Value, 0)
	fmt.Println(len(values))
	for i := 0; i < len(number); i++ {
		value := chart.Value{}
		if i%20 == 0 {
			value = chart.Value{
				Value: float64(float64(number[i]) / float64(len(values))),
				Label: strconv.Itoa(i),
			}
		} else {
			value = chart.Value{
				Value: float64(float64(number[i]) / float64(len(values))),
			}
		}
		barr = append(barr, value)
	}
	graph := chart.BarChart{
		Title: "Histogram",
		Background: chart.Style{
			Padding: chart.Box{
				Top:   40,
				Right: -100,
			},
		},
		BarSpacing:   5,
		ColorPalette: chart.DefaultColorPalette,
		Height:       600,
		Width:        810,
		BarWidth:     15,
		Bars:         barr,
	}

	f, _ := os.Create(".tmp/hist.png")
	graph.Render(chart.PNG, f)
	f.Close()
}
