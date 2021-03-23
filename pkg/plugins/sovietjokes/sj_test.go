package sovietjokes

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestSj(t *testing.T) {

	jokes := getJokes()
	rand.Seed(time.Now().UnixNano())
	v := rand.Intn(len(jokes) - 1)
	randJoke := jokes[v]
	fmt.Printf("%s\n", randJoke.Title)
	fmt.Printf("%s", randJoke.Content)
}
