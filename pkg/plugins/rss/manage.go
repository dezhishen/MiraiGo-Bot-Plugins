package rss

import (
	"encoding/json"
	"fmt"

	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
	"github.com/mmcdole/gofeed"
)

var rss_prefix string = "rss.url:"
var rss_url_distributor string = "rss-url.distributor:"

func listenFeed(url string, req *plugins.MessageRequest) error {
	rss_url_distributor_key :=
		rss_url_distributor + url + string(req.MessageType)
	distributorInfo := &info{
		Type: string(req.MessageType),
	}
	if req.MessageType == "group" {
		rss_url_distributor_key += fmt.Sprintf(":%v", req.GroupCode)
		distributorInfo.Code = req.GroupCode
	} else {
		rss_url_distributor_key += fmt.Sprintf(":%v", req.Sender.Uin)
		distributorInfo.Code = req.Sender.Uin
	}
	storage.Put([]byte(pluginId), []byte(rss_prefix+url), []byte(url))
	jsonBytes, _ := json.Marshal(distributorInfo)
	storage.Put([]byte(pluginId), []byte(rss_url_distributor_key), jsonBytes)
	return nil
}
func unListenFeed(url string, req *plugins.MessageRequest) error {
	rss_url_distributor_key :=
		rss_url_distributor + url + string(req.MessageType)
	if req.MessageType == "group" {
		rss_url_distributor_key += fmt.Sprintf(":%v", req.GroupCode)
	} else {
		rss_url_distributor_key += fmt.Sprintf(":%v", req.Sender.Uin)
	}
	storage.Delete([]byte(pluginId), []byte(rss_url_distributor_key))
	var hasRss = false
	storage.GetByPrefix([]byte(pluginId), []byte(rss_url_distributor+url), func(b1, b2 []byte) error {
		if hasRss {
			return nil
		}
		hasRss = true
		return nil
	})
	if !hasRss {
		storage.Delete([]byte(pluginId), []byte(rss_prefix+url))
	}
	return nil
}

func getFeed(url string) (*gofeed.Feed, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		logger.Infof("抓取Feed时发生异常:%v", url)
	}
	return feed, err
}

func updateFeed(url string, d int64) ([]*gofeed.Item, error) {
	logger.Infof("开始抓取更新:%s", url)
	feed, err := getFeed(url)
	if err != nil {
		return nil, err
	}
	var results []*gofeed.Item
	// lastDateByte, _ := storage.GetValue([]byte(pluginId), []byte(rss_url_date+url))
	if d >= feed.UpdatedParsed.Unix() {
		return results, nil
	}
	for _, e := range feed.Items {
		var now = e.PublishedParsed.Unix()
		if now <= d {
			continue
		}
		results = append(results, e)
	}
	logger.Infof("数量:%s", len(results))
	logger.Infof("结束更新:%s", url)
	return results, nil
}

func getAllFeed(req *plugins.MessageRequest) []string {
	var urls []string
	storage.GetByPrefix([]byte(pluginId), []byte(rss_prefix), func(b1, b2 []byte) error {
		urls = append(urls, string(b2))
		return nil
	})
	var result []string
	for _, url := range urls {
		rss_url_distributor_key :=
			rss_url_distributor + url + string(req.MessageType)
		if req.MessageType == "group" {
			rss_url_distributor_key += fmt.Sprintf(":%v", req.GroupCode)
		} else {
			rss_url_distributor_key += fmt.Sprintf(":%v", req.Sender.Uin)
		}
		v, _ := storage.GetValue([]byte(pluginId), []byte(rss_url_distributor_key))
		if v != nil {
			result = append(result, url)
		}
	}
	return result
}
