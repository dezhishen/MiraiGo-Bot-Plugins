package calendar

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
)

// Plugin 日历插件
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (p Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "calendar",
		Name: "日期插件",
	}
}

// IsFireEvent 是否触发
func (p Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".calendar")
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (p Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	msg, err := getDate(time.Now())
	if err != nil {
		return nil, err
	}
	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	params := strings.Split(context, " ")
	if len(params) > 1 && request.MessageType == plugins.GroupMessage {
		bucket := []byte(p.PluginInfo().ID)
		key := []byte(fmt.Sprintf("calendar.enable.%v", request.GroupCode))
		if params[1] == "Y" {
			storage.Put(bucket, key, storage.IntToBytes(1))
			msg += "\n已启用定时发送日历"
		} else if params[1] == "N" {
			storage.Delete(bucket, key)
			msg += "\n已禁用定时发送日历"
		}
	}
	return &plugins.MessageResponse{
		Elements: []message.IMessageElement{message.NewText(msg)},
	}, nil
}

// Cron cron表达式
func (p Plugin) Cron() string {
	return "0 0 6 * * ?"
}

// Run 回调
func (p Plugin) Run(bot *bot.Bot) error {
	text, _ := getDate(time.Now())
	sendingMessage := &message.SendingMessage{}
	sendingMessage.Append(message.NewText(text))
	groups, err := bot.GetGroupList()
	if err != nil {
		fmt.Printf("calendar send msg err %v", err)
	}
	bucket := []byte(p.PluginInfo().ID)
	for _, g := range groups {
		key := []byte(fmt.Sprintf("calendar.enable.%v", g.Code))
		value, _ := storage.GetValue(bucket, key)
		if value != nil && storage.BytesToInt(value) == 1 {
			go bot.QQClient.SendGroupMessage(g.Code, sendingMessage)
		}
	}
	return nil
}

type resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Time int    `json:"time"`
	Data *data  `json:"data"`
}

type data struct {
	// 农历年
	IYear int `json:"iYear"`
	// 农历月
	IMonth int `json:"iMonth"`
	// 农历日
	IDay int `json:"iDay"`
	// 农历月 汉字
	IMonthChinese string `json:"iMonthChinese"`
	// 农历日 汉字
	IDayChinese string `json:"iDayChinese"`
	// 阳历年
	SYear int `json:"sYear"`
	// 阳历月
	SMonth int `json:"sMonth"`
	// 阳历日
	SDay int `json:"sDay"`
	// 天干地支年
	CYear string `json:"cYear"`
	// 天干地支月
	CMonth string `json:"cMonth"`
	// 天干地支日
	CDay string `json:"cDay"`
	// 是否为节假日
	IsHoliday bool `json:"isHoliday"`
	IsLeap    bool `json:"isLeap"`
	// 阳历假日
	SolarFestival string `json:"solarFestival"`
	SolarTerms    string `json:"solarTerms"`
	// 农历假日
	LunarFestival string `json:"lunarFestival"`
	// 周,汉字 一~日
	Week string `json:"week"`
}

func getDate(time time.Time) (string, error) {
	url := fmt.Sprintf("http://www.autmone.com/openapi/icalendar/queryDate?date=%v-%v-%v", time.Local().Year(), int(time.Local().Month()), time.Local().Day())
	r, err := http.DefaultClient.Get(url)
	if err != nil {
		return "", err
	}
	robots, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return "", err
	}
	if robots == nil {
		return "", nil
	}

	var resp resp
	err = json.Unmarshal(robots, &resp)
	if err != nil {
		return "", err
	}
	if resp.Code != 0 {
		return "", errors.New(resp.Msg)
	}
	if resp.Data == nil {
		return "", errors.New("data is nil")
	}
	buf := new(bytes.Buffer)
	str := strings.Join([]string{
		strconv.Itoa(resp.Data.SYear),
		"年",
		strconv.Itoa(resp.Data.SMonth),
		"月",
		strconv.Itoa(resp.Data.SDay),
		"日,周",
		resp.Data.Week,
	}, "")
	buf.WriteString(str)
	if resp.Data.SolarFestival != "" {
		buf.WriteString("\n")
		buf.WriteString(strings.TrimSpace(resp.Data.SolarFestival))
	}
	buf.WriteString("\n农历")
	buf.WriteString(resp.Data.IMonthChinese)
	buf.WriteString(resp.Data.IDayChinese)
	if resp.Data.LunarFestival != "" {
		buf.WriteString("\n")
		buf.WriteString(strings.TrimSpace(resp.Data.LunarFestival))
	}
	return buf.String(), nil
}
func init() {
	plugin := Plugin{}
	plugins.RegisterOnMessagePlugin(plugin)
	plugins.RegisterSchedulerPlugin(plugin)
}
