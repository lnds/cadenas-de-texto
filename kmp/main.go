package main

import (
	"fmt"

	"kmp/kmp"
)

func main() {
	k := kmp.New("cab")
	fmt.Println("in abacab: ", k.Match("abacab"))
}
