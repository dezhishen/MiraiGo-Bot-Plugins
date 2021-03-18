package pixiv

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var randomUrl = "https://open.pixivic.net/wallpaper/%v/random?size=%v&domain=https://i.pixiv.cat&webp=0&detail=1"

func randomImage(platform, size string) (*[]byte, error) {
	if platform == "" {
		platform = "mobile"
	}
	r, err := http.DefaultClient.Get(fmt.Sprintf(randomUrl, platform, size))
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
