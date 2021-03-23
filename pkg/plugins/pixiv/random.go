package pixiv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/dezhiShen/MiraiGo-Bot/pkg/cache"
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
		if r.StatusCode == 404 {
			log.Print("发生404错误")
			return nil, nil
		}
		robots, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return nil, err
		}
		r.Body.Close()
		return &robots, nil
	}
	c := http.Client{}
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

type RankingData struct {
	Title    string `json:"title"`
	UserId   string `json:"user_id"`
	UserName string `json:"user_name"`
	IllustID string `json:"illust_id"`
	Url      string `json:"url"`
}

type PictureData struct {
	Error             string   `json:"error"`
	Success           bool     `json:"success"`
	OriginalUrlsProxy []string `json:"original_urls_proxy"`
	OriginalUrlProxy  string   `json:"original_url_proxy"`
}

type Picture struct {
	Title    string   `json:"title"`
	UserId   string   `json:"userId"`
	UserName string   `json:"userName"`
	IllustID string   `json:"illustId"`
	Urls     []string `json:"urls"`
}

func randomAImage() (*RankingData, error) {
	randomUrl := "https://api.loli.st/pixiv/random.php?type=json&r18=true"
	r, err := http.DefaultClient.Get(randomUrl)
	if err != nil {
		return nil, err
	}

	robots, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	var data RankingData
	err = json.Unmarshal(robots, &data)
	return &data, err
}

func random(key string) (*Picture, error) {
	var data *RankingData
	var err error
	for i := 0; i < 10; i++ {
		data, err = randomAImage()
		if err != nil {
			return nil, err
		}
		theKey := fmt.Sprintf("pixiv_exists.%v.%v", key, data.IllustID)
		v, ok := cache.Get(theKey)
		if !ok || v == "N" {
			cache.Set(theKey, "Y", 24*time.Hour)
			break
		}
	}
	urlData := url.Values{
		"p": []string{data.IllustID},
	}
	r, err := http.DefaultClient.PostForm("https://api.pixiv.cat/v1/generate", urlData)
	if err != nil {
		return nil, err
	}
	robots, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	var pictureData PictureData
	err = json.Unmarshal(robots, &pictureData)
	if err != nil {
		return nil, err
	}
	if !pictureData.Success {
		return nil, errors.New(pictureData.Error)
	}
	result := &Picture{
		Title:    data.Title,
		UserId:   data.UserId,
		UserName: data.UserName,
		IllustID: data.IllustID,
	}
	if len(pictureData.OriginalUrlsProxy) > 0 {
		result.Urls = append(result.Urls, pictureData.OriginalUrlsProxy...)
	} else {
		result.Urls = append(result.Urls, pictureData.OriginalUrlProxy)
	}
	return result, nil
}

func getRank() (*RankingData, error) {
	r18Url := "https://api.loli.st/pixiv/?mode=daily_r18"
	r, err := http.DefaultClient.Get(r18Url)
	if err != nil {
		return nil, err
	}

	robots, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	var data RankingData
	err = json.Unmarshal(robots, &data)
	return &data, err
}

func randomR18(key string) (*Picture, error) {
	var data *RankingData
	var err error
	for i := 0; i < 10; i++ {
		data, err = getRank()
		if err != nil {
			return nil, err
		}
		theKey := fmt.Sprintf("pixiv_r18_exists.%v.%v", key, data.IllustID)
		v, ok := cache.Get(theKey)
		if !ok || v == "N" {
			cache.Set(theKey, "Y", 24*time.Hour)
			break
		}
	}
	urlData := url.Values{
		"p": []string{data.IllustID},
	}
	r, err := http.DefaultClient.PostForm("https://api.pixiv.cat/v1/generate", urlData)
	if err != nil {
		return nil, err
	}
	robots, err := ioutil.ReadAll(r.Body)
	r.Body.Close()
	if err != nil {
		return nil, err
	}
	var pictureData PictureData
	err = json.Unmarshal(robots, &pictureData)
	if err != nil {
		return nil, err
	}
	if !pictureData.Success {
		return nil, errors.New(pictureData.Error)
	}
	result := &Picture{
		Title:    data.Title,
		UserId:   data.UserId,
		UserName: data.UserName,
		IllustID: data.IllustID,
	}
	if len(pictureData.OriginalUrlsProxy) > 0 {
		result.Urls = append(result.Urls, pictureData.OriginalUrlsProxy...)
	} else {
		result.Urls = append(result.Urls, pictureData.OriginalUrlProxy)
	}
	return result, nil
}
