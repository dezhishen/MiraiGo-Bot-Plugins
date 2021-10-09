package rss

import "testing"

func TestFeed(t *testing.T) {
	url := "https://rssfeed.today/weibo/rss/7408951128"
	Feed(url)
}
