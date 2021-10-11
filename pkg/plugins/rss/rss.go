package rss

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/PuerkitoBio/goquery"
	"github.com/SlyMarbo/rss"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/command"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
)

// Plugin menhear
type Plugin struct {
	plugins.NoSortPlugin
	plugins.NoInitPlugin
	plugins.AlwaysNotFireNextEventPlugin
}

// PluginInfo PluginInfo
func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   ".rss",
		Name: "rss",
	}
}

// IsFireEvent 是否触发
func (w Plugin) IsFireEvent(msg *plugins.MessageRequest) bool {
	if len(msg.Elements) > 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && strings.HasPrefix(field.Content, ".rss")
	}
	return false
}

type rssReq struct {
	Event string `short:"e" long:"event" description:"动作" default:"add"`
}

var rss_prefix string = "rss.url:"
var rss_url_distributor string = "rss-url.distributor:"

// OnMessageEvent OnMessageEvent
func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	var elements []message.IMessageElement
	context := ""
	for _, v := range request.Elements {
		if v.Type() == message.Text {
			field, _ := v.(*message.TextElement)
			context += (field.Content)
		}
	}
	req := rssReq{}
	commands, err := command.Parse(".rss", &req, strings.Split(context, " "))
	if err != nil {
		return nil, err
	}
	if req.Event == "add" {
		for i := 1; i < len(commands); i++ {
			url := commands[i]
			// prefix := []byte(rss_prefix + url)
			// storage.Put([]byte(w.PluginInfo().ID), prefix, []byte(url))
			feed, err := setFeed(url, request)
			if err != nil {
				return nil, err
			}
			if feed != nil {
				elements = append(elements, message.NewText("订阅成功:"+feed.Title))
			}
		}
	} else if req.Event == "del" {
		for i := 1; i < len(commands); i++ {
			url := commands[i]
			// prefix := []byte(rss_prefix + url)
			// storage.Put([]byte(w.PluginInfo().ID), prefix, []byte(url))
			feed, err := removeFeed(url, request)
			if err != nil {
				return nil, err
			}
			if feed != nil {
				elements = append(elements, message.NewText("移除成功:"+feed.Title))
			}
		}
	}
	return &plugins.MessageResponse{
		Elements: elements,
	}, nil
}

func init() {
	plugin := Plugin{}
	plugins.RegisterOnMessagePlugin(plugin)
	plugins.RegisterSchedulerPlugin(plugin)
}

type info struct {
	Type string `json:"type"`
	Code int64  `json:"code"`
}
type oneOfFeed struct {
	Title     string
	Link      string
	CoverSrc  string
	CoverByte []byte
}

// Run 回调
func (t Plugin) Run(bot *bot.Bot) error {
	prefix := []byte(rss_prefix)
	var urls []string
	storage.GetByPrefix([]byte(t.PluginInfo().ID), prefix, func(key, v []byte) error {
		url := string(v)
		urls = append(urls, url)
		return nil
	})
	for _, url := range urls {
		items, _ := update(url)
		var feeds []oneOfFeed
		for _, feedItem := range items {
			doc, err := goquery.NewDocumentFromReader(strings.NewReader(feedItem.Summary))
			if err != nil {
				continue
			}
			e := oneOfFeed{
				Title: feedItem.Title,
				Link:  feedItem.Link,
			}
			images := doc.Find("img")
			if images != nil && len(images.Nodes) > 0 {
				coverSrc, exists := goquery.NewDocumentFromNode(images.Nodes[0]).Attr("src")
				if exists {
					e.CoverSrc = coverSrc
					r, _ := http.DefaultClient.Get(coverSrc)
					content, _ := ioutil.ReadAll(r.Body)
					e.CoverByte = content
				}
			}
			feeds = append(feeds, e)
		}
		storage.GetByPrefix([]byte(t.PluginInfo().ID), []byte(rss_url_distributor+url), func(key, v []byte) error {
			var info info
			json.Unmarshal(v, &info)

			for _, e := range feeds {
				if info.Type == "group" {
					sendingMessage := &message.SendingMessage{}
					sendingMessage.Append(message.NewText(e.Title + "\n"))
					if e.CoverByte != nil {
						image, _ := bot.QQClient.UploadGroupImage(int64(info.Code), bytes.NewReader(e.CoverByte))
						sendingMessage.Append(image)
					}
					sendingMessage.Append(message.NewText(e.Link))
					bot.SendGroupMessage(int64(info.Code), sendingMessage)
				} else {
					sendingMessage := &message.SendingMessage{}
					sendingMessage.Append(message.NewText(e.Title + "\n"))
					if e.CoverByte != nil {
						image, _ := bot.QQClient.UploadPrivateImage(int64(info.Code), bytes.NewReader(e.CoverByte))
						sendingMessage.Append(image)
					}
					sendingMessage.Append(message.NewText(e.Link))
					bot.SendPrivateMessage(int64(info.Code), sendingMessage)
				}
			}
			return nil
		})
	}
	// bot.QQClient.Send
	return nil
}

// Cron cron表达式
func (t Plugin) Cron() string {
	return "0 */5 * * * ?"
}

var allFeed = make(map[string]*rss.Feed)

func update(url string) ([]*rss.Item, error) {
	feed, ok := getFeed(url)
	if !ok {
		feed, _ = rss.Fetch(url)
		allFeed[url] = feed
	}
	feed.Update()
	var results []*rss.Item
	lastDate, _ := storage.GetValue([]byte(".rss"), []byte(rss_prefix+url+":last"))
	for i, e := range feed.Items {
		if lastDate != nil {
			tDate := storage.BytesToInt(lastDate)
			if int(e.Date.Unix()) <= tDate {
				break
			}
		}
		if i == 0 {
			storage.Put([]byte(".rss"), []byte(rss_prefix+url+":last"), storage.IntToBytes(int(e.Date.Unix())))
		}
		results = append(results, e)
	}
	return results, nil
}

func getFeed(url string) (*rss.Feed, bool) {
	feed, ok := allFeed[url]
	return feed, ok
}

func setFeed(url string, req *plugins.MessageRequest) (*rss.Feed, error) {
	feed, ok := getFeed(url)
	rss_url_distributor_key :=
		rss_url_distributor + url + string(req.MessageType)

	distributorInfo := &info{
		Type: string(req.MessageType),
	}
	if req.MessageType == "group" {
		rss_url_distributor_key += string(rune(req.GroupCode))
		distributorInfo.Code = req.GroupCode
	} else {
		rss_url_distributor_key += string(rune(req.Sender.Uin))
		distributorInfo.Code = req.Sender.Uin
	}
	if !ok {
		feed, err := rss.Fetch(url)
		if err != nil {
			return nil, err
		}
		allFeed[url] = feed
		storage.Put([]byte(".rss"), []byte(rss_prefix+url), []byte(url))
	}
	jsonBytes, _ := json.Marshal(distributorInfo)
	storage.Put([]byte(".rss"), []byte(rss_url_distributor_key), jsonBytes)
	return feed, nil
}

func removeFeed(url string, req *plugins.MessageRequest) (*rss.Feed, error) {
	feed, ok := getFeed(url)
	if ok {
		rss_url_distributor_key :=
			rss_url_distributor + url + string(req.MessageType)
		if req.MessageType == "group" {
			rss_url_distributor_key += string(rune(req.GroupCode))
		} else {
			rss_url_distributor_key += string(rune(req.Sender.Uin))
		}
		storage.Delete([]byte(".rss"), []byte(rss_url_distributor_key))
		var hasRss = false
		storage.GetByPrefix([]byte(".rss"), []byte(rss_url_distributor+url), func(b1, b2 []byte) error {
			if hasRss {
				return nil
			}
			hasRss = true
			return nil
		})
		if !hasRss {
			storage.Delete([]byte(".rss"), []byte(rss_prefix+url))
		}
	}
	return feed, nil
}
