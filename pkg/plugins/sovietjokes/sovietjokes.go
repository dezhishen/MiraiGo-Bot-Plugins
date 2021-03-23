package sovietjokes

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"
	"unsafe"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/antchfx/htmlquery"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   "sovietjokes",
		Name: "苏联笑话插件",
	}
}

func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && field.Content == ".sj"
	}
	return false
}

func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	var elements []message.IMessageElement
	//todo:
	randJoke := getRandomJoke()
	if unsafe.Sizeof(randJoke) == 0 {
		elements = append(elements, message.NewText("没有笑话库存了,你有兴趣加入吗?"))
	} else {
		elements = append(elements, message.NewText(fmt.Sprintf("%s\n", randJoke.Title)))
		elements = append(elements, message.NewText(fmt.Sprintf("%s", randJoke.Content)))
	}

	result := &plugins.MessageResponse{
		Elements: elements,
	}
	return result, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}

func getRandomJoke() Joke {

	jokes := getJokes()
	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(len(jokes) - 1)
	randJoke := jokes[v]

	return randJoke
}

func getJokes() []Joke {

	var jokes []Joke
	var emptyJk []Joke
	//检测根目录是否存在 sovietjokes.json ,不存在进行一次爬取
	if !sjExists() {

		jokes = analyzeJokes()

		jbytes, err := json.Marshal(&jokes)

		if err != nil {
			return emptyJk
		}

		err = ioutil.WriteFile(_joke_path_, jbytes, 0777)
		if err != nil {
			return emptyJk
		}
	} else {
		//解析文件
		jData, err := ioutil.ReadFile(_joke_path_)

		if err != nil {
			return emptyJk
		}

		datajson := []byte(jData)

		err = json.Unmarshal(datajson, &jokes)

		if err != nil {
			return emptyJk
		}
	}
	return jokes
}

func analyzeJokes() []Joke {

	var jokes []Joke

	sjUrl := "https://library.moegirl.org.cn/%E8%8B%8F%E8%81%94%E7%AC%91%E8%AF%9D"
	html := getHtml(sjUrl)
	if html == "" {
		return jokes
	}

	root, _ := htmlquery.Parse(strings.NewReader(html))
	heads := htmlquery.Find(root, "//span[@class='mw-headline']")

	for _, row := range heads {

		var j Joke
		//ttstr := strings.Split(fmt.Sprintf("%s", htmlquery.InnerText(row)), " ")[1]

		j.Title = fmt.Sprintf("%s", htmlquery.InnerText(row))
		jokes = append(jokes, j)
	}

	contents := htmlquery.Find(root, "//div[@class='poem']/p")
	for index, row := range contents {

		jokes[index].Content = fmt.Sprintf("%s", htmlquery.InnerText(row))
	}

	return jokes
}

func sjExists() bool {

	_, err := os.Stat(_joke_path_)

	if err != nil {
		return false
	}

	if os.IsNotExist(err) {
		return false
	}

	return true
}

func getHtml(url_ string) string {
	req, _ := http.NewRequest("GET", url_, nil)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/76.0.3776.0 Safari/537.36")
	client := &http.Client{Timeout: time.Second * 5}
	resp, err := client.Do(req)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil && data == nil {
		return ""
	}
	return fmt.Sprintf("%s", data)
}

type Joke struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

const _joke_path_ = "./sovietJokes.json"
