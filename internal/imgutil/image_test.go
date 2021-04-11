package imgutil

import "testing"

func TestLoadImage(t *testing.T) {
	img, err := LoadImage("/Users/zengqiang96/Pictures/壁纸/13kkm3.png")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v\n", img.Bounds())                          // 图片长宽
	t.Logf("%v %T\n", img.ColorModel(), img.ColorModel()) // 图片颜色模型
	t.Logf("%v\n", img.At(100, 100))                      // 该像素点的颜色
}
