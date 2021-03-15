package calendar

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
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
		return ok && field.Content == ".calendar"
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (p Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	msg, err := getDate(time.Now())
	if err != nil {
		return nil, err
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
	getDate(time.Now())
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
	return fmt.Sprintf(
		"%v-%v-%v,周%v %v,农历%v%v %v, %v,%v,%v",
		resp.Data.SYear,
		resp.Data.SMonth,
		resp.Data.SDay,
		resp.Data.Week,
		resp.Data.SolarFestival,
		resp.Data.IMonthChinese,
		resp.Data.IDayChinese,
		resp.Data.LunarFestival,
		resp.Data.CYear,
		resp.Data.CMonth,
		resp.Data.CDay,
	), nil
}
func init() {
	plugin := Plugin{}
	plugins.RegisterOnMessagePlugin(plugin)
	plugins.RegisterSchedulerPlugin(plugin)
}
