package tips

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
	"github.com/go-basic/uuid"
)

// Info 消息详情
type Info struct {
	ID      string `json:"ID"`
	Content string `json:"content"`
	// 发送类型 1 私聊 2群聊
	SendType  int   `json:"sendType"`
	GroupCode int64 `json:"groupCode"`
	SenderUID int64 `json:"senderUID"`
	Hour      int   `json:"hour"`
	Minute    int   `json:"minute"`
	EveryDay  bool  `json:"everyDay"`
}

// Tips 提示
type Tips struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (t Tips) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "tips",
		Name: "提醒插件",
	}
}

// IsFireEvent 是否触发
func (t Tips) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".tips ")
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (t Tips) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{
		Elements: make([]message.IMessageElement, 1),
	}
	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	params := strings.Split(context, " ")
	if params[1] == "help" {
		result.Elements[0] = message.NewText(fmt.Sprintf("请输入 .tips 小时:分钟 提示内容 每天重复(Y/N[默认])\n .tips remove 小时:分钟 将会删除由你创建的该时间点的提醒"))
		return result, nil
	}
	if params[1] == "remove" {
		timestr := params[2]
		hAndM := strings.Split(timestr, ":")
		if len(hAndM) != 2 {
			return nil, errors.New("错误的时间格式")
		}
		hour, err := strconv.Atoi(hAndM[0])
		if err != nil {
			return nil, errors.New("错误的时间格式")
		}
		min, err := strconv.Atoi(hAndM[1])
		if err != nil {
			return nil, errors.New("错误的时间格式")
		}
		prefix := []byte(fmt.Sprintf("tips.%v.%v", hour, min))
		var keys []([]byte)
		storage.GetByPrefix([]byte(t.PluginInfo().ID), prefix, func(k, v []byte) error {
			var info Info
			err := json.Unmarshal(v, &info)
			if err != nil {
				return err
			}
			if info.SenderUID == request.Sender.Uin {
				keys = append(keys, k)
			}
			return nil
		})
		for _, k := range keys {
			storage.Delete([]byte(t.PluginInfo().ID), k)
		}
		result.Elements[0] = message.NewText(fmt.Sprintf("已经移除由你创建的%v的提醒", timestr))
		return result, nil

	}
	if len(params) < 3 {
		return nil, errors.New("请输入 .tips 小时:分钟 提示内容 每天重复(Y/N[默认])")
	}
	timestr := params[1]
	content := params[2]
	sendType := 2
	if len(params) > 3 {
		if params[3] == "1" {
			sendType = 1
		} else {
			sendType = 2
		}
	}
	everyDay := len(params) > 4 && params[4] == "Y"
	hAndM := strings.Split(timestr, ":")
	if len(hAndM) != 2 {
		return nil, errors.New("错误的时间格式")
	}
	hour, err := strconv.Atoi(hAndM[0])
	if err != nil {
		return nil, errors.New("错误的时间格式")
	}
	min, err := strconv.Atoi(hAndM[1])
	if err != nil {
		return nil, errors.New("错误的时间格式")
	}
	if err != nil {
		return nil, err
	}
	info := Info{
		ID:        uuid.New(),
		Content:   content,
		SendType:  sendType,
		SenderUID: request.Sender.Uin,
		Hour:      hour,
		Minute:    min,
		EveryDay:  everyDay,
		GroupCode: request.GroupCode,
	}
	jsonBytes, _ := json.Marshal(info)
	err = storage.Put([]byte(t.PluginInfo().ID), []byte(fmt.Sprintf("tips.%v.%v.%v", info.Hour, info.Minute, info.ID)), jsonBytes)
	if err != nil {
		return nil, err
	}
	result.Elements[0] = message.NewText(fmt.Sprintf("For %v:提醒创建成功!", request.GetNickName()))
	return result, nil
}

// Run 回调
func (t Tips) Run(bot *bot.Bot) error {
	nowLocal := time.Now().Local()
	prefix := []byte(fmt.Sprintf("tips.%v.%v.", nowLocal.Hour(), nowLocal.Minute()))
	storage.GetByPrefix([]byte(t.PluginInfo().ID), prefix, func(key, v []byte) error {
		var info Info
		err := json.Unmarshal(v, &info)
		if err != nil {
			return err
		}
		sendingMessage := &message.SendingMessage{}
		sendingMessage.Append(message.NewAt(info.SenderUID))
		sendingMessage.Append(message.NewText(info.Content))
		go bot.QQClient.SendGroupMessage(info.GroupCode, sendingMessage)
		if !info.EveryDay {
			go storage.Delete([]byte(t.PluginInfo().ID), key)
		}
		return nil
	})
	// bot.QQClient.Send
	return nil
}

// Cron cron表达式
func (t Tips) Cron() string {
	return "0 */1 * * * ?"
}
func init() {
	tips := Tips{}
	plugins.RegisterOnMessagePlugin(tips)
	plugins.RegisterSchedulerPlugin(tips)
}
