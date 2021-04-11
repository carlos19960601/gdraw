package core

import "github.com/zengqiang96/gdraw/internal/numutil"

type Scanline struct {
	Y, X1, X2 int
	Alpha     uint32
}

func cropScanlines(lines []Scanline, w, h int) []Scanline {
	i := 0
	for _, line := range lines {
		if line.Y < 0 || line.Y >= h {
			continue
		}
		if line.X1 >= w {
			continue
		}
		line.X1 = numutil.ClampInt(line.X1, 0, w-1)
		line.X2 = numutil.ClampInt(line.X2, 0, w-1)
		if line.X1 > line.X2 {
			continue
		}
		lines[i] = line
		i++
	}
	return lines[:i]
}
