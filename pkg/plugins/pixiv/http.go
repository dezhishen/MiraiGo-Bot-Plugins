package pixiv

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
)

func GetImage() (*[]byte, error) {
	return nil, nil
}

var setTuUrl = "https://api.acgmx.com/public/setu?type=json&ranking_type=illust"

type Resp struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data *data  `json:"data"`
}

type data struct {
	Illust string `json:"illust"`
	Title  string `json:"title"`
	Large  string `json:"large"`
	User   *user  `json:"user"`
}

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type ImageForSend struct {
	Title      string
	SourceFrom string
	Image      *[]byte
	UserName   string
	IllustID   string
}

func getSetu() (*ImageForSend, error) {
	req, err := http.NewRequest("GET", setTuUrl, nil)
	if err != nil {
		return nil, err
	}
	token := os.Getenv("BOT_PIXIV_TOKEN")
	req.Header.Add("token", token)
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	robots, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var resp Resp
	err = json.Unmarshal(robots, &resp)
	if err != nil {
		return nil, err
	}
	if !(resp.Code == 200 || resp.Code == 201) {
		return nil, errors.New(resp.Msg)
	}
	r, err = http.Get(resp.Data.Large)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	robots, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	result := &ImageForSend{}
	result.IllustID = resp.Data.Illust
	result.UserName = resp.Data.User.Name
	result.Image = &robots
	_ = ioutil.WriteFile("./test/test", *result.Image, 0644)
	return result, nil
}
