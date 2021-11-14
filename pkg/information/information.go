package information

import (
	"fmt"
	"strconv"
)

func Information(format string, width, height, min, max int, brightness,
	contrast, entropy float64) string {
	information := "Format: " + format + "\nSize:\n - Width: " +
		strconv.Itoa(width) + "\n - Height: " + strconv.Itoa(height) +
		"\nScale Gray Range:\n - Min: " + strconv.Itoa(min) + "\n - Max: " +
		strconv.Itoa(max) + "\nEntropy: " + fmt.Sprintf("%f", entropy) +
		"\nNumber of bits: " + strconv.Itoa(int(entropy+1)) + "\nBrightness: " +
		strconv.Itoa(int(brightness)) + "\nContrast: " + strconv.Itoa(int(contrast))
	return information
}
