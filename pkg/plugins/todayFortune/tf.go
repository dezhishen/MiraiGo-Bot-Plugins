package todayFortune

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image/jpeg"
	"io/ioutil"
	"math"
	"math/rand"

	"github.com/fogleman/gg"
)

// @function randtext 随机选取签文
// @param file 文件路径
// @param seed 随机数种子
// @return 签名 & 签文 & 错误信息
func randtext(file string, seed int64) (string, string, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return "", "", err
	}
	temp := []map[string]string{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return "", "", err
	}
	rand.Seed(seed)
	r := rand.Intn(len(temp))
	return temp[r]["title"], temp[r]["content"], nil
}

// @function draw 绘制运势图
// @param background 背景图片路径
// @param title 签名
// @param text 签文
// @return 错误信息
func draw(background, title, text string) ([]byte, error) {
	// 加载背景
	back, err := gg.LoadImage(background)
	if err != nil {
		return nil, err
	}
	canvas := gg.NewContext(back.Bounds().Size().Y, back.Bounds().Size().X)
	canvas.DrawImage(back, 0, 0)
	// 写标题
	canvas.SetRGB(1, 1, 1)
	// if err := canvas.LoadFontFace(base+"sakura.ttf", 45); err != nil {
	// 	return nil, err
	// }
	sw, _ := canvas.MeasureString(title)
	canvas.DrawString(title, 140-sw/2, 112)
	// 写正文
	canvas.SetRGB(0, 0, 0)
	// if err := canvas.LoadFontFace(base+"sakura.ttf", 23); err != nil {
	// 	return nil, err
	// }
	tw, th := canvas.MeasureString("测")
	tw, th = tw+10, th+10
	r := []rune(text)
	xsum := rowsnum(len(r), 9)
	switch xsum {
	default:
		for i, o := range r {
			xnow := rowsnum(i+1, 9)

			offIt := (float64)(len(r) - (xnow-1)*9)
			ysum := (int)(math.Min(offIt, 9))
			ynow := i%9 + 1
			canvas.DrawString(string(o), -offest(xsum, xnow, tw)+115, offest(ysum, ynow, th)+320.0)
		}
	case 2:
		div := rowsnum(len(r), 2)
		for i, o := range r {
			xnow := rowsnum(i+1, div)

			offIt := (float64)(len(r) - (xnow-1)*div)
			flt64Div := (float64)(div)
			ysum := (int)(math.Min(offIt, flt64Div))
			ynow := i%div + 1
			switch xnow {
			case 1:
				canvas.DrawString(string(o), -offest(xsum, xnow, tw)+115, offest(9, ynow, th)+320.0)
			case 2:
				canvas.DrawString(string(o), -offest(xsum, xnow, tw)+115, offest(9, ynow+(9-ysum), th)+320.0)
			}
		}
	}
	// 转成 base64
	buffer := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	var opt jpeg.Options
	opt.Quality = 70
	err = jpeg.Encode(encoder, canvas.Image(), &opt)
	if err != nil {
		return nil, err
	}
	encoder.Close()
	return buffer.Bytes(), nil
}

func rowsnum(total, div int) int {
	temp := total / div
	if total%div != 0 {
		temp++
	}
	return temp
}

func offest(total, now int, distance float64) float64 {
	if total%2 == 0 {
		return (float64(now-total/2) - 1) * distance
	}
	return (float64(now-total/2) - 1.5) * distance
}
