package information

import (
	"strconv"
)

func Information(format string, width, height, min, max int, brightness, contrast float64,
	entropy float64) string {
	information := "Format: " + format + "\nSize:\n - Width: " +
		strconv.Itoa(width) + "\n - Height: " + strconv.Itoa(height) +
		"\nScale Gray Range:\n - Min: " + strconv.Itoa(min) + "\n - Max: " +
		strconv.Itoa(max) + "\nEntropy: " + strconv.Itoa(int(entropy)) +
		"\nBrightness: " + strconv.Itoa(int(brightness)) + "\nContrast: " +
		strconv.Itoa(int(contrast))
	return information

}
