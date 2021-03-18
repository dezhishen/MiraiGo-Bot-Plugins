package pixiv

import (
	"bytes"
	"errors"
	"strconv"
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
			var loop int
			if len(params) == 2 {
				platform = "mobile"
			} else if len(params) == 3 {
				loop, _ = strconv.Atoi(params[2])
				if loop == 0 {
					if params[2] == "p" {
						platform = "pc"
					} else if params[2] == "m" {
						platform = "mobile"
					} else {
						return nil, errors.New(".pixiv r m/p 数量")
					}
					loop = 1
				}
			} else if len(params) == 4 {
				if params[2] == "p" {
					platform = "pc"
				} else if params[2] == "m" {
					platform = "mobile"
				} else {
					return nil, errors.New(".pixiv r m/p 数量")
				}
				loop, _ = strconv.Atoi(params[3])
				if loop == 0 {
					loop = 1
				}
			}
			if loop > 5 {
				loop = 5
			}
			for i := 0; i < loop; i++ {
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
	}
	result.Elements = elements
	return result, nil
}
