package vader

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/jonreiter/govader"
)

var analyzer = govader.NewSentimentIntensityAnalyzer()

// Plugin vader
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (p Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "vader",
		Name: "情感分析插件",
	}
}

// IsFireEvent 是否触发
func (p Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".vader")
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (p Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{}
	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	params := strings.TrimSpace(strings.TrimPrefix(context, ".vader"))
	if params == "" {
		return nil, errors.New("请输入要分析的话")
	}
	var elements []message.IMessageElement
	text := strings.TrimSpace(params)
	sentiment := analyzer.PolarityScores(text)
	elements = append(elements, message.NewText(text))
	elements = append(elements, message.NewText(fmt.Sprintf("\n综合分数: %v", sentiment.Compound)))
	elements = append(elements, message.NewText(fmt.Sprintf("\n积极得分: %v", sentiment.Positive)))
	elements = append(elements, message.NewText(fmt.Sprintf("\n中立得分: %v", sentiment.Neutral)))
	elements = append(elements, message.NewText(fmt.Sprintf("\n消极得分: %v", sentiment.Negative)))
	result.Elements = elements
	return result, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}
