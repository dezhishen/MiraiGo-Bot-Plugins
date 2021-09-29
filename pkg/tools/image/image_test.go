package image

import (
	"log"
	"testing"
)

func TestCreatImage(t *testing.T) {
	err := CreatImage("吉吉吉", "out.png")
	if err != nil {
		log.Fatal(err)
	}
}
