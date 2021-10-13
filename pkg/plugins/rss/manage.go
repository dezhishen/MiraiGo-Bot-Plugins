package rss

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/SlyMarbo/rss"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/plugins"
	"github.com/dezhiShen/MiraiGo-Bot/pkg/storage"
)

var rss_prefix string = "rss.url:"
var rss_url_distributor string = "rss-url.distributor:"

var feeds = make(map[string]*rss.Feed)

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

func getFeed(url string) (*rss.Feed, bool, error) {
	feed, ok := feeds[url]
	if !ok {
		feed = initFeed(url)
		if feed == nil {
			return nil, false, errors.New("未订阅的地址")
		}
	}
	return feed, ok, nil
}

func updateFeed(url string, d int64) ([]*rss.Item, error) {
	logger.Infof("开始抓取更新:%s", url)
	feed, ok, err := getFeed(url)
	if err != nil {
		return nil, err
	}
	if ok {
		feed.Update()
	}
	var results []*rss.Item
	// lastDateByte, _ := storage.GetValue([]byte(pluginId), []byte(rss_url_date+url))
	lastDate := d
	for _, e := range feed.Items {
		var now = e.Date.Unix()
		if now <= lastDate {
			continue
		}
		results = append(results, e)
	}
	logger.Infof("数量:%s", len(results))
	logger.Infof("结束更新:%s", url)
	return results, nil
}

func initFeed(url string) *rss.Feed {
	var feed *rss.Feed
	b, _ := storage.GetValue([]byte(pluginId), []byte(rss_prefix+url))
	if b == nil {
		return nil
	}
	var err error
	feed, err = rss.Fetch(url)
	if err != nil {
		logger.Infof("初始化feed流,发生异常[%s],url:[%s]", err.Error(), url)
		return nil
	}
	return feed
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
			feed, ok := feeds[url]
			if ok {
				result = append(result, feed.Title+":"+url)
			} else {
				result = append(result, url)
			}
		}
	}
	return result
}
