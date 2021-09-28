package image

import (
	"image"
	"image/color"
	"image/png"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype"
)

func CreatImage(rp string) {
	//图片的宽度
	srcWidth := 200
	//图片的高度
	srcHeight := 200
	img := image.NewRGBA(image.Rect(0, 0, srcWidth, srcHeight))

	//为背景图片设置颜色
	for y := 0; y < srcWidth; y++ {
		for x := 0; x < srcHeight; x++ {
			img.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}
	f := freetype.NewContext()
	//设置分辨率
	f.SetDPI(100)
	//设置尺寸
	f.SetFontSize(26)
	//读取字体数据  http://fonts.mobanwang.com/201503/12436.html
	fontBytes, err := ioutil.ReadFile("C:\\github\\plugins\\assert\\fonts\\yasqht.ttf")
	if err != nil {
		log.Println(err)
	}
	//载入字体数据
	font, err := freetype.ParseFont(fontBytes)
	if err != nil {
		log.Println("载入字体失败", err)
	}
	f.SetFont(font)
	f.SetClip(img.Bounds())
	//设置输出的图片
	f.SetDst(img)
	//设置字体颜色(红色)
	f.SetSrc(image.NewUniform(color.RGBA{255, 0, 0, 255}))
	//设置字体的位置
	pt := freetype.Pt(55, 40)
	_, err = f.DrawString(rp, pt)
	if err != nil {
		log.Fatal(err)
	}
	//以png 格式写入文件

	imgfile, _ := os.Create("out.png")
	defer imgfile.Close()
	err = png.Encode(imgfile, img)
	if err != nil {
		log.Fatal(err)
	}
}
