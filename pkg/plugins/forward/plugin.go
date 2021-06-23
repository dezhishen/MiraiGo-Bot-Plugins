package forward

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/command"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

// Plugin 消息转发插件
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

var pluginID = "forward"

// PluginInfo PluginInfo
func (p Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   pluginID,
		Name: "消息转发插件",
	}
}

// IsFireEvent 是否触发
func (p Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	// 验证消息来源
	if !contains(accouts, msg.Sender.Uin) {
		return false
	}
	if msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".forward")
	}
	return false
}

type forwardCommand struct {
	To    int64 `short:"t" long:"to" description:"转发对象" required:"true"`
	Group bool  `short:"g" long:"group" description:"是否转发到群组"`
}

func (p Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	fc := forwardCommand{}
	commands, err := command.Parse(".forward", &fc, strings.Split(context, " "))
	if err != nil {
		return nil, err
	}
	var q string
	for i := 1; i < len(commands); i++ {
		q = q + " " + commands[i]
	}
	m := message.NewSendingMessage()

	if plugins.GroupMessage == request.MessageType {
		m.Append(message.NewText(fmt.Sprintf("来自群[%v(%v)]的转发消息:\n", request.GroupName, request.GroupCode)))
	} else {
		m.Append(message.NewText(fmt.Sprintf("来自私聊[%v(%v)]的转发消息:\n", request.GetNickName(), request.Sender.Uin)))
	}
	if q != "" {
		m.Append(message.NewText(strings.TrimSpace(q)))
	}
	for i := 1; i < len(request.Elements); i++ {
		if request.Elements[i].Type() == message.Image {
			field, _ := request.Elements[i].(*message.ImageElement)
			b, _ := getImage(field.Url)
			var image message.IMessageElement
			if fc.Group {
				image, err = request.QQClient.UploadGroupImage(request.GroupCode, bytes.NewReader(b))
				if err != nil {
					return nil, err
				}
				// request.QQClient.SendGroupMessage(fc.To, m)
			} else {
				image, err = request.QQClient.UploadPrivateImage(request.Sender.Uin, bytes.NewReader(b))
				if err != nil {
					return nil, err
				}
				// request.QQClient.SendPrivateMessage(fc.To, m)
			}
			m.Append(image)
		} else {
			m.Append(request.Elements[i])
		}
	}
	if fc.Group {
		request.QQClient.SendGroupMessage(fc.To, m)
	} else {
		request.QQClient.SendPrivateMessage(fc.To, m)
	}
	return nil, nil
}

var accouts []int64

func contains(s []int64, e int64) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
	accouts = getAdminUid()
}

func getImage(url string) ([]byte, error) {
	r, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	robots, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	r.Body.Close()
	return robots, nil
}

func getAdminUid() []int64 {
	str := os.Getenv("BOT_FORWARD_ADMIN")
	if str == "" {
		return nil
	}
	accouts := strings.Split(str, ",")
	if len(accouts) == 0 {
		return nil
	}
	var result []int64
	for _, a := range accouts {
		e, err := strconv.ParseInt(a, 10, 64)
		if err != nil {
			continue
		}
		result = append(result, e)
	}
	return result
}
