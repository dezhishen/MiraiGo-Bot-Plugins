package pixiv

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
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
	plugins.RegisterSchedulerPlugin(p)
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
			res, err := random()
			if err == nil && res != nil {
				elements = append(elements, message.NewText(fmt.Sprintf("标题:%v\n作者:%v\n原地址:https://www.pixiv.net/artworks/%v\n", res.Title, res.UserName, res.IllustID)))
				for _, url := range res.Urls {
					r, _ := http.DefaultClient.Get(url)
					robots, _ := ioutil.ReadAll(r.Body)
					r.Body.Close()
					var image message.IMessageElement
					if plugins.GroupMessage == request.MessageType {
						image, err = request.QQClient.UploadGroupImage(request.GroupCode, bytes.NewReader(robots))
						if err != nil {
							elements = append(elements, message.NewText("\n"+url+"\n"))
						} else {
							elements = append(elements, image)
						}
					} else {
						image, err = request.QQClient.UploadPrivateImage(request.Sender.Uin, bytes.NewReader(robots))
						if err != nil {
							elements = append(elements, message.NewText("\n"+url+"\n"))
						} else {
							elements = append(elements, image)
						}
					}
				}
			}
			if plugins.GroupMessage == request.MessageType && len(params) > 2 {
				bucket := []byte(w.PluginInfo().ID)
				key := []byte(fmt.Sprintf("pixiv.enable.%v", request.GroupCode))
				if params[2] == "Y" {
					storage.Put(bucket, key, storage.IntToBytes(1))
					elements = append(elements, message.NewText("\n已开启定时发送"))
				} else if params[2] == "N" {
					storage.Delete(bucket, key)
					elements = append(elements, message.NewText("\n已关闭定时发送"))
				}
			}
		case "r18":
			var cacheKey string
			if plugins.GroupMessage == request.MessageType {
				cacheKey = fmt.Sprintf("%v%v", request.MessageType, request.GroupCode)
			} else {
				cacheKey = fmt.Sprintf("%v%v", request.MessageType, request.Sender.Uin)
			}
			res, err := randomR18(string(cacheKey))
			if err == nil && res != nil {
				elements = append(elements, message.NewText(fmt.Sprintf("标题:%v\n作者:%v\n原地址:https://www.pixiv.net/artworks/%v\n", res.Title, res.UserName, res.IllustID)))
				for _, url := range res.Urls {
					r, _ := http.DefaultClient.Get(url)
					robots, _ := ioutil.ReadAll(r.Body)
					r.Body.Close()
					var image message.IMessageElement
					if plugins.GroupMessage == request.MessageType {
						elements = append(elements, message.NewText("\n"+url+"\n"))
					} else {
						image, err = request.QQClient.UploadPrivateImage(request.Sender.Uin, bytes.NewReader(robots))
						if err != nil {
							elements = append(elements, message.NewText("\n"+url+"\n"))
						} else {
							elements = append(elements, image)
						}
					}
				}
			}
		default:
			elements = append(elements, message.NewText("错误的命令"))
		}
	}
	result.Elements = elements
	return result, nil
}

// Cron cron表达式
func (p Plugin) Cron() string {
	return "0 */5 * * * ?"
}

// Run 回调
func (p Plugin) Run(bot *bot.Bot) error {
	bucket := []byte(p.PluginInfo().ID)
	groups, err := bot.GetGroupList()
	if err != nil {
		fmt.Printf("pixiv r send msg err %v", err)
	}
	for _, g := range groups {
		key := []byte(fmt.Sprintf("pixiv.enable.%v", g.Code))
		value, _ := storage.GetValue(bucket, key)
		if value != nil && storage.BytesToInt(value) == 1 {
			sendingMessage := &message.SendingMessage{}
			var elements []message.IMessageElement
			res, err := random()
			if err != nil {
				continue
			}
			elements = append(elements, message.NewText(fmt.Sprintf("标题:%v\n作者:%v\n原地址:https://www.pixiv.net/artworks/%v\n", res.Title, res.UserName, res.IllustID)))
			for _, url := range res.Urls {
				r, _ := http.DefaultClient.Get(url)
				robots, _ := ioutil.ReadAll(r.Body)
				defer r.Body.Close()
				var image message.IMessageElement
				image, err = bot.UploadGroupImage(g.Code, bytes.NewReader(robots))
				if err != nil {
					continue
				}
				elements = append(elements, image)
			}
			elements = append(elements, message.NewText("\n来自定时发送,可以发送[.pixiv r N]关闭"))
			sendingMessage.Elements = elements
			go bot.QQClient.SendGroupMessage(g.Code, sendingMessage)
		}
	}
	return nil
}
