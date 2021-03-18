package pixiv

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var randomUrl = "https://open.pixivic.net/wallpaper/%v/random?size=%v&domain=https://i.pixiv.cat&webp=0&detail=1"

func randomImage(platform, size, msgType string) (*[]byte, error) {
	url := fmt.Sprintf(randomUrl, platform, size)
	if msgType == "image" {
		if platform == "" {
			platform = "mobile"
		}
		r, err := http.DefaultClient.Get(url)
		if err != nil {
			return nil, err
		}
		robots, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		r.Body.Close()
		return &robots, nil
	}
	c := http.DefaultClient
	var imageSrc []byte
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		imageSrc = []byte(req.URL.String())
		fmt.Print(req.URL.String())
		return errors.New("stop by sendType")
	}
	_, err := c.Get(url)
	if imageSrc == nil {
		fmt.Println(err)
		return nil, errors.New("发生意外,请重试")
	}
	return &imageSrc, nil
}
