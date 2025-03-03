package kmp

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
	t := -1
	f[0] = -1
	for s := 0; s < m-1; s++ {
		for t >= 0 && k.pattern[s+1] != k.pattern[t+1] {
			t = f[t]
		}
		if k.pattern[s+1] == k.pattern[t+1] {
			t = t + 1
			f[s+1] = t
		} else {
			f[s+1] = -1
		}
	}
	k.fail = f
}

func (k *KMP) Match(text string) int {
	n := len(text)
	m := len(k.pattern)
	s := -1
	for i := 0; i < n; i++ {
		for s > 0 && text[i] != k.pattern[s+1] {
			s = k.fail[s]
		}
		if text[i] == k.pattern[s+1] {
			s = s + 1
		}
		if s == m-1 {
			return i - s
		}
	}
	return -1
}
