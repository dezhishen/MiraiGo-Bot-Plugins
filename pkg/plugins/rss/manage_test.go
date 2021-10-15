package rss

import (
	"testing"
)

func Test_getFeed(t *testing.T) {
	feed1, err := getFeed("https://nitter.namazso.eu/212moving/rss")
	if err != nil {
		panic(err)
	}
	println(feed1)
	feed2, err := getFeed("https://nitter.namazso.eu/212moving/rss")
	if err != nil {
		panic(err)
	}
	print(feed2)
}
