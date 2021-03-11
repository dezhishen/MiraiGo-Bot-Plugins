package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

// Hitokoto 一言插件
type Hitokoto struct {
}

// HitokotoPlugin 一言插件
type HitokotoPlugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (w HitokotoPlugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "hitokoto",
		Name: "一言插件",
	}
}

// IsFireEvent 是否出发
func (w HitokotoPlugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && field.Content == ".hitokoto"
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (w HitokotoPlugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{
		Elements: make([]message.IMessageElement, 1),
	}
	r, err := http.DefaultClient.Get("https://v1.hitokoto.cn/")
	if err != nil {
		return nil, err
	}
	robots, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	resp := string(robots)
	if resp == "" {
		return nil, nil
	}
	var mapResult map[string]interface{}
	err = json.Unmarshal([]byte(resp), &mapResult)
	if err != nil {
		return nil, err
	}
	name := request.Sender.CardName
	if name == "" {
		name = request.Sender.Nickname
	}
	txt := fmt.Sprintf(
		"For %v : %v",
		name,
		mapResult["hitokoto"],
	)
	result.Elements[0] = message.NewText(txt)
	return result, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(HitokotoPlugin{})
}
