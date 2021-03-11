package random

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

// Plugin Random插件
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "Random",
		Name: "Random插件",
	}
}

// IsFireEvent 是否出发
func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && field.Content == ".r"
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{
		Elements: make([]message.IMessageElement, 1),
	}
	name := request.Sender.CardName
	if name == "" {
		name = request.Sender.Nickname
	}
	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(100)
	result.Elements[0] = message.NewText(fmt.Sprintf("[%v]掷出了: %v", name, v))
	return result, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}
