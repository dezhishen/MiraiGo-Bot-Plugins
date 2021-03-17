package pixiv

import (
	"bytes"
	"errors"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

// Plugin  Pixiv助手插件
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}

// PluginInfo PluginInfo
func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   ".pixiv",
		Name: "Pixiv助手",
	}
}

// IsFireEvent 是否触发
func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".pixiv")
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{}
	var elements []message.IMessageElement

	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	params := strings.Split(context, " ")
	if len(params) > 1 {
		command := strings.TrimSpace(params[1])
		switch command {
		case "r":
			var platform string
			if len(params) > 2 {
				if params[2] == "p" {
					platform = "pc"
				} else if params[2] == "m" {
					platform = "mobile"
				} else {
					return nil, errors.New("不支持的类型,只支持 m/pc (mobile/pc)")
				}
			}
			b, err := randomImage(platform)
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
			elements = append(elements, image)
		}
	}
	result.Elements = elements
	return result, nil
}
