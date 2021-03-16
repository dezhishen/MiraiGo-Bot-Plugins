package calendar

import (
	"testing"
	"time"
)

func TestR(t *testing.T) {
	r, err := getDate(time.Now())
	if err != nil {
		t.Error(err)
	}
	t.Log(r)
}
