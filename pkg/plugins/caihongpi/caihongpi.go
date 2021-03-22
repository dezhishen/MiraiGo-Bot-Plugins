package caihongpi

import (
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

// Plugin Random插件
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

var pluginID = "caihongpi"

// PluginInfo PluginInfo
func (p Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   pluginID,
		Name: "彩虹屁",
	}
}

// IsFireEvent 是否触发
func (p Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".chp")
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (p Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{}
	caihongpi, err := getcaihongpi()
	if err != nil {
		return nil, err
	}
	result.Elements = append(result.Elements, message.NewText(caihongpi))
	return result, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}

var url = "https://api.muxiaoguo.cn/api/caihongpi?api_key=%v"

type resp struct {
	Data *data  `json:"data"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type data struct {
	Comment string `json:"comment"`
}

func getcaihongpi() (string, error) {
	key, err := getKey()
	if err != nil {
		return "", err
	}
	r, err := http.DefaultClient.Get(fmt.Sprintf(url, key))
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	robots, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	var resp resp
	err = json.Unmarshal(robots, &resp)
	if err != nil {
		return "", err
	}
	if resp.Code != 200 {
		return "", errors.New(resp.Msg)
	}
	return resp.Data.Comment, nil
}

func getKey() (string, error) {
	str := os.Getenv("BOT_CHP_KEY")
	if str == "" {
		return "", errors.New("未配置毒鸡汤的key")
	}
	return str, nil
}
