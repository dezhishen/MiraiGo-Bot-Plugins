package weibolisten

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
)

// ListenUser 监听的用户
type ListenUser struct {
	//用户UID
	UID string `json:"uid"`
	//用户ContainerID
	ContainerID string `json:"containerId"`
	//用户最后微博ID
	LastWeiboID string `json:"lastWeiboId"`
}
type messageType string

const (
	//GroupMessage 群消息
	GroupMessage = messageType("group")
	//PrivateMessage 私聊消息
	PrivateMessage = messageType("private")
)

// ListenUserMessage 监听用户需要发送的群或者其他
type ListenUserMessage struct {
	// 微博用户ID
	UID string `json:"id"`
	// 1 私聊/2 群聊
	ReciveType messageType `json:"type"`
	// QQ的 私聊 用户ID / 群聊用户ID
	Reciver int64 `json:"senderId"`
}

// Plugin 微博监听
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:          "weibo-listen",
		Name:        "微博监听插件",
		Description: "涩图感应器(√)",
	}
}

// IsFireEvent 是否触发
func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".weibo-l ")
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	var elements []message.IMessageElement

	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	params := strings.Split(context, " ")
	command := params[1]
	switch command {
	case "add":
		if len(params) < 3 {
			return nil, errors.New("请输入要添加的微博账户的UID")
		}
		key := fmt.Sprintf("weibo-listen.user.%v", params[2])
		err := storage.Get(w.PluginInfo().ID, key, func(s string) error {
			var listenUser ListenUser
			if s == "" {
				listenUser = ListenUser{
					UID: params[2],
				}
				err := setContainerId(&listenUser)
				if err != nil {
					return err
				}
				jsonBytes, _ := json.Marshal(listenUser)
				storage.Put(w.PluginInfo().ID, key, string(jsonBytes))
			} else {
				_ = json.Unmarshal([]byte(s), &listenUser)
			}
			return nil
		})
		if err != nil {
			return nil, err
		}
		listenUserMessage := ListenUserMessage{
			UID:        params[2],
			ReciveType: messageType(request.MessageType),
		}
		if request.MessageType == plugins.GroupMessage {
			listenUserMessage.Reciver = request.GroupCode
		} else {
			listenUserMessage.Reciver = request.Sender.Uin
		}
		messageKey := fmt.Sprintf("weibo-listen.user-sender.%v.%v.%v", listenUserMessage.UID, listenUserMessage.ReciveType, listenUserMessage.Reciver)
		jsonBytes, _ := json.Marshal(listenUserMessage)
		err = storage.Put(w.PluginInfo().ID, messageKey, string(jsonBytes))
		if err != nil {
			return nil, err
		}
		incrKey := fmt.Sprintf("weibo-listen.user-count.%v", params[2])
		_, err = storage.Incr(w.PluginInfo().ID, incrKey, 1)
		if err != nil {
			return nil, err
		}
		elements = append(elements, message.NewText("添加成功!"))
	case "remove":
		if len(params) < 3 {
			return nil, errors.New("请输入要移除的微博账户的UID")
		}
		listenUserMessage := ListenUserMessage{
			UID:        params[2],
			ReciveType: messageType(request.MessageType),
		}
		if request.MessageType == plugins.GroupMessage {
			listenUserMessage.Reciver = request.GroupCode
		} else {
			listenUserMessage.Reciver = request.Sender.Uin
		}
		messageKey := fmt.Sprintf("weibo-listen.user-sender.%v.%v.%v", listenUserMessage.UID, listenUserMessage.ReciveType, listenUserMessage.Reciver)
		err := storage.Delete(w.PluginInfo().ID, messageKey)
		if err != nil {
			return nil, err
		}
		incrKey := fmt.Sprintf("weibo-listen.user-count.%v", params[2])
		_, err = storage.Incr(w.PluginInfo().ID, incrKey, -1)
		if err != nil {
			return nil, err
		}
		elements = append(elements, message.NewText("移除成功!"))
	case "help":
		elements = append(elements, message.NewText("add uid 增加一个监听的微博用户\nremove uid 移除一个监听的微博用户"))
	default:
		elements = append(elements, message.NewText("只支持 add remove help命令"))
	}
	result := &plugins.MessageResponse{
		Elements: elements,
	}
	return result, nil
}

// Run 回调
func (p Plugin) Run(bot *bot.Bot) error {
	prefix := "weibo-listen.user."
	storage.GetByPrefix(p.PluginInfo().ID, prefix, func(key, v string) error {
		var info ListenUser
		err := json.Unmarshal([]byte(v), &info)
		if err != nil {
			return err
		}
		countKey := fmt.Sprintf("weibo-listen.user-count.%v", info.UID)
		var count int
		err = storage.Get(p.PluginInfo().ID, countKey, func(s string) error {
			if s == "" {
				count = 0
			} else {
				count, err = strconv.Atoi(s)
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}
		if count == 0 {
			return nil
		}
		// sendingMessage := &message.SendingMessage{}
		return nil
	})
	// bot.QQClient.Send
	return nil
}

// Cron cron表达式
func (t Plugin) Cron() string {
	return "0 0/5 * * * ?"
}

func init() {
	p := Plugin{}
	plugins.RegisterOnMessagePlugin(p)
	plugins.RegisterSchedulerPlugin(p)
}

func getContainerIDByUid(uid string) (string, error) {
	url := fmt.Sprintf("https://m.weibo.cn/api/container/getIndex?type=uid&value=%v", uid)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", err
	}
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", err
	}
	respBodyStr := string(robots)
	if respBodyStr == "" {
		return "", err
	}
	var containerResp ContainerResp
	err = json.Unmarshal(robots, &containerResp)
	if err != nil {
		return "", err
	}
	for _, v := range containerResp.Data.TabsInfo.Tabs {
		if v.TabKey == "weibo" {
			return v.ContainerID, nil

		}
	}
	return "", nil
}

func setContainerId(info *ListenUser) error {
	if info.ContainerID == "" {
		id, err := getContainerIDByUid(info.ContainerID)
		if err != nil {
			return err
		}
		info.ContainerID = id
	}
	return nil
}

type ContainerResp struct {
	Data ContainerData `json:"data"`
}

type ContainerData struct {
	TabsInfo TabsInfo `json:"tabsInfo"`
}

type Tab struct {
	ID          string `json:"id"`
	TabKey      string `json:"tabKey"`
	TabType     string `json:"tab_type"`
	ContainerID string `json:"containerid"`
}

type TabsInfo struct {
	SelectedTab string `json:"selectedTab"`
	Tabs        []Tab
}

type User struct {
}
