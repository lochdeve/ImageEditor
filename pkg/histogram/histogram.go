package histogram

import (
	"math"
	"os"

	"github.com/wcharczuk/go-chart/v2"
)

func HistogramMap(colors []uint64) map[int]int {
	histogram := make(map[int]int)
	for i := 0; i <= 255; i++ {
		cont := 0
		for j := 0; j < len(colors); j++ {
			if i == int(colors[j]) {
				cont++
			}
		}
		histogram[i] = cont
	}
	return histogram
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

func Plote(histogramMap map[int]int, size int, cumulative bool) {
	number := histogramMap
	if cumulative {
		number = CumulativeHistogram(histogramMap)
	}
	Xaxis := []float64{}
	Yaxis := []float64{}
	// fmt.Println(len(values))
	value := chart.ContinuousSeries{}
	for i := 0; i < len(number); i++ {
		Yaxis = append(Yaxis, float64(float64(number[i])/float64(size)))
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

func Plotesections(histogramMap map[int]int) {
	Xaxis := []float64{}
	Yaxis := []float64{}
	// fmt.Println(len(values))
	value := chart.ContinuousSeries{}
	for index, element := range histogramMap {
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
	ecualizeLut := make(map[int]int)
	for i := 0; i < len(cumulativeHistogram); i++ {
		mult := float64(cumulativeHistogram[i]) * 256.0
		quotient := float64(mult / float64(size))
		ecualizeLut[i] = int(math.Max(0.0, math.Round(float64(quotient)-1.0)))
	}
	return ecualizeLut
}
