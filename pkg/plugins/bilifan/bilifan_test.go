package bilifan

import (
	"testing"
	"fmt"
)

func TestR(t *testing.T) {
	r, err := getBiliFan("1")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(r)
	t.Log(r)
}
