package util

import (
	"image"
	"os"

	"go.uber.org/zap"
)

func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	zap.S().Error(err)
	return img, err
}
