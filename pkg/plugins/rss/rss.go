package rss

import (
	"bytes"
	"encoding/json"
	"fmt"
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
	if len(msg.Elements) == 1 && msg.Elements[0].Type() == message.Text {
		v := msg.Elements[0]
		field, ok := v.(*message.TextElement)
		return ok && field.Content == ".rss"
	}
	return false
}

type rssReq struct {
	Event string `short:"e" long:"event" description:"动作" default:"add"`
}

// OnMessageEvent OnMessageEvent
func (w Plugin) OnMessageEvent(request *plugins.MessageRequest) (*plugins.MessageResponse, error) {
	var elements []message.IMessageElement
	v := request.Elements[0]
	field, _ := v.(*message.TextElement)
	context := field.Content
	req := rssReq{}
	commands, err := command.Parse(".rss", &req, strings.Split(context, " "))
	if err != nil {
		return nil, err
	}
	if req.Event == "add" {
		for i := 1; i < len(commands); i++ {
			url := commands[i]
			prefix := []byte("rss.url:" + url)
			storage.Put([]byte(w.PluginInfo().ID), prefix, []byte(url))
		}
	}
	elements = append(elements, message.NewText("订阅成功"))
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
	Type string
	Code uint64
}
type oneOfFeed struct {
	Title     string
	Link      string
	CoverSrc  string
	CoverByte []byte
}

// Run 回调
func (t Plugin) Run(bot *bot.Bot) error {
	prefix := []byte("rss.url")
	var urls []string
	storage.GetByPrefix([]byte(t.PluginInfo().ID), prefix, func(key, v []byte) error {
		url := string(v)
		urls = append(urls, url)
		return nil
	})
	for _, url := range urls {
		prefix := []byte(fmt.Sprintf("rss-url.distributor.%v", url))
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
					r, _ := http.DefaultClient.Get(url)
					content, _ := ioutil.ReadAll(r.Body)
					e.CoverByte = content
				}
			}
			feeds = append(feeds, e)
		}
		storage.GetByPrefix([]byte(t.PluginInfo().ID), prefix, func(key, v []byte) error {
			var info info
			json.Unmarshal(v, &info)
			if info.Type == "group" {
				sendingMessage := &message.SendingMessage{}
				for _, e := range feeds {
					sendingMessage.Append(message.NewText(e.Title + "\n"))
					image, _ := bot.QQClient.UploadGroupImage(int64(info.Code), bytes.NewReader(e.CoverByte))
					sendingMessage.Append(image)
				}
				bot.SendGroupMessage(int64(info.Code), sendingMessage)
			} else {
				sendingMessage := &message.SendingMessage{}
				for _, e := range feeds {
					sendingMessage.Append(message.NewText(e.Title + "\n"))
					image, _ := bot.QQClient.UploadPrivateImage(int64(info.Code), bytes.NewReader(e.CoverByte))
					sendingMessage.Append(image)
				}
				bot.SendGroupMessage(int64(info.Code), sendingMessage)
			}
			return nil
		})
	}
	// bot.QQClient.Send
	return nil
}

// Cron cron表达式
func (t Plugin) Cron() string {
	return "0 */1 * * * ?"
}

var allFeed = make(map[string]*rss.Feed)

func update(url string) ([]*rss.Item, error) {
	feed, ok := allFeed[url]
	if !ok {
		feed, _ = rss.Fetch(url)
		allFeed[url] = feed
	}
	feed.Update()
	return feed.Items, nil
}
