package lpl

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func TestLpl(t *testing.T) {
	recentGamesUrl := "https://lpl.qq.com/web201612/data/LOL_MATCH2_MATCH_HOMEPAGE_BMATCH_LIST_148.js"
	resp, err := http.DefaultClient.Get(recentGamesUrl)

	if err != nil {
		return
	}
	robots, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return
	}
	respBodyStr := string(robots)
	if respBodyStr == "" {
		return
	}
	var mBox MatchBox
	err = json.Unmarshal(robots, &mBox)
	if err != nil {
		return
	}

	nowLocal := time.Now().Local()

	for _, game := range mBox.GameInfo {

		matchDate, err := time.Parse("2006-01-02 15:04:05", game.MatchDate)
		if err != nil {
			continue
		}
		if matchDate.Before(nowLocal) {
			continue
		}
		if matchDate.Day() == nowLocal.Day() {
			todayInfo := createMatchInfo(game, true)
			fmt.Printf("%s", todayInfo)
		}
		if matchDate.Day() == nowLocal.Day()+1 {
			tomorrowInfo := createMatchInfo(game, false)
			fmt.Printf("%s", tomorrowInfo)
		}

	}

}
