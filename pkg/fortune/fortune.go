package fortune

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"image/jpeg"
	"io/ioutil"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fogleman/gg"
)

type FortuneResult struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

var root = "./fortune"
var site = "https://pan.dihe.moe/fortune"

// @function randtext 随机选取签文
// @param file 文件路径
// @param seed 随机数种子
// @return 运势结果 & 错误信息
func Randtext() (*FortuneResult, error) {
	file := "运势签文.json"
	seed := time.Now().UnixNano()
	path, err := getFile(file)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	temp := []map[string]string{}
	if err := json.Unmarshal(data, &temp); err != nil {
		return nil, err
	}
	rand.Seed(seed)
	r := rand.Intn(len(temp))
	return &FortuneResult{
		Title:   temp[r]["title"],
		Content: temp[r]["content"],
	}, nil
}

// @function draw 绘制运势图
// @param background 背景图片路径
// @param title 签名
// @param text 签文
// @return 错误信息
func Draw(background string, fortuneResult *FortuneResult) ([]byte, error) {
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
	sw, _ := canvas.MeasureString(fortuneResult.Title)
	canvas.DrawString(fortuneResult.Title, 140-sw/2, 112)
	// 写正文
	canvas.SetRGB(0, 0, 0)
	// if err := canvas.LoadFontFace(base+"sakura.ttf", 23); err != nil {
	// 	return nil, err
	// }
	tw, th := canvas.MeasureString("测")
	tw, th = tw+10, th+10
	r := []rune(fortuneResult.Content)
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

func getFileName(name string) string {
	i := strings.LastIndex(name, "/")
	return fmt.Sprintf(root+"/%v", name[i+1:])
}
func getSiteUrl(name string) string {
	i := strings.LastIndex(name, "/")
	return fmt.Sprintf(site+"/%v", name[i+1:])
}

func getFile(name string) (string, error) {
	path := getFileName(name)
	exists, _ := pathExists(path)
	if exists {
		return path, nil
	}
	err := downloadFile(name)
	if err != nil {
		return "", err
	}
	return path, nil
}

func downloadFile(name string) error {
	r, err := http.DefaultClient.Get(getSiteUrl(name))
	if err != nil {
		return err
	}
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	_ = ioutil.WriteFile(getFileName(name), content, 0644)
	return nil
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func init() {
	exists, _ := pathExists(root)
	if !exists {
		os.Mkdir(root, 0777)
	}
}
