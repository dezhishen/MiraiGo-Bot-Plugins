package rss

import (
	"github.com/SlyMarbo/rss"
)

func Feed(url string) (*rss.Feed, error) {
	return rss.Fetch(url)
}
