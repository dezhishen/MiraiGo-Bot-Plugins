package rss

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/Logiase/MiraiGo-Template/bot"
	"github.com/Mrs4s/MiraiGo/client"
	"github.com/Mrs4s/MiraiGo/message"
	"github.com/PuerkitoBio/goquery"
	"github.com/SlyMarbo/rss"
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

var logger = logrus.WithField("bot-plugin", "rss")
var cron = "0 */15 * * * ?"
var pluginId = ".rss"

// PluginInfo PluginInfo
func (w Plugin) PluginInfo() *plugins.PluginInfo {
	return &plugins.PluginInfo{
		ID:   pluginId,
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
			err := listenFeed(url, request)
			if err != nil {
				return nil, err
			}
			elements = append(elements, message.NewText("订阅成功:"+url))
		}
	} else if req.Event == "del" {
		for i := 1; i < len(commands); i++ {
			url := commands[i]
			// prefix := []byte(rss_prefix + url)
			// storage.Put([]byte(w.PluginInfo().ID), prefix, []byte(url))
			err := unListenFeed(url, request)
			if err != nil {
				return nil, err
			}
			elements = append(elements, message.NewText("移除成功:"+url))
		}
	} else if req.Event == "list" {
		urls := getAllFeed(request)
		if len(urls) > 0 {
			for _, e := range urls {
				elements = append(elements, message.NewText(e+"\n"))
			}
		} else {
			elements = append(elements, message.NewText("当前无订阅"))
		}
	} else if req.Event == "update" {
		now := int64(time.Now().Unix()/900) * 900
		for i := 1; i < len(commands); i++ {
			url := commands[i]
			items, err := updateFeed(url, now)
			if err != nil {
				return nil, errors.New("更新失败" + err.Error())
			}
			feeds := items2Feeds(items)
			for _, f := range feeds {
				if request.MessageType == plugins.GroupMessage {
					telements, _ := feed2MessageElements(f, request.QQClient, "group", request.GroupCode)
					elements = append(elements, telements...)
				} else {
					telements, _ := feed2MessageElements(f, request.QQClient, "private", request.Sender.Uin)
					elements = append(elements, telements...)
				}
				if err != nil {
					return nil, err
				}
			}
		}
		if len(elements) == 0 {
			elements = append(elements, message.NewText("暂无更新"))
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
	theDate := time.Now().Unix() - 15*60
	for _, url := range urls {
		items, _ := updateFeed(url, theDate)
		feeds := items2Feeds(items)
		var infos []info
		storage.GetByPrefix([]byte(t.PluginInfo().ID), []byte(rss_url_distributor+url), func(key, v []byte) error {
			var info info
			json.Unmarshal(v, &info)
			infos = append(infos, info)
			return nil
		})

		for _, info := range infos {
			for _, oneFeed := range feeds {
				elemens, err := feed2MessageElements(oneFeed, bot.QQClient, info.Type, info.Code)
				if err != nil {
					return err
				}
				if info.Type == "group" {
					bot.SendGroupMessage(info.Code, &message.SendingMessage{
						Elements: elemens,
					})
				} else {
					bot.SendPrivateMessage(info.Code, &message.SendingMessage{
						Elements: elemens,
					})
				}
			}

		}
	}
	// bot.QQClient.Send
	return nil
}

// Cron cron表达式
func (t Plugin) Cron() string {
	return cron
}

func items2Feeds(items []*rss.Item) []oneOfFeed {
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
	return feeds
}
func feed2MessageElements(oneOfFeed oneOfFeed, client *client.QQClient, messageType string, code int64) ([]message.IMessageElement, error) {
	var messageElement []message.IMessageElement
	messageElement = append(messageElement, message.NewText(oneOfFeed.Title+"\n"))
	if messageType == "group" {
		// sendingMessage := &message.SendingMessage{}
		if oneOfFeed.CoverByte != nil {
			image, _ := client.UploadGroupImage(code, bytes.NewReader(oneOfFeed.CoverByte))
			messageElement = append(messageElement, image)
		}
		messageElement = append(messageElement, message.NewText(oneOfFeed.Link))
	} else {
		if oneOfFeed.CoverByte != nil {
			image, _ := client.UploadPrivateImage(code, bytes.NewReader(oneOfFeed.CoverByte))
			messageElement = append(messageElement, image)
		}
		messageElement = append(messageElement, message.NewText(oneOfFeed.Link))
	}
	return messageElement, nil
}
