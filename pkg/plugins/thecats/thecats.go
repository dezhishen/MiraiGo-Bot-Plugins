package thecats

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

var url = "https://api.thecatapi.com/v1/images/search?limit=1&size=med"

// Plugin 猫猫图
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "thecat",
		Name: "猫猫图",
	}
}

// IsFireEvent 是否触发
func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && field.Content == ".thecat"
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{
		Elements: make([]message.IMessageElement, 1),
	}
	b, err := getCatPic()
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

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}

type Data struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

func getCatPic() ([]byte, error) {
	r, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	robots, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body.Close()
	resp := string(robots)
	if resp == "" {
		return nil, errors.New("请稍后重试")
	}

	var datas []Data
	err = json.Unmarshal(robots, &datas)
	if err != nil {
		return nil, err
	}
	if datas == nil || len(datas) == 0 {
		return nil, errors.New("请稍后重试")
	}
	imageResp, err := http.DefaultClient.Get(datas[0].Url)
	if err != nil {
		return nil, err
	}
	imageBytes, err := ioutil.ReadAll(imageResp.Body)
	if err != nil {
		return nil, err
	}
	imageResp.Body.Close()
	return imageBytes, nil
}
