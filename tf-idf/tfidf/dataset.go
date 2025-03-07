package tfidf

import (
	"math"
	"os"
	"regexp"
	"strings"
)

type (
	Term       = string
	File       = string
	BagOfWords map[Term]int
	Documents  map[File]BagOfWords
	Vocabulary map[Term]float64
	Vector     map[int]float64
	Vectors    map[File]Vector
	TermIndex  map[Term]int
)

type DataSet struct {
	files   []File
	vectors Vectors
	sims    map[string]float64
}

func New(files []File) (*DataSet, error) {
	documents := Documents{}
	for _, f := range files {
		doc, err := loadFile(f)
		if err != nil {
			return nil, err
		}
		documents[f] = doc
	}

	vocabulary, terms := extractTerms(documents)
	n := float64(len(terms))
	idfs := Vocabulary{}
	for t, f := range vocabulary {
		idfs[t] = math.Log(n / f)
	}
	termIndex := TermIndex{}
	for i, t := range terms {
		termIndex[t] = i
	}

	vectors := calcVectors(documents, idfs, termIndex)
	return &DataSet{
		files:   files,
		vectors: vectors,
		sims:    map[string]float64{},
	}, nil
}

func (ds *DataSet) CalcSimils(f File) map[File]float64 {
	result := map[File]float64{}
	for _, g := range ds.files {
		if g != f {
			result[g] = ds.calcSimil(f, g)
		}
	}
	return result
}

func (ds *DataSet) calcSimil(f, g File) float64 {
	key1 := f + ":" + g
	key2 := g + ":" + f
	v, ok1 := ds.sims[key1]
	_, ok2 := ds.sims[key2]
	if ok1 || ok2 {
		return v
	}
	if f == g {
		return 1.0
	}
	result := calcSimilitude(ds.vectors[f], ds.vectors[g])
	ds.sims[key1] = result
	ds.sims[key2] = result
	return result
}

func calcSimilitude(vec1, vec2 Vector) float64 {
	sum := 0.0
	for k, v1 := range vec1 {
		if v2, ok := vec2[k]; ok {
			sum += v1 * v2
		}
	}
	return sum
}

func loadFile(file File) (BagOfWords, error) {
	content, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}
	terms := contentToTerms(string(content))
	bag := BagOfWords{}
	for _, term := range terms {
		bag[term]++
	}
	return bag, nil
}

func contentToTerms(doc string) []string {
	re := regexp.MustCompile("[\\w|ñ|á|é|í|ó|ú]+")
	return re.FindAllString(strings.ToLower(doc), -1)
}

func extractTerms(docs Documents) (Vocabulary, []Term) {
	vocabulary := Vocabulary{}
	for _, doc := range docs {
		for w, f := range doc {
			vocabulary[w] += float64(f)
		}
	}

	terms := []Term{}
	for w := range vocabulary {
		terms = append(terms, w)
	}
	return vocabulary, terms
}

func calcVectors(docs Documents, idfs Vocabulary, termIndex TermIndex) Vectors {
	vectors := Vectors{}
	for file, doc := range docs {
		vectors[file] = normalize(doc, idfs, termIndex)
	}
	return vectors
}

func normalize(doc BagOfWords, idfs Vocabulary, termIndex TermIndex) Vector {
	vector := Vector{}
	sum := 0.0
	for term, freq := range doc {
		tfi := 1 + math.Log(float64(freq))
		idf := idfs[term]
		f := tfi * idf
		sum += f * f
		vector[termIndex[term]] = f
	}
	for i, x := range vector {
		vector[i] = x / math.Sqrt(sum)
	}
	return vector
}
