package pixiv

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
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
	Illust    string      `json:"illust"`
	Title     string      `json:"title"`
	Large     string      `json:"large"`
	User      *user       `json:"user"`
	Originals []*original `json:"originals"`
}

type user struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type original struct {
	Url string `json:"url"`
}
type ImageForSend struct {
	Title      string
	SourceFrom string
	Images     *[]*[]byte
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
	result := &ImageForSend{}
	result.IllustID = resp.Data.Illust
	result.UserName = resp.Data.User.Name
	result.Images = getImages(&resp)
	return result, nil
}

func getImages(resp *Resp) *[]*[]byte {
	var images []*[]byte
	for index, v := range resp.Data.Originals {
		image, err := getImage(resp.Data.Illust, index, v.Url)
		if err != nil {
			continue
		}
		images = append(images, image)
	}
	return &images
}

func getImage(id string, index int, url string) (*[]byte, error) {
	path := fmt.Sprintf("./pixiv/%v", id)
	fileName := getFileName(url)
	filePath := fmt.Sprintf("%v/%v", path, fileName)
	ok, _ := pathExists(filePath)
	if ok {
		file, err := os.Open(filePath)
		if err != nil {
			panic(err)
		}
		defer file.Close()
		content, err := ioutil.ReadAll(file)
		return &content, err

	}

	exists, _ := pathExists(path)
	if !exists {
		os.Mkdir(path, 0777)
	}
	r, err := http.DefaultClient.Get(url)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	robots, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	ioutil.WriteFile(filePath, robots, 0644)
	return &robots, nil
}

func getFileName(url string) string {
	i := strings.LastIndex(url, "/")
	return url[i+1:]
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func init() {
	exists, _ := pathExists("./pixiv")
	if !exists {
		os.Mkdir("./pixiv", 0777)
	}
}
