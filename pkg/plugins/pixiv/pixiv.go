package pixiv

import (
	"bytes"
	"fmt"
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
	p := Plugin{}
	plugins.RegisterOnMessagePlugin(p)
	// plugins.RegisterSchedulerPlugin(p)
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
	res, err := getSetu()
	if err != nil {
		return nil, err
	}
	elements = append(elements, message.NewText(fmt.Sprintf("标题:%v\n作者:%v\n原地址:https://www.pixiv.net/artworks/%v\n", res.Title, res.UserName, res.IllustID)))
	if plugins.GroupMessage == request.MessageType {
		for _, image := range *res.Images {
			imageElement, err := request.QQClient.UploadGroupImage(request.GroupCode, bytes.NewReader(*image))
			if err != nil {
				return nil, err
			}
			elements = append(elements, imageElement)
		}

	} else {
		for _, image := range *res.Images {
			imageElement, err := request.QQClient.UploadPrivateImage(request.Sender.Uin, bytes.NewReader(*image))
			if err != nil {
				return nil, err
			}
			elements = append(elements, imageElement)
		}
	}
	result.Elements = elements
	return result, nil
}
