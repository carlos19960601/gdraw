package imgutil

import (
	"fmt"
	"image"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"image/color"
	"image/draw"
	"image/png"
	_ "image/png"
)

func LoadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, _, err := image.Decode(file)
	return img, err
}

// 所有像素点的平均值
func AverageImageColor(img image.Image) color.NRGBA {
	rgba := Image2RGBA(img)
	size := rgba.Bounds().Size()
	w, h := size.X, size.Y
	var r, g, b int
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			c := rgba.RGBAAt(x, y)
			r += int(c.R)
			g += int(c.G)
			b += int(c.B)
		}
	}
	r /= w * h
	g /= w * h
	b /= w * h
	return color.NRGBA{uint8(r), uint8(g), uint8(b), 255}
}

func Image2RGBA(src image.Image) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	draw.Draw(dst, dst.Rect, src, image.Point{}, draw.Src)
	return dst
}

func UniformRGBA(r image.Rectangle, c color.Color) *image.RGBA {
	img := image.NewRGBA(r)
	draw.Draw(img, img.Bounds(), &image.Uniform{c}, image.Point{}, draw.Src)
	return img
}

func CopyRGBA(src *image.RGBA) *image.RGBA {
	dst := image.NewRGBA(src.Bounds())
	copy(dst.Pix, src.Pix)
	return dst
}

func SavePNG(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return png.Encode(file, img)
}

func SaveGIFImageMagick(path string, frames []image.Image, delay, lastDelay int) error {
	dir, err := ioutil.TempDir("", "")
	if err != nil {
		return err
	}
	for i, im := range frames {
		path := filepath.Join(dir, fmt.Sprintf("%06d.png", i))
		_ = SavePNG(path, im)
	}
	args := []string{
		"-loop", "0",
		"-delay", fmt.Sprint(delay),
		filepath.Join(dir, "*.png"),
		"-delay", fmt.Sprint(lastDelay - delay),
		filepath.Join(dir, fmt.Sprintf("%06d.png", len(frames)-1)),
		path,
	}
	cmd := exec.Command("convert", args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return os.RemoveAll(dir)
}
