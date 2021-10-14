package weather

import "testing"

func Test_getWeather(t *testing.T) {
	got, _ := getWeather("岳麓区")
	println(got)
}
