package translate

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/antchfx/htmlquery"
	"github.com/go-basic/uuid"
)

func TestTr(t *testing.T) {
	uri := "http://dict-co.iciba.com/search.php?word=" + "dog"
	resp, err := http.DefaultClient.Get(uri)
	if err != nil {
		return
	}
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	respBodyStr := string(robots)
	root, _ := htmlquery.Parse(strings.NewReader(respBodyStr))
	brs := htmlquery.Find(root, "/html/body/text()")
	for _, row := range brs {
		ele := htmlquery.InnerText(row)
		ele = strings.TrimSpace(ele)
		if ele != "" {
			vals := strings.Split(ele, " ")
			for _, val := range vals {
				fmt.Printf("%v\n", val)
			}

		}

	}
}

func testYouDao() {

	//http://fanyi.youdao.com/openapi.do?keyfrom=<keyfrom>&key=<key>&type=data&doctype=<doctype>&version=1.1&q=要翻译的文本
	uri := "https://openapi.youdao.com/api"

	salt := uuid.New()
	curtime := strconv.FormatInt(time.Now().UTC().Unix(), 10)
	fmt.Printf("%v\n", curtime)
	appkey := "0b5c8081839859d7"
	appsecret := "51Wil4okuw5LIHohef2Zc3FPimaGgpDi"
	q := "test"
	q = truncate(q)
	fmt.Printf("%v\n", q)
	sign := appkey + q + salt + curtime + appsecret //APP_KEY + truncate(q) + salt + curtime + APP_SECRET
	encSign := fmt.Sprintf("%x", sha256.Sum256([]byte(sign)))
	fmt.Printf("%v\n", encSign)

	var dic = map[string]string{
		"q":        q, //待转文字
		"from":     "en",
		"to":       "zh-CHS",
		"appKey":   appkey,
		"salt":     salt,
		"sign":     encSign, //hash
		"signType": "v3",
		"curtime":  curtime, //当前时间戳
	}

	loop := 0
	data := ""
	for k, v := range dic {
		if loop > 0 {
			data += "&"
		}
		data += fmt.Sprintf("%v=%v", k, v)
		loop++
	}
	fmt.Printf("%v\n", data)

	resData, _ := gbkToUtf8([]byte(data))
	resp, err := http.DefaultClient.Post(uri, "application/x-www-form-urlencoded", strings.NewReader(string(resData)))

	if err != nil {
		return
	}
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	respBodyStr := string(robots)
	fmt.Printf("%v\n", respBodyStr)
}
