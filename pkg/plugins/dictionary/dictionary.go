package dictionary

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

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
		ID:   "dictionary",
		Name: "字典插件",
	}
}

func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".dict ")
	}
	return false
}

func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {

	var elements []message.IMessageElement

	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	params := strings.Split(context, " ")
	q := params[1]

	uri := "http://dict-co.iciba.com/search.php?word=" + q
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
	root, _ := htmlquery.Parse(strings.NewReader(respBodyStr))
	brs := htmlquery.Find(root, "/html/body/text()")
	for _, row := range brs {
		ele := htmlquery.InnerText(row)
		ele = strings.TrimSpace(ele)
		if ele != "" {
			vals := strings.Split(ele, " ")
			for _, val := range vals {
				elements = append(elements, message.NewText(fmt.Sprintf("%v\n", val)))
			}
		}

	}

	result := &plugins.MessageResponse{
		Elements: elements,
	}

	return result, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}
