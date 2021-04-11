package core

import (
	"image"
	"math"

	"github.com/zengqiang96/gdraw/internal/color"
	"github.com/zengqiang96/gdraw/internal/numutil"
)

func DifferenceFull(one, another *image.RGBA) float64 {
	size := one.Bounds().Size()
	w, h := size.X, size.Y
	var total uint64
	for y := 0; y < h; y++ {
		i := one.PixOffset(0, y)
		for x := 0; x < w; x++ {
			oneR := int(one.Pix[i])
			oneG := int(one.Pix[i+1])
			oneB := int(one.Pix[i+2])
			oneA := int(one.Pix[i+3])
			anotherR := int(another.Pix[i])
			anotherG := int(another.Pix[i+1])
			anotherB := int(another.Pix[i+2])
			anotherA := int(another.Pix[i+3])
			i += 4
			dr := oneR - anotherR
			dg := oneG - anotherG
			db := oneB - anotherB
			da := oneA - anotherA
			total += uint64(dr*dr + dg*dg + db*db + da*da)
		}
	}
	return math.Sqrt(float64(total)/float64(w*h*4)) / 255
}

func ComputeColor(target, current *image.RGBA, lines []Scanline, alpha int) color.Color {
	var rsum, gsum, bsum, count int64
	a := 0x101 * 255 / alpha
	for _, line := range lines {
		i := target.PixOffset(line.X1, line.Y)
		for x := line.X1; x < line.X2; x++ {
			tr := int(target.Pix[i])
			tg := int(target.Pix[i+1])
			tb := int(target.Pix[i+2])
			cr := int(current.Pix[i])
			cg := int(current.Pix[i+1])
			cb := int(current.Pix[i+2])
			i += 4
			rsum += int64((tr-cr)*a + cr*0x101)
			gsum += int64((tg-cg)*a + cg*0x101)
			bsum += int64((tb-cb)*a + cb*0x101)
			count++
		}
	}

	if count == 0 {
		return color.Color{}
	}

	r := numutil.ClampInt(int(rsum/count)>>8, 0, 255)
	g := numutil.ClampInt(int(gsum/count)>>8, 0, 255)
	b := numutil.ClampInt(int(bsum/count)>>8, 0, 255)
	return color.Color{R: r, G: g, B: b, A: alpha}
}

func CopyLines(dst, src *image.RGBA, lines []Scanline) {
	for _, line := range lines {
		a := dst.PixOffset(line.X1, line.Y)
		b := a + (line.X2-line.X1+1)*4
		copy(dst.Pix[a:b], src.Pix[a:b])
	}
}

func DrawLines(img *image.RGBA, c color.Color, lines []Scanline) {
	const m = 0xffff
	sr, sg, sb, sa := c.NRGBA().RGBA()
	for _, line := range lines {
		ma := line.Alpha
		a := (m - sa*ma/m) * 0x101
		i := img.PixOffset(line.X1, line.Y)
		for x := line.X1; x <= line.X2; x++ {
			dr := uint32(img.Pix[i])
			dg := uint32(img.Pix[i+1])
			db := uint32(img.Pix[i+2])
			da := uint32(img.Pix[i+3])
			img.Pix[i] = uint8((dr*a + sr*ma) / m >> 8)
			img.Pix[i+1] = uint8((dg*a + sg*ma) / m >> 8)
			img.Pix[i+2] = uint8((db*a + sb*ma) / m >> 8)
			img.Pix[i+3] = uint8((da*a + sa*ma) / m >> 8)
			i += 4
		}
	}
}

func DifferencePartial(target, before, after *image.RGBA, score float64, lines []Scanline) float64 {
	size := target.Bounds().Size()
	w, h := size.X, size.Y
	total := uint64(math.Pow(score*255, 2) * float64(w*h*4))
	for _, line := range lines {
		i := target.PixOffset(line.X1, line.Y)
		for x := line.X1; x < line.X2; x++ {
			tr := int(target.Pix[i])
			tg := int(target.Pix[i+1])
			tb := int(target.Pix[i+2])
			ta := int(target.Pix[i+3])
			br := int(before.Pix[i])
			bg := int(before.Pix[i+1])
			bb := int(before.Pix[i+2])
			ba := int(before.Pix[i+3])
			ar := int(after.Pix[i])
			ag := int(after.Pix[i+1])
			ab := int(after.Pix[i+2])
			aa := int(after.Pix[i+3])
			i += 4
			dr1 := tr - br
			dg1 := tg - bg
			db1 := tb - bb
			da1 := ta - ba
			dr2 := tr - ar
			dg2 := tg - ag
			db2 := tb - ab
			da2 := ta - aa
			total -= uint64(dr1*dr1 + dg1*dg1 + db1*db1 + da1*da1)
			total += uint64(dr2*dr2 + dg2*dg2 + db2*db2 + da2*da2)
		}
	}
	return math.Sqrt(float64(total)/float64(w*h*4)) / 255
}
