package haimage

import (
	"fmt"
	"testing"
)

func TestGetHaImage(t *testing.T) {
	robots, _ := getHaImage()
	fmt.Print(robots)
}
