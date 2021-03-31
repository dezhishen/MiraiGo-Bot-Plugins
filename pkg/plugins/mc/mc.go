package mc

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

// Plugin menhear
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   ".mc",
		Name: "menhear",
	}
}

// IsFireEvent 是否触发
func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && field.Content == ".mc"
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{
		Elements: make([]message.IMessageElement, 1),
	}
	b, err := randomImage()
	if err != nil {
		return nil, err
	}
	var image message.IMessageElement
	if plugins.GroupMessage == request.MessageType {
		image, err = request.QQClient.UploadGroupImage(request.GroupCode, bytes.NewReader(*b))
	} else {
		image, err = request.QQClient.UploadPrivateImage(request.Sender.Uin, bytes.NewReader(*b))
	}
	if err != nil {
		return nil, err
	}
	result.Elements[0] = image
	return result, nil
}

func init() {
	exists, _ := pathExists("./mc")
	if !exists {
		os.Mkdir("./mc", 0777)
	}
	plugins.RegisterOnMessagePlugin(Plugin{})
}

type Resp struct {
	ImgUrl string `json:"imgurl"`
	Code   string `json:"code"`
}

var randomUrl = "https://api.ixiaowai.cn/mcapi/mcapi.php?return=json"

var client = http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func randomImage() (*[]byte, error) {
	r, err := client.Get(randomUrl)
	if err != nil {
		return nil, err
	}
	robots, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var resp Resp
	err = json.Unmarshal(robots, &resp)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	if resp.Code != "200" {
		return nil, errors.New("请求接口失败")
	}
	b, e := getFile(resp.ImgUrl)
	return &b, e
}

func getFile(url string) ([]byte, error) {
	path := getFileName(url)

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
	r, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	_ = ioutil.WriteFile(path, content, 0644)
	return content, err
}

func getFileName(url string) string {
	i := strings.LastIndex(url, "/")
	return fmt.Sprintf("./mc/%v", url[i+1:])
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
