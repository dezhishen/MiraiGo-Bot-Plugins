package bilifan

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"

	"github.com/Mrs4s/MiraiGo/message"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
)

// Plugin b站粉丝数查看插件
type Plugin struct{
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:		"bilifan",
		Name: 	"b站粉丝数查看插件",
	}
}

// IsFireEvent 是否触发
func (p Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".bilifan")
	}
	return false
}

// OnMessageEvent OnMessageEvent
func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	result := &plugins.MessageResponse{
		Elements: make([]message.IMessageElement, 1),
	}
	// fetch and split the parameter
	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	params := strings.Split(context, " ")
	uid := ""
	if params[1] == "help"{
		result.Elements[0] = message.NewText(".bilifan UP主UID")
	}
	if len(params) < 2 {
		// return nil, errors.New("请输入需要查询的UP主的uid")
		// 不输入查询目标那就查我的吧嘤嘤嘤，走过路过点个关注不迷路鸭
		uid = "7528659"
	} else {
		uid = params[1]
	}
	txt,  _ := getBiliFan(uid)
	result.Elements[0] = message.NewText(txt)
	return result, nil
}

func init() {
	plugins.RegisterOnMessagePlugin(Plugin{})
}

func getBiliFan(uid string) (string, error){
	// request
	upNameUri := fmt.Sprintf("https://api.bilibili.com/x/space/acc/info?mid=%v&jsonp=jsonp", uid)
	fanUri := fmt.Sprintf("https://api.bilibili.com/x/relation/stat?vmid=%v&jsonp=jsonp", uid)
	
	// fetch the name of UP
	resp, err := http.DefaultClient.Get(upNameUri)
	if err != nil {
		return "", nil
	}
	// read the Body of requested data
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", nil
	}
	// convert the code to string
	respBodyStr := string(robots)
	if respBodyStr == "" {
		return "", nil
	}
	// decode json
	var biliupdata BiliUpData
	err = json.Unmarshal([]byte(respBodyStr), &biliupdata)
	if err != nil {
		return "", err
	}

	// fetch number of fans
	resp, err = http.DefaultClient.Get(fanUri)
	if err != nil {
		return "", nil
	}
	// read the Body of requested data
	robots, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return "", nil
	}
	// convert the code to string
	respBodyStr = string(robots)
	if respBodyStr == "" {
		return "", nil
	}
	// decode json
	var bilifandata BiliFanData
	err = json.Unmarshal([]byte(respBodyStr), &bilifandata)
	if err != nil {
		return "", err
	}

	// Produce the prefix
	prefix := ""
	switch uid {
	case "7528659":
		prefix = "脑子有饼的"
	case "2929582":
		prefix = "折纸之光"
	}

	txt := fmt.Sprintf(
		"%v %v\nuid：%v\n当前粉丝数：%v",
		prefix,
		biliupdata.Data.Name,
		uid,
		bilifandata.Data.Follower,
	)
	return txt, nil
}

type FanData struct{
    Mid             uint64
    Following       uint16
    Whisper         uint16
    Blcak           uint16
    Follower        uint64
}

type UpData struct{
    Mid             uint64
    Name            string
    Sex             string
    Face            string
    Sign            string
    Rank            uint32
    Level           uint8
    Jointime        uint8
    Moral           uint8
    Silence         uint8
    Birthday        string
    Coins           int8
    Fans_badge      bool
    Official        interface{}
    Vip             interface{}
    Pendant         interface{}
    Live_room       interface{}
}


type BiliFanData struct{
    Code            int8
    Message         string
    Ttl             int8
    Data            FanData
}

type BiliUpData struct{
    Code            int8
    Messages        string
    Ttl             int8
    Data            UpData
}