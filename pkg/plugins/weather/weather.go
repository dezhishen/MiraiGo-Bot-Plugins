package weather

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
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
	localtion := strings.TrimSpace(strings.ReplaceAll(context, ".weather", ""))
	resp, err := getWeather(localtion)
	if err != nil {
		return nil, err
	}
	result.Elements[0] = message.NewText(resp)
	return result, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}

func getWeather(localtion string) (string, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://wttr.in/~%v?0&AT&lang=zh", localtion), nil)
	if err != nil {
		return "", err
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer r.Body.Close()
	robots, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	resp := string(robots)
	return resp, nil
}
