package image

import (
	"log"
	"testing"
)

func TestCreatImage(t *testing.T) {
	err := CreatImage("☆★", "out.png")
	if err != nil {
		log.Fatal(err)
	}
}
