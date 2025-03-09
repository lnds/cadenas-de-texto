package main

import (
	"container/heap"
	"fmt"
	"log"
	"os"

	"tf-idf/tfidf"
)

const threshold = 0.050

func main() {
	if len(os.Args) < 2 {
		return
	}

	var pattern string
	pos := 1
	if os.Args[1] == "-s" {
		pattern = os.Args[2]
		pos = 3
	}

	files := os.Args[pos:]

	searchFile := ""
	if pattern != "" {
		searchFile = "/tmp/qry_tfidf.txt"
		os.WriteFile(searchFile, []byte(pattern), 0644)
		files = append(files, searchFile)
	}

	ds, err := tfidf.New(files)
	if err != nil {
		log.Fatal(err)
	}
	n := min(len(files), 10)
	if searchFile != "" {
		simil := ds.CalcSimils(searchFile)
		top := getTop(simil, n)
		for _, kv := range top {
			if kv.Value > threshold {
				fmt.Println(kv.Key, kv.Value)
			}
		}
		return
		// ifmt.Println("found:")
		// fmt.Printf("%#v\n", simil[searchFile])
	}
	for _, f := range files {
		simil := ds.CalcSimils(f)
		top := getTop(simil, n)
		fmt.Println(f, fmt.Sprintf("%v", top))
	}
}

func getTop(m map[string]float64, n int) []kv {
	h := getHeap(m)
	result := []kv{}
	for range n {
		result = append(result, heap.Pop(h).(kv))
	}
	return result
}

func getHeap(m map[string]float64) *KVHeap {
	h := &KVHeap{}
	heap.Init(h)
	for k, v := range m {
		heap.Push(h, kv{k, v})
	}
	return h
}

type kv struct {
	Key   string
	Value float64
}

type KVHeap []kv

func (h KVHeap) Len() int           { return len(h) }
func (h KVHeap) Less(i, j int) bool { return h[i].Value > h[j].Value }
func (h KVHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *KVHeap) Push(x any) {
	*h = append(*h, x.(kv))
}

func (h *KVHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
