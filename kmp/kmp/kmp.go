package kmp

import "fmt"

type KMP struct {
	pattern string
	fail    []int
}

func New(pattern string) *KMP {
	kmp := &KMP{
		pattern: pattern,
	}
	kmp.buildFailureTable()
	return kmp
}

func (k *KMP) buildFailureTable() {
	m := len(k.pattern)
	f := make([]int, m)
	t := 0
	f[0] = -1
	for s := 1; s < m; s++ {
		if k.pattern[s] == k.pattern[t] {
			f[s] = f[t]
		} else {
			f[s] = t
			for t >= 0 && k.pattern[s] != k.pattern[t] {
				t = f[t]
			}
		}
		t++
	}
	k.fail = f
	fmt.Println(k.fail)
}

func (k *KMP) Match(text string) int {
	n := len(text)
	m := len(k.pattern)
	s := 0
	i := 0
	for i < n {
		if k.pattern[s] == text[i] {
			s++
			i++
			if s == m {
				return i - m
			}
		} else {
			s = k.fail[s]
			if s < 0 {
				i++
				s++
			}
		}
	}
	return -1
}
