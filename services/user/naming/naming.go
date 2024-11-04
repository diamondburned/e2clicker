package naming

import (
	_ "embed"
	"math/rand/v2"
	"strings"
)

var (
	//go:embed flowers.txt
	flowersTXT string
	//go:embed adjectives.txt
	adjectivesTXT string
)

var (
	flowers    = strings.Split(flowersTXT, "\n")
	adjectives = strings.Split(adjectivesTXT, "\n")
)

// RandomName generates a random name.
func RandomName() string {
	return pick(adjectives) + " " + pick(flowers)
}

func pick[T any](items []T) T {
	return items[rand.N(len(items))]
}
