package cosplay

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

// Plugin cosplay图片
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "cosplay",
		Name: "cosplay图片",
	}
}

// IsFireEvent 是否触发
func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && field.Content == ".cosplay"
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{
		Elements: make([]message.IMessageElement, 1),
	}
	b, err := getPic()
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

type response struct {
	Msg  string `json:"msg"`
	Code int    `json:"code"`
	Data []data `json:"data"`
}

type data struct {
	ID  string `json:"id"`
	Url string `json:"url"`
}

var url = "http://api.isoyu.com/api/picture/index?page=%v"

func getPic() ([]byte, error) {
	page := rand.Intn(190)
	r, err := http.DefaultClient.Get(fmt.Sprintf(url, page))
	if err != nil {
		return nil, err
	}
	robots, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	if robots == nil {
		return nil, errors.New("抓取列表失败")
	}
	r.Body.Close()
	resp := string(robots)
	if resp == "" {
		return nil, errors.New("请稍后重试")
	}

	var data response
	err = json.Unmarshal(robots, &data)
	if err != nil {
		return nil, err
	}
	return nil, nil
}
