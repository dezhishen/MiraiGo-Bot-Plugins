package image

import (
	"image/color"

	"github.com/fogleman/gg"
)

func CreatImage(text string, path string) error {
	//图片的宽度
	var srcWidth float64 = 200
	//图片的高度
	var srcHeight float64 = 200
	dc := gg.NewContext(int(srcWidth), int(srcHeight))
	//设置背景色
	dc.SetColor(color.White)
	dc.Clear()
	dc.SetRGB255(255, 0, 0)
	if err := dc.LoadFontFace("C:\\github\\plugins\\assert\\fonts\\yasqht.ttf", 25); err != nil {
		return err
	}
	sWidth, _ := dc.MeasureString(text)
	dc.DrawString(text, (srcWidth-sWidth)/2, 40)
	err := dc.SavePNG(path)
	if err != nil {
		return err
	}
	return nil
}
