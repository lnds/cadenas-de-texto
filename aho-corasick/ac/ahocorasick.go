package ac

import (
	"fmt"
	"log"
)

const ALPHABET_SIZE = 26

type AC struct {
	keywords []string
	g        [][]int
	fail     []int
	final    map[int]string
}

func New(keywords []string) *AC {
	ac := &AC{
		keywords: keywords,
		g:        [][]int{},
		final:    map[int]string{},
	}
	s := ac.buildTrie()
	ac.buildFail(s)
	return ac
}

func (ac *AC) addState(fill int) int {
	state := make([]int, ALPHABET_SIZE, ALPHABET_SIZE)
	for i := range state {
		state[i] = fill
	}
	ac.g = append(ac.g, state)
	return len(ac.g) - 1
}

func (ac *AC) buildTrie() int {
	maxS := ac.addState(0)
	for _, word := range ac.keywords {
		s := 0
		for _, c := range []byte(word) {
			if c < 'a' || c > 'z' {
				log.Fatal("alphabet is 'a' to 'z'")
			}
			c = c - 'a'
			if ac.g[s][c] <= 0 {
				maxS = ac.addState(-1)
				ac.g[s][c] = maxS
				s = maxS
			} else {
				s = ac.g[s][c]
			}
		}
		ac.final[s] = word
	}
	ac.dumpTrie()
	return maxS + 1
}

func (ac *AC) dumpTrie() {
	for i, s := range ac.g {
		fmt.Printf("%2d: ", i)
		endNode := true
		for c, ns := range s {
			if ns > 0 {
				fmt.Printf("%c -> %d ", c+'a', ns)
				endNode = false
			}
		}
		if endNode {
			fmt.Print("$")
		}
		fmt.Println("")
	}
	fmt.Println(ac.final)
}

func (ac *AC) buildFail(states int) {
	ac.fail = make([]int, states)
	for d := 1; d < states; d++ {
		for a, si := range ac.g[d] {
			if si < 0 {
				continue
			}
			s := ac.fail[d]
			for ac.g[s][a] < 0 {
				s = ac.fail[s]
			}
			ac.fail[si] = ac.g[s][a]
		}
	}
	ac.dumpFail()
}

func (ac *AC) dumpFail() {
	fmt.Println(ac.fail)
}

func (ac *AC) Match(text string) map[string][]int {
	s := 0
	result := map[string][]int{}
	for pos, a := range []byte(text) {
		if a < 'a' || a > 'z' {
			continue
		}
		a = a - 'a'
		for ac.g[s][a] < 0 {
			s = ac.fail[s]
		}
		s = ac.g[s][a]
		if w, ok := ac.final[s]; ok {
			result[w] = append(result[w], pos-len(w)+1)
		}
	}
	return result
}
