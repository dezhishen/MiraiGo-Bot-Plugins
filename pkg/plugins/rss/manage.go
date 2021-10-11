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
var rss_url_date string = "rss-url.date:"

func update(url string) ([]*rss.Item, error) {
	feed, ok := getFeed(url)
	if !ok {
		return nil, errors.New("订阅地址不存在")
	}
	feed.Update()
	var results []*rss.Item
	lastDate, _ := storage.GetValue([]byte(".rss"), []byte(rss_url_date+url))
	for i, e := range feed.Items {
		if lastDate != nil {
			tDate := storage.BytesToInt(lastDate)
			if int(e.Date.Unix()) <= tDate {
				break
			}
		}
		if i == 0 {
			storage.Put([]byte(".rss"), []byte(rss_url_date+url), storage.IntToBytes(int(e.Date.Unix())))
		}
		results = append(results, e)
	}
	return results, nil
}

func getFeed(url string) (*rss.Feed, bool) {
	var feed *rss.Feed
	b, _ := storage.GetValue([]byte(".rss"), []byte(rss_prefix+url))
	if b == nil {
		return nil, false
	}
	tUrl := string(b)
	var err error
	feed, err = rss.Fetch(tUrl)
	if err != nil {
		return nil, false
	}
	return feed, true
}

func setFeed(url string, req *plugins.MessageRequest) (*rss.Feed, error) {
	feed, ok := getFeed(url)
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
	if !ok {
		var err error
		feed, err = rss.Fetch(url)
		if err != nil {
			return nil, err
		}
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
			rss_url_distributor_key += fmt.Sprintf(":%v", req.GroupCode)
		} else {
			rss_url_distributor_key += fmt.Sprintf(":%v", req.Sender.Uin)
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

func getAllFeed(req *plugins.MessageRequest) []*rss.Feed {
	var urls []string
	storage.GetByPrefix([]byte(".rss"), []byte(rss_prefix), func(b1, b2 []byte) error {
		urls = append(urls, string(b2))
		return nil
	})
	var result []*rss.Feed
	for _, url := range urls {
		rss_url_distributor_key :=
			rss_url_distributor + url + string(req.MessageType)
		if req.MessageType == "group" {
			rss_url_distributor_key += fmt.Sprintf(":%v", req.GroupCode)
		} else {
			rss_url_distributor_key += fmt.Sprintf(":%v", req.Sender.Uin)
		}
		v, _ := storage.GetValue([]byte(".rss"), []byte(rss_url_distributor_key))
		if v != nil {
			f, ok := getFeed(url)
			if ok {
				result = append(result, f)
			}
		}
	}
	return result
}
