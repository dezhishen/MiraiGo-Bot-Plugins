package segmentation

import (
	"errors"
	"fmt"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/go-ego/gse"
)

// Plugin segmentation
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

var (
	seg gse.Segmenter
)

// PluginInit 简单插件初始化
func (p Plugin) PluginInit() {
	// 加载默认字典
	seg.LoadDict()
}

// PluginInfo PluginInfo
func (p Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "segmentation",
		Name: "分词插件",
	}
}

// IsFireEvent 是否触发
func (p Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".segment")
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (p Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{}
	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	params := strings.TrimSpace(strings.TrimPrefix(context, ".segment"))
	if params == "" {
		return nil, errors.New("请输入要分析的文本")
	}
	var elements []message.IMessageElement
	text := strings.TrimSpace(params)
	sentiment := seg.Slice(text)
	elements = append(elements, message.NewText(text))
	elements = append(elements, message.NewText(fmt.Sprintf("\n结果: %v", sentiment)))
	result.Elements = elements
	return result, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}
