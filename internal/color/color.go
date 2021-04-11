package color

import (
	"fmt"
	"image/color"
	"strings"
)

type Color struct {
	R, G, B, A int
}

func (c Color) NRGBA() color.NRGBA {
	return color.NRGBA{uint8(c.R), uint8(c.G), uint8(c.B), uint8(c.A)}
}

func MakeColor(c color.Color) Color {
	r, g, b, a := c.RGBA()
	return Color{int(r / 257), int(g / 257), int(b / 257), int(a / 257)}
}

func MakeHexColor(hex string) Color {
	hex = strings.Trim(hex, "#")
	var r, g, b, a int
	a = 255
	switch len(hex) {
	case 3:
		fmt.Sscanf(hex, "%1x%1x%1x", &r, &g, &b)
		r |= r << 4
		g |= g << 4
		b |= b << 4
	case 4:
		fmt.Sscanf(hex, "%1x%1x%1x%1x", &r, &g, &b, &a)
		r |= r << 4
		g |= g << 4
		b |= b << 4
		a |= a << 4
	case 6:
		fmt.Sscanf(hex, "%2x%2x%2x", &r, &g, &b)
	case 8:
		fmt.Sscanf(hex, "%2x%2x%2x%2x", &r, &g, &b, &a)
	}
	return Color{r, g, b, a}
}
