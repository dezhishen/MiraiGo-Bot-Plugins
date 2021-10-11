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
	"github.com/dezhiShen/MiraiGo-Bot/pkg/command"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
	"github.com/sirupsen/logrus"
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
	v := msg.Elements[0]
	field, ok := v.(*message.TextElement)
	return ok && strings.HasPrefix(field.Content, ".rss")
}

type rssReq struct {
	Event string `short:"e" long:"event" description:"动作" default:"add"`
}

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
			if strings.TrimSpace(url) == "" {
				continue
			}
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
				elements = append(elements, message.NewText("移除成功:"+feed.Title+"\n"))
			} else {
				elements = append(elements, message.NewText("当前未订阅"))
			}
		}
	} else if req.Event == "list" {
		feeds := getAllFeed(request)
		if len(feeds) > 0 {
			for _, e := range feeds {
				elements = append(elements, message.NewText(
					e.Title+": "+e.UpdateURL+"\n"))
			}
		} else {
			elements = append(elements, message.NewText("当前无订阅"))
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

var logger = logrus.WithField("bot-plugin", "rss")

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
		var infos []info

		storage.GetByPrefix([]byte(t.PluginInfo().ID), []byte(rss_url_distributor+url), func(key, v []byte) error {
			var info info
			json.Unmarshal(v, &info)
			infos = append(infos, info)
			return nil
		})

		for _, e := range feeds {
			for _, info := range infos {
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
		}
	}
	// bot.QQClient.Send
	return nil
}

// Cron cron表达式
func (t Plugin) Cron() string {
	return "0 */5 * * * ?"
}
