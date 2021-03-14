package jrrp

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"
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
	result := &plugins.MessageResponse{}
	timeNow := time.Now().Local()
	score, err := getScore(timeNow, p.PluginInfo().ID, request.Sender.Uin, true)
	if err != nil {
		return nil, err
	}
	var elements []message.IMessageElement
	elements = append(elements, message.NewText(fmt.Sprintf("[%v]今日人品: %v", request.GetNickName(), score)))

	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	params := strings.Split(context, " ")

	if len(params) > 1 {
		preDays, err := strconv.Atoi(params[1])
		if err != nil {
			return nil, errors.New("请输入一个7以内的正整数")
		}
		if preDays > 7 {
			preDays = 7
		}
		for i := 1; i <= preDays; i++ {
			timeNow = timeNow.Add(-1 * 24 * time.Hour)
			score, err := getScore(timeNow, p.PluginInfo().ID, request.Sender.Uin, true)
			if err != nil {
				return nil, err
			}
			if score == 0 {
				break
			}
			elements = append(elements, message.NewText(fmt.Sprintf("\n[%v]历史人品: %v", request.GetNickName(), score)))
		}
	}
	result.Elements = elements
	return result, nil
}

func getScore(t time.Time, pid string, uid int64, genIfNil bool) (int, error) {
	timestr := fmt.Sprintf("%v-%v-%v", t.Year(), t.Month(), t.Day())
	key := []byte(fmt.Sprintf("jrrp.%v.%v", uid, timestr))
	var score int
	err := storage.Get([]byte(pid), key, func(b []byte) error {
		if b != nil {
			score = storage.BytesToInt(b)
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	if genIfNil && score == 0 {
		rand.Seed(time.Now().UnixNano())
		score = rand.Intn(100) + 1
		storage.Put([]byte(pid), key, storage.IntToBytes(score))
		keyLast7Day := fmt.Sprintf("jrrp.%v.%v", uid, t.AddDate(0, 0, -7).Format("2020-02-08"))
		storage.Delete([]byte(pid), []byte(keyLast7Day))
	}
	return score, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}
