package weibolisten

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
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
func (p Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
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
		key := []byte(fmt.Sprintf("weibo-listen.user.%v", params[2]))
		err := storage.Get([]byte(p.PluginInfo().ID), key, func(s []byte) error {
			var listenUser ListenUser
			if s == nil {
				listenUser = ListenUser{
					UID: params[2],
				}
				err := setContainerId(&listenUser)
				if err != nil {
					return err
				}
				jsonBytes, _ := json.Marshal(listenUser)
				storage.Put([]byte(p.PluginInfo().ID), key, jsonBytes)
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
		messageKey := []byte(
			fmt.Sprintf(
				"weibo-listen.user-sender.%v.%v.%v",
				listenUserMessage.UID,
				listenUserMessage.ReciveType,
				listenUserMessage.Reciver,
			))
		notExists := false
		storage.Get([]byte(p.PluginInfo().ID), messageKey, func(s []byte) error {
			notExists = s == nil
			return nil
		})
		if notExists {
			jsonBytes, _ := json.Marshal(listenUserMessage)
			err = storage.Put([]byte(p.PluginInfo().ID), messageKey, jsonBytes)
			if err != nil {
				return nil, err
			}
			incrKey := []byte(fmt.Sprintf("weibo-listen.user-count.%v", params[2]))
			_, err = storage.Incr([]byte(p.PluginInfo().ID), incrKey, 1)
		}
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
		messageKey := []byte(fmt.Sprintf("weibo-listen.user-sender.%v.%v.%v", listenUserMessage.UID, listenUserMessage.ReciveType, listenUserMessage.Reciver))
		err := storage.Delete([]byte(p.PluginInfo().ID), messageKey)
		if err != nil {
			return nil, err
		}
		incrKey := []byte(fmt.Sprintf("weibo-listen.user-count.%v", params[2]))
		doIncr := false
		storage.Get([]byte(p.PluginInfo().ID), incrKey, func(s []byte) error {
			if s == nil {
				return nil
			}
			count := storage.BytesToInt(s)
			doIncr = count > 0
			return nil
		})
		if doIncr {
			_, err = storage.Incr([]byte(p.PluginInfo().ID), incrKey, -1)
		}
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
	prefix := []byte("weibo-listen.user.")
	var listenUsers []ListenUser
	storage.GetByPrefix([]byte(p.PluginInfo().ID), prefix, func(key, v []byte) error {
		var info ListenUser
		err := json.Unmarshal(v, &info)
		if err != nil {
			return err
		}
		countKey := []byte(fmt.Sprintf("weibo-listen.user-count.%v", info.UID))
		var count int
		err = storage.Get([]byte(p.PluginInfo().ID), countKey, func(s []byte) error {
			if s == nil {
				count = 0
			} else {
				count = storage.BytesToInt(s)
				return nil
			}
			return nil
		})
		if err != nil {
			return err
		}
		if count == 0 {
			return nil
		}
		listenUsers = append(listenUsers, info)
		return nil
	})
	for _, info := range listenUsers {
		key := []byte(fmt.Sprintf("weibo-listen.user.%v", info.UID))
		if info.ContainerID == "" {
			setContainerId(&info)
			jsonBytes, _ := json.Marshal(info)
			storage.Put([]byte(p.PluginInfo().ID), key, jsonBytes)
		}
		card, err := getLastItemsByUIDAndContainerID(info.UID, info.ContainerID)
		if err != nil {
			fmt.Printf("%v", err)
			return err
		}
		if card.Mblog.BID == info.LastWeiboID {
			return nil
		}
		info.LastWeiboID = card.Mblog.BID
		messageKey := []byte(fmt.Sprintf("weibo-listen.user-sender.%v.", info.UID))
		m := message.SendingMessage{}
		m.Elements = append(m.Elements, message.NewText(fmt.Sprintf("%v 发了微博 %v", card.Mblog.User.ScreenName, card.Mblog.Text)))
		m.Elements = append(m.Elements, message.NewText(card.Scheme))
		err = storage.GetByPrefix([]byte(p.PluginInfo().ID), messageKey, func(k, senderValue []byte) error {
			var sender ListenUserMessage
			json.Unmarshal([]byte(senderValue), &sender)
			if &sender == nil {
				return nil
			}
			if sender.ReciveType == GroupMessage {
				bot.QQClient.SendGroupMessage(sender.Reciver, &m)
			} else if sender.ReciveType == PrivateMessage {
				bot.QQClient.SendPrivateMessage(sender.Reciver, &m)
			}
			return nil
		})
		if err != nil {
			return err
		}
		jsonBytes, _ := json.Marshal(info)
		storage.Put([]byte(p.PluginInfo().ID), key, jsonBytes)
	}
	// bot.QQClient.Send
	return nil
}

// Cron cron表达式
func (t Plugin) Cron() string {
	// return "0 0/5 * * * ?"
	return "0 */1 * * * ?"
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
		id, err := getContainerIDByUid(info.UID)
		if err != nil {
			return err
		}
		if id == "" {
			return errors.New("UID错误,或者该微博不可见")
		}
		info.ContainerID = id
	}
	return nil
}

func getLastItemsByUIDAndContainerID(UID, ContainerID string) (*WeiboContentCard, error) {
	var contentResp WeiboContentResp
	url := fmt.Sprintf(
		"https://m.weibo.cn/api/container/getIndex?uid=%v&t=0&type=uid&containerid=%v",
		UID,
		ContainerID,
	)
	resp, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	respBodyStr := string(robots)
	if respBodyStr == "" {
		return nil, err
	}
	err = json.Unmarshal(robots, &contentResp)
	if err != nil {
		return nil, err
	}
	if &contentResp == nil {
		return nil, nil
	}
	for _, c := range contentResp.Data.Cards {
		if c.Mblog.IsTop == 1 {
			continue
		}
		return &c, nil
	}
	return nil, nil
}

type WeiboContentResp struct {
	Data *WeiboContentData `json:"data"`
}

type WeiboContentData struct {
	Cards []WeiboContentCard `json:"cards"`
}

type WeiboContentCard struct {
	CardType int                `json:"card_type"`
	ItemID   string             `json:"itemid"`
	Scheme   string             `json:"scheme"`
	Mblog    *WeiboContentMblog `json:"mblog"`
}

type WeiboContentMblog struct {
	IsTop int               `json:"isTop"`
	Text  string            `json:"text"`
	BID   string            `json:"bid"`
	Pics  []WeiboContentPic `json:"pics"`
	User  *WeiboUser        `json:"user"`
}

type WeiboUser struct {
	ScreenName string `json:"screen_name"`
}
type WeiboContentPic struct {
	PID string `json:"pid"`
	Url string `json:"url"`
}

type ContainerResp struct {
	Data ContainerData `json:"data"`
}

type ContainerData struct {
	TabsInfo TabsInfo `json:"tabsInfo"`
}

type Tab struct {
	ID          int    `json:"id"`
	TabKey      string `json:"tabKey"`
	TabType     string `json:"tab_type"`
	ContainerID string `json:"containerid"`
}

type TabsInfo struct {
	Tabs []Tab `json:"tabs"`
}

type User struct {
}
