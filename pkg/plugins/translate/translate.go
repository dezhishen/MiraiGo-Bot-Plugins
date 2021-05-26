package translate

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/command"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "translate(baidu)",
		Name: "百度文本翻译插件",
	}
}

func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".tr") || field.Content == ".tr--help"
	}
	return false
}

type TranslateReq struct {
	From string `short:"f" long:"from" description:"来源语言" default:"auto"`
	To   string `short:"t" long:"to" description:"目标语言" default:"auto"`
}

var dictI18 = map[string]string{
	"zh":  "中文",
	"en":  "英语",
	"yue": "粤语",
	"wyw": "文言文",
	"jp":  "日语",
	"kor": "韩语",
	"fra": "法语",
	"spa": "西班牙语",
	"th":  "泰语",
	"ara": "阿拉伯语",
	"ru":  "俄语",
	"pt":  "葡萄牙语",
	"de":  "德语",
	"it":  "意大利语",
	"el":  "希腊语",
	"nl":  "荷兰语",
	"pl":  "波兰语",
	"bul": "保加利亚语",
	"est": "爱沙尼亚语",
	"dan": "丹麦语",
	"fin": "芬兰语",
	"cs":  "捷克语",
	"rom": "罗马尼亚语",
	"slo": "斯洛文尼亚语",
	"swe": "瑞典语",
	"hu":  "匈牙利语",
	"cht": "繁体中文",
	"vie": "越南语",
}

func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {

	var elements []message.IMessageElement

	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	trReq := TranslateReq{}
	commands, err := command.Parse(".tr", &trReq, strings.Split(context, " "))
	if err != nil {
		return nil, err
	}
	// 测试 根命令 .dict|.tr
	// rootCommand := commands[0]
	// 文本
	var q string
	for i := 1; i < len(commands); i++ {
		q = q + " " + commands[i]
	}
	q = strings.TrimSpace(q)
	from := trReq.From
	to := trReq.To
	transInfo, err := callHttp(q, from, to)
	if err != nil {
		return nil, err
	}
	re := transInfo.TransRe
	elements = append(elements, message.NewText(fmt.Sprintf("%v=>%v\n", dictI18[transInfo.FromLan], dictI18[transInfo.ToLan])))
	elements = append(elements, message.NewText(fmt.Sprintf("源文本\n%v\n", re[0].Source)))
	elements = append(elements, message.NewText(fmt.Sprintf("翻译文本\n%v\n", re[0].Destination)))

	result := &plugins.MessageResponse{
		Elements: elements,
	}

	return result, nil
}

func callHttp(q, from, to string) (*TransStruct, error) {
	salt := strconv.Itoa(rand.Intn(100000))
	uri := "http://api.fanyi.baidu.com/api/trans/vip/translate?"
	data := appid + q + salt + key
	signMd5 := md5.New()
	signMd5.Write([]byte(data))
	sign := hex.EncodeToString(signMd5.Sum(nil))
	uri += fmt.Sprintf("q=%v", url.QueryEscape(q))
	uri += fmt.Sprintf("&from=%v", from)
	uri += fmt.Sprintf("&to=%v", to)
	uri += fmt.Sprintf("&appid=%v", appid)
	uri += fmt.Sprintf("&salt=%v", salt)
	uri += fmt.Sprintf("&sign=%v", sign)
	fmt.Printf("%v\n", uri)
	resp, err := http.DefaultClient.Get(uri)
	if err != nil {
		return nil, err
	}
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	respBodyStr := string(robots)
	fmt.Printf("%v\n", respBodyStr)
	var transInfo TransStruct
	err = json.Unmarshal(robots, &transInfo)
	if err != nil {
		return nil, err
	}
	return &transInfo, nil
}

var key string

var appid string

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
	var e error
	key, e = getKey()
	if e != nil {
		fmt.Printf("读取百度翻译的key发生错误:[%v]", e.Error())
	}
	appid, e = getID()
	if e != nil {
		fmt.Printf("读取百度翻译的ID发生错误:[%v]", e.Error())
	}
}

type TransResult struct {
	Source      string `json:"src"`
	Destination string `json:"dst"`
}

type TransStruct struct {
	FromLan string        `json:"from"`
	ToLan   string        `json:"to"`
	TransRe []TransResult `json:"trans_result"`
}

func getID() (string, error) {
	str := os.Getenv("BOT_BAIDU_FANYI_ID")
	if str == "" {
		return "", errors.New("未配置百度翻译的appID")
	}
	return str, nil
}

func getKey() (string, error) {
	str := os.Getenv("BOT_BAIDU_FANYI_KEY")
	if str == "" {
		return "", errors.New("未配置百度翻译的key")
	}
	return str, nil
}
