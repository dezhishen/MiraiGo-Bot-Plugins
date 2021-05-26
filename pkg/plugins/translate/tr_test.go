package translate

import (
	"testing"
)

func TestTr(t *testing.T) {
	tr, err := callHttp("test", "en", "zh")
	print(tr)
	print(err)
}
