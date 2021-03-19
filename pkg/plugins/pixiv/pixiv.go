package pixiv

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

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
			platform := "mobile"
			var loop int
			var err error
			size := "large"
			messageType := "image"
			for i := 2; i < len(params); i++ {
				if params[i] == "-h" || params[i] == "--help" {
					elements = append(elements, message.NewText(
						".pixiv r "+
							"\n-p,--pc/-m,--mobile 指定pc格式还是mobile格式 "+
							"\n-original/-large/-medium/-squareMedium 指定尺寸 "+
							"\n-n$num 指定数量,超过10则为10"+
							"\n-t,--text 指定返回地址而非图片"))
					result.Elements = elements
					return result, nil
				}
				if params[i] == "-p" || params[i] == "--pc" {
					platform = "pc"
					continue
				}
				if params[i] == "-m" || params[i] == "--mobile" {
					platform = "mobile"
					continue
				}
				if strings.HasPrefix(params[i], "-n") {
					loop, err = strconv.Atoi(strings.TrimPrefix(params[i], "-n"))
					if err != nil {
						return nil, err
					}
					continue
				}
				if params[i] == "-original" {
					size = "original"
					continue
				}
				if params[i] == "-large" {
					size = "large"
					continue
				}
				if params[i] == "-medium" {
					size = "medium"
					continue
				}
				if params[i] == "-squareMedium" {
					size = "squareMedium"
					continue
				}
				if params[i] == "-t" || params[i] == "--text" {
					messageType = "text"
				}
			}
			if messageType == "text" {
				size = "original"
			}
			if loop < 1 {
				loop = 1
			} else if loop > 10 {
				loop = 10
			}
			for i := 0; i < loop; i++ {
				b, err := randomImage(platform, size, messageType)
				if err != nil {
					return nil, err
				}
				if b == nil {
					continue
				}
				if messageType == "image" {
					var image message.IMessageElement
					if plugins.GroupMessage == request.MessageType {
						image, err = request.QQClient.UploadGroupImage(request.GroupCode, bytes.NewReader(*b))
					} else {
						image, err = request.QQClient.UploadPrivateImage(request.Sender.Uin, bytes.NewReader(*b))
					}
					if image == nil {
						i--
						continue
					}
					if err != nil {
						out, err := os.Create(fmt.Sprintf("%v.jpg", time.Now().Unix()))
						if err == nil {
							io.Copy(out, bytes.NewReader(*b))
							out.Close()
						}
						continue
					}
					elements = append(elements, image)
				} else {
					elements = append(elements, message.NewText(string(*b)+"\n"))
				}
			}
		case "r18":
			res, err := randomR18()
			if err == nil && res != nil {
				elements = append(elements, message.NewText(fmt.Sprintf("标题:%v\n作者:%v\n原地址:https://www.pixiv.net/artworks/%v\n", res.Title, res.UserName, res.IllustID)))
				for _, url := range res.Urls {
					// r, _ := http.DefaultClient.Get(url)
					// robots, _ := ioutil.ReadAll(r.Body)
					// r.Body.Close()
					// var image message.IMessageElement
					// if plugins.GroupMessage == request.MessageType {
					// 	image, err = request.QQClient.UploadGroupImage(request.GroupCode, bytes.NewReader(robots))
					// } else {
					// 	image, err = request.QQClient.UploadPrivateImage(request.Sender.Uin, bytes.NewReader(robots))
					// }
					// if err != nil {
					// 	log.Print(err)
					// }
					// elements = append(elements, image)
					elements = append(elements, message.NewText("\n"+url+"\n"))
				}
			}
		default:
			elements = append(elements, message.NewText("错误的命令"))
		}
	}
	result.Elements = elements
	return result, nil
}
