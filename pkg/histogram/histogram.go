package histogram

import (
	"math"
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
	number := numbersOfPixel
	if cumulative {
		number = CumulativeHistogram(numbersOfPixel)
	}
	Xaxis := []float64{}
	Yaxis := []float64{}
	// fmt.Println(len(values))
	value := chart.ContinuousSeries{}
	for i := 0; i < len(number); i++ {
		Yaxis = append(Yaxis, float64(float64(number[i])/float64(len(values))))
		Xaxis = append(Xaxis, float64(i))
	}
	if Yaxis[0] >= 0.01 {
		Yaxis[0] -= 0.01
	}
	value = chart.ContinuousSeries{
		XValues: Xaxis,
		YValues: Yaxis,
	}

	graph := chart.Chart{
		Series: []chart.Series{
			value,
		},
	}

	ouputFile, _ := os.Create(".tmp/hist.png")
	graph.Render(chart.PNG, ouputFile)
	ouputFile.Close()
}

func Plotesections(numbersOfPixel map[int]int) {

	Xaxis := []float64{}
	Yaxis := []float64{}
	// fmt.Println(len(values))
	value := chart.ContinuousSeries{}
	for index, element := range numbersOfPixel {
		Yaxis = append(Yaxis, float64(float64(element)))
		Xaxis = append(Xaxis, float64(index))
	}

	value = chart.ContinuousSeries{
		XValues: Xaxis,
		YValues: Yaxis,
	}

	graph := chart.Chart{
		Series: []chart.Series{
			value,
		},
	}

	ouputFile, _ := os.Create(".tmp/sectHist.png")
	graph.Render(chart.PNG, ouputFile)
	ouputFile.Close()
}

func Equalization(histogram map[int]int, size int) map[int]int {
	cumulativeHistogram := CumulativeHistogram(histogram)
	ecualizeHistogram := make(map[int]int, 255)

	for i := 0; i < len(cumulativeHistogram); i++ {
		mul := float64(cumulativeHistogram[i]) * float64(256.0)
		cociente := float64(mul / float64(size))
		ecualizeHistogram[i] = int(math.Max(float64(0), math.Round(float64(cociente)-1.0)))
	}
	return ecualizeHistogram
}
