package main

import (
	"fmt"

	"aho-corasick/ac"
)

const lyrics = `
he's a real nowhere man
sitting in his nowhere land
making all his nowhere plans for nobody
`

func main() {
	ac := ac.New([]string{"he", "she", "his", "hers"})
	fmt.Println(ac.Match(lyrics))
}
