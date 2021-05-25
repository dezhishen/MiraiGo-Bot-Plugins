package translate

import (
	"bytes"
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
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
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

func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {

	var elements []message.IMessageElement

	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	dic := getLanDic()
	var q = ""
	from := "auto"
	to := "auto"
	salt := strconv.Itoa(rand.Intn(100000))

	if context == ".tr--help" {

		elements = append(elements, message.NewText(fmt.Sprintf("指令形式 .tr [目标语言代码] [文本]。 如.tr [en] [baba]。\n可以不指定目标语言，如.tr [baba]。\n")))
		helpLan := "支持语言代码:\n"
		for k, v := range dic {
			helpLan += fmt.Sprintf("%v : %v\n", k, v)
		}
		elements = append(elements, message.NewText(helpLan))

		result := &plugins.MessageResponse{
			Elements: elements,
		}
		return result, nil
	}
	test := context
	lanset, lanStr := languageSet(test, dic)
	fmt.Printf("%v\n", lanset)
	if lanset {
		to = lanStr
		lanText := context
		fmt.Printf("%v\n", lanText)
		q = lanText[strings.Index(lanText, "]")+3 : strings.LastIndex(lanText, "]")]
	} else {
		pureText := context
		fmt.Printf("%v\n", pureText)
		q = pureText[strings.Index(pureText, "[")+1 : strings.LastIndex(pureText, "]")]
	}

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

	re := transInfo.TransRe
	elements = append(elements, message.NewText(fmt.Sprintf("%v=>%v\n", dic[transInfo.FromLan], dic[transInfo.ToLan])))
	elements = append(elements, message.NewText(fmt.Sprintf("源文本\n%v\n", re[0].Source)))
	elements = append(elements, message.NewText(fmt.Sprintf("翻译文本\n%v\n", re[0].Destination)))

	result := &plugins.MessageResponse{
		Elements: elements,
	}

	return result, nil
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

func gbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GBK.NewDecoder())
	d, e := ioutil.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	return d, nil
}

func languageSet(lanstr string, dic map[string]string) (bool, string) {
	lanstr = lanstr[strings.Index(lanstr, "[")+1 : strings.Index(lanstr, "]")]
	if lanstr == "" {
		return false, ""
	}
	fmt.Println(dic[lanstr])
	if _, ok := dic[lanstr]; ok {
		return true, lanstr
	}
	return false, ""
}

func getLanDic() map[string]string {

	dicLan := map[string]string{
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
	return dicLan

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
