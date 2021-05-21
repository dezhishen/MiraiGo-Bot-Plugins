package translate

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"testing"
)

func TestTr(t *testing.T) {
	uri := "http://api.fanyi.baidu.com/api/trans/vip/translate?"

	qbt, _ := gbkToUtf8([]byte("徐思佳是我儿"))
	q := string(qbt)
	q = "测试"
	appid := "20210521000836370"
	from := "auto"
	to := "auto"
	salt := strconv.Itoa(rand.Intn(100000))
	data := appid + q + salt + "f6ctAeXzVZNaknrgqiKs"
	w := md5.New()
	w.Write([]byte(data))
	sign := hex.EncodeToString(w.Sum(nil))

	uri += fmt.Sprintf("q=%v", q)
	uri += fmt.Sprintf("&from=%v", from)
	uri += fmt.Sprintf("&to=%v", to)
	uri += fmt.Sprintf("&appid=%v", appid)
	uri += fmt.Sprintf("&salt=%v", salt)
	uri += fmt.Sprintf("&sign=%v", sign)

	fmt.Printf("%v\n", uri)
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
	fmt.Printf("%v\n", respBodyStr)
}
