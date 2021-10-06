package almanac

import (
	"bytes"
	"fmt"
	"image/color"
	"io/ioutil"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/fogleman/gg"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
)

// Plugin jrrp
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (p Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "jrrp",
		Name: "jrrp",
	}
}

// IsFireEvent 是否触发
func (p Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		if !ok {
			return false
		}
		if strings.HasPrefix(field.Content, "签到") {
			return true
		}
		if strings.HasPrefix(field.Content, ".jrrp") {
			return true
		}
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (p Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{
		Elements: make([]message.IMessageElement, 1),
	}
	b, err := getImage(request.Sender.Uin)
	if err != nil {
		return nil, err
	}
	var image message.IMessageElement
	if plugins.GroupMessage == request.MessageType {
		image, err = request.QQClient.UploadGroupImage(request.GroupCode, bytes.NewReader(b))
	} else {
		image, err = request.QQClient.UploadPrivateImage(request.Sender.Uin, bytes.NewReader(b))
	}
	if err != nil {
		return nil, err
	}
	result.Elements[0] = image
	return result, nil
}

func getImage(id int64) ([]byte, error) {
	timeNow := time.Now().Local()
	path := getFileName(id, timeNow)
	b, err := getFile(id, path)
	if err != nil {
		return nil, err
	}
	if b != nil {
		return b, err
	}
	err = randomFile(timeNow, "jrrp", id, true, path)
	if err != nil {
		return nil, err
	}
	return getFile(id, path)
}

func getFileName(id int64, t time.Time) string {
	return fmt.Sprintf("./jrrp/%v-%v.png", id, timeToStr(t))
}

func getFile(id int64, path string) ([]byte, error) {
	// url = strings.Replace(url, "large", "original", -1)
	exists, _ := pathExists(path)
	if exists {
		file, err := os.Open(path)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		return content, err
	}
	return nil, nil
}

func randomFile(t time.Time, pid string, uid int64, genIfNil bool, path string) error {
	r, err := getScore(t, pid, uid, true)
	if err != nil {
		return err
	}
	var score = r / 20
	full := "★"
	empty := "☆"
	var text string
	for i := 1; i <= score; i++ {
		text += full
	}
	for i := 1; i <= 5-score; i++ {
		text += empty
	}
	err = CreatImage(text, path)
	if err != nil {
		return err
	}
	return err
}

func getScore(t time.Time, pid string, uid int64, genIfNil bool) (int, error) {
	timestr := timeToStr(t)
	key := []byte(fmt.Sprintf("jrrp.%v.%v", uid, timestr))
	var score int
	err := storage.Get([]byte(pid), key, func(b []byte) error {
		if b != nil {
			score = storage.BytesToInt(b)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	if genIfNil && score == 0 {
		rand.Seed(time.Now().UnixNano())
		score = rand.Intn(100) + 1
		storage.Put([]byte(pid), key, storage.IntToBytes(score))
		theTime := t.Add(-7 * 24 * time.Hour)
		theTimestr := timeToStr(theTime)
		keyLast7Day := fmt.Sprintf("jrrp.%v.%v", uid, theTimestr)
		storage.Delete([]byte(pid), []byte(keyLast7Day))
	}
	return score, nil
}

func timeToStr(t time.Time) string {
	return fmt.Sprintf("%v-%v-%v", t.Year(), int(t.Month()), t.Day())
}

func init() {
	exists, _ := pathExists("./jrrp")
	if !exists {
		os.Mkdir("./jrrp", 0777)
	}
	plugins.RegisterOnMessagePlugin(Plugin{})
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

func CreatImage(text string, path string) error {
	//图片的宽度
	var srcWidth float64 = 100
	//图片的高度
	var srcHeight float64 = 100
	dc := gg.NewContext(int(srcWidth), int(srcHeight))
	//设置背景色
	dc.SetColor(color.White)
	dc.Clear()
	dc.SetRGB255(255, 0, 0)
	if err := dc.LoadFontFace("/assert/fonts/Symbola.ttf", 25); err != nil {
		return err
	}
	sWidth, sHeight := dc.MeasureString(text)
	dc.DrawString(text, (srcWidth-sWidth)/2, (srcHeight-sHeight)/2)
	err := dc.SavePNG(path)
	if err != nil {
		return err
	}
	return nil
}
