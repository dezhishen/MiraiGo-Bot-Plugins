package todayFortune

import (
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot-Plugins/pkg/fortune"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

// Plugin Random插件
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (p Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "todayFortune",
		Name: "今日运势",
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
		return strings.HasPrefix(field.Content, ".tf") || field.Content == "运势" || field.Content == "签到"
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (p Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result, err := fortune.Randtext()
	if err != nil {
		return nil, err
	}
	print(result.Title, result.Content)
	return nil, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}
