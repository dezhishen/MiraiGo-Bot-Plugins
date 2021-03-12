package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/tools"
)

// Plugin 天气插件
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "weather",
		Name: "天气插件",
	}
}

// IsFireEvent 是否触发
func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".weather ")
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{
		Elements: make([]message.IMessageElement, 1),
	}
	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	params := strings.Split(context, " ")
	if len(params) < 3 {
		return nil, errors.New("请至少输入省份和城市")
	}
	if len(params) == 3 {
		params = append(params, "")
	}
	uri := fmt.Sprintf("https://wis.qq.com/weather/common?source=pc&weather_type=observe&province=%v&city=%v&county=%v", params[1], params[2], params[3])
	r, err := http.DefaultClient.Get(uri)
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
	flatMap := tools.FlatMap("", mapResult)
	txt := fmt.Sprintf(
		"%v%v%v现在天气是%v,温度为%v",
		params[1],
		params[2],
		params[3],
		flatMap["data.observe.weather_short"],
		flatMap["data.observe.degree"],
	)
	result.Elements[0] = message.NewText(txt)
	return result, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}
