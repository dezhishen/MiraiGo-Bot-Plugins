package weather

import (
	"io/ioutil"
	"testing"
)

func Test_getWeather(t *testing.T) {
	got, _ := getWeather("岳麓区")
	_ = ioutil.WriteFile("weather.png", got, 0644)
}
