package information

import (
	"fmt"
	"strconv"
	imagecontent "vpc/pkg/imageContent"
)

func Information(content imagecontent.InformationImage) string {
	information := "Format: " + content.Format() + "\nSize:\n - Width: " +
		strconv.Itoa(content.Image().Bounds().Dx()) + "\n - Height: " +
		strconv.Itoa(content.Image().Bounds().Dy()) + "\nScale Gray Range:\n - Min: " +
		strconv.Itoa(content.Min()) + "\n - Max: " + strconv.Itoa(content.Max()) +
		"\nEntropy: " + fmt.Sprintf("%f", content.Entropy()) + "\nNumber of bits: " +
		strconv.Itoa(int(content.Entropy()+1)) + "\nBrightness: " +
		strconv.Itoa(int(content.Brigthness())) + "\nContrast: " +
		strconv.Itoa(int(content.Contrast()))
	return information
}
