package lpl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"unsafe"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
)

// Plugin lpl今日明日赛事
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "lpl",
		Name: "lpl今日明日赛事查询插件",
	}
}

// IsFireEvent 是否触发
func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && field.Content == ".lpl"
	}
	return false
}

func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	var elements []message.IMessageElement

	post := handleMatchMessage(w.PluginInfo().ID)
	if unsafe.Sizeof(post) == 0 {
		elements = append(elements, message.NewText("无可用信息"))
	} else {
		elements = append(elements, message.NewText(fmt.Sprintf("%s", post.TodayPost)))
		elements = append(elements, message.NewText(" "))
		elements = append(elements, message.NewText(fmt.Sprintf("%s", post.TommorrowPost)))
	}
	result := &plugins.MessageResponse{
		Elements: elements,
	}
	return result, nil
}

// Run 回调
// func (t Plugin) Run(bot *bot.Bot) error {

// 	m := message.SendingMessage{}
// 	post := handleMatchMessage(t.PluginInfo().ID)
// 	if unsafe.Sizeof(post) == 0 {
// 		m.Elements = append(m.Elements, message.NewText("无可用信息"))
// 		return nil
// 	}

// 	m.Elements = append(m.Elements, message.NewText(fmt.Sprintf("%s", post.TodayPost)))
// 	m.Elements = append(m.Elements, message.NewText(fmt.Sprintf("%s", post.TommorrowPost)))

// 	bot.QQClient.SendGroupMessage(sender.Reciver, &m)

// 	// bot.QQClient.Send
// 	return nil
// }

func (t Plugin) Cron() string {
	return "0 */1 * * * ?"
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}

func handleMatchMessage(pluginId string) MatchPost {

	prefix := []byte("lpl.recent_match")
	var post MatchPost
	storage.GetByPrefix([]byte(pluginId), prefix, func(key, v []byte) error {
		err := json.Unmarshal(v, &post)
		if err != nil {
			return err
		}
		return nil
	})

	if unsafe.Sizeof(post) == 0 || isTimeExpired(post) {
		//查询并填充
		post = createMatchPost()
		jsonBytes, _ := json.Marshal(post)

		err := storage.Delete([]byte(pluginId), prefix)
		if err != nil {
			return post
		}

		err = storage.Put([]byte(pluginId), prefix, jsonBytes)
		if err != nil {
			return post
		}
	}
	return post
}

func isTimeExpired(post MatchPost) bool {

	nowLocal := time.Now().Local()

	lastUpTime, err := time.Parse("2006-01-02 15:04:05", post.LastUpdate)
	if err != nil {
		return true
	}
	lutUnix := time.Unix(lastUpTime.Unix(), 0)
	subHours := int(nowLocal.Sub(lutUnix))

	if subHours > 8 {
		return true
	}
	return false
}

func createMatchInfo(game Game, isToday bool) string {
	var matchDay string
	if isToday {
		matchDay = "今日"
	} else {
		matchDay = "明日"
	}

	teamBox := strings.Split(game.MatchNameB, "vs")
	teamA := strings.TrimSpace(teamBox[0])
	teamB := strings.TrimSpace(teamBox[1])

	scoreA := game.ScoreA
	scoreB := game.ScoreB

	matchStatus := game.MatchStatus
	gameTypeName := game.GameTypeName //e.g 常规赛

	var status string
	if matchStatus == "1" {
		status = "未开始"
	} else if matchStatus == "2" {
		status = "进行中"
	} else {
		status = "已结束"
	}

	return fmt.Sprintf("%s比赛(%s,%s) %s %s:%s %s", matchDay, gameTypeName, status, teamA, scoreA, scoreB, teamB)
}

func createMatchPost() MatchPost {
	var post MatchPost
	post.TodayPost = ""
	post.TommorrowPost = ""

	recentGamesUrl := "https://lpl.qq.com/web201612/data/LOL_MATCH2_MATCH_HOMEPAGE_BMATCH_LIST_148.js"
	resp, err := http.DefaultClient.Get(recentGamesUrl)

	if err != nil {
		return post
	}
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return post
	}
	respBodyStr := string(robots)
	if respBodyStr == "" {
		return post
	}
	var mBox MatchBox
	err = json.Unmarshal(robots, &mBox)
	if err != nil {
		return post
	}

	nowLocal := time.Now().Local()
	for _, game := range mBox.GameInfo {

		matchDate, err := time.Parse("2006-01-02 15:04:05", game.MatchDate)
		if err != nil {
			continue
		}
		if matchDate.Before(nowLocal) {
			continue
		}
		if matchDate.Day() == nowLocal.Day() {
			todayInfo := createMatchInfo(game, true)
			post.TodayPost += fmt.Sprintf("%s\n", todayInfo)
		}
		if matchDate.Day() == nowLocal.Day()+1 {
			tomorrowInfo := createMatchInfo(game, false)
			post.TommorrowPost += fmt.Sprintf("%s\n", tomorrowInfo)
		}

	}

	post.LastUpdate = nowLocal.String()
	if post.TodayPost == "" {
		post.TodayPost = "今日无赛事"
	}
	if post.TommorrowPost == "" {
		post.TommorrowPost = "明日无赛事"
	}

	return post
}

type MatchPost struct {
	TodayPost     string `json:"todayPost"`
	TommorrowPost string `json:"tomorrowPost"`
	LastUpdate    string `json:"lastUpTime"`
}

type MatchBox struct {
	Status     string `json:"status"`
	LastUpdate string `json:"lastUpTime"`
	GameInfo   []Game `json:"msg"`
}

type Game struct {
	IsTft         string `json:"isTft"`
	TftInfos      string `json:"tftInfos"`
	GameIdB       string `json:"bGameId"`
	MatchIdB      string `json:"bMatchId"`
	MatchNameB    string `json:"bMatchName"` //xx vs xx
	GameId        string `json:"GameId"`
	GameName      string `json:"GameName"`
	GameMode      string `json:"GameMode"`
	GameModeName  string `json:"GameModeName"`
	GameTypeId    string `json:"GameTypeId"`
	GameTypeName  string `json:"GameTypeName"` //常规赛
	GameProcId    string `json:"GameProcId"`
	GameProcName  string `json:"GameProcName"` //第几周
	TeamA         string `json:"TeamA"`
	ScoreA        string `json:"ScoreA"`
	TeamB         string `json:"TeamB"`
	ScoreB        string `json:"ScoreB"`
	MatchDate     string `json:"MatchDate"`   //比赛时间
	MatchStatus   string `json:"MatchStatus"` // 1 未开始 2 进行中
	MatchWin      string `json:"MatchWin"`
	IQTMatchId    string `json:"iQTMatchId"`
	AppTopicId    string `json:"AppTopicId"`
	AppShowFlag_m string `json:"AppShowFlag_m"`
	AppShowFlag_n string `json:"AppShowFlag_n"`
	NewsId        string `json:"NewsId"`
	ExtMsg        string `json:"sExt1"`
	Video1        string `json:"Video1"`
	Video2        string `json:"Video2"`
	Video3        string `json:"Video3"`
	Chat1         string `json:"Chat1"`
	Chat2         string `json:"Chat2"`
	Chat3         string `json:"Chat3"`
	News1         string `json:"News1"`
	News2         string `json:"News2"`
	News3         string `json:"News3"`
	IsFocus       string `json:"isFocus"`
}
