package jrrp

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
)

// Plugin Random插件
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (p Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "Jrrp",
		Name: "今日人品插件",
	}
}

// IsFireEvent 是否触发
func (p Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".jrrp")
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (p Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{
		Elements: make([]message.IMessageElement, 1),
	}
	timeNow := time.Now().Local()
	timestr := timeNow.Format("2020-02-08")
	key := []byte(fmt.Sprintf("jrrp.%v.%v", request.Sender.Uin, timestr))
	var score int
	err := storage.Get([]byte(p.PluginInfo().ID), key, func(b []byte) error {
		if b != nil {
			score = storage.BytesToInt(b)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if score == 0 {
		rand.Seed(time.Now().UnixNano())
		score = rand.Intn(100) + 1
		storage.Put([]byte(p.PluginInfo().ID), key, storage.IntToBytes(score))
		keyLast7Day := timeNow.AddDate(0, 0, -7).Format("2020-02-08")
		storage.Delete([]byte(p.PluginInfo().ID), []byte(keyLast7Day))
	}
	result.Elements[0] = message.NewText(fmt.Sprintf("[%v]今日人品: %v", request.GetNickName(), score))
	return result, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}
