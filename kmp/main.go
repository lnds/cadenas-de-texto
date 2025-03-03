package main

import (
	"fmt"

	"kmp/kmp"
)

func main() {
	k := kmp.New("ababaa")
	fmt.Println("in abacababaac: ", k.Match("abacababaac"))
	fmt.Println("in abacabab: ", k.Match("abacab"))
	k = kmp.New("cab")
	fmt.Println("cab in abacab", k.Match("abacab"))
}
