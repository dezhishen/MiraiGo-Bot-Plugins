package rss

import (
	"testing"
	"time"
)

func Test_getFeed(t *testing.T) {
	// feed1, err := getFeed("http://192.168.31.104/weibo/user/7505824632")
	// if err != nil {
	// 	panic(err)
	// }
	// println(feed1)
	feed2, err := updateFeed("https://nitter.namazso.eu/212moving/rss", time.Now().Unix())
	if err != nil {
		panic(err)
	}
	print(feed2)
}
