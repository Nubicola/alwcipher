package corpus

import (
	"bufio"
	"os"
	"strings"
	"testing"
	"calculator"
)

func TestAddWholeString(t *testing.T) {
	var calculator = new(EQBaseCalculator)
	var corpus = NewCorpus(calculator)

	var s = "hello"
	corpus.Add(s)

	var v = calculator.StringValue(s)

	// there is exactly one string in there; check that it's the right one
	if corpus.eqs[v][0] != strings.ToUpper(s) {
		t.Errorf("expected %v but found %v", s, corpus.eqs[v][0])
	}
}

func TestGetOnce(t *testing.T) {
	var calculator = new(EQBaseCalculator)
	corpus := NewCorpus(calculator)
	var s = "hello"

	corpus.Add(s)

	var v = calculator.StringValue(s)
	var xs = corpus.Get(v)
	if xs[0] != strings.ToUpper(s) {
		t.Errorf("expected %v but found %v", s, xs[0])
	}
}

func TestGetMany(t *testing.T) {
	var calculator = new(EQBaseCalculator)
	corpus := NewCorpus(calculator)
	var s = "hello"
	var p = "foo"
	var g = "gah"
	var k = "kah"
	xs := []string{}
	xs = append(xs, s, p, g, k)

	for _, l := range xs {
		corpus.Add(l)
	}

	var ns = make([]string, 1)
	for _, l := range xs {
		var v = calculator.StringValue(l)
		ns = corpus.Get(v)
		if ns[0] != strings.ToUpper(l) {
			t.Errorf("expected %v but found %v", l, ns[0])
		}
	}
}

func TestCalculate(t *testing.T) {
	var calculator = new(EQBaseCalculator)
	corpus := NewCorpus(calculator)
	var s = "hello"
	var p = "foo"
	var g = "gah"
	var k = "kah"
	xs := []string{}
	xs = append(xs, s, p, g, k)
	for _, l := range xs {
		var cs = calculator.StringValue(l)
		var cc = corpus.Calculate(l)
		if cc != cs {
			t.Errorf("expected %v but got %v", cs, cc)
		}
	}
}

func TestReadOneValidCorpusFile(t *testing.T) {
	var calculator = new(EQBaseCalculator)
	corpus := NewCorpus(calculator)
	f, err := os.Open("test/corpus_1.md")
	if err != nil {
		t.Error("failed to open test corpus file")
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(bufio.NewReader(f))
	scanner.Split(bufio.ScanLines)
	err = corpus.Read(scanner)
	if err != nil {
		t.Errorf("failed to read corpus file, %v", err)
		return
	}
	if len(corpus.eqs) == 0 {
		t.Error("read file but corpus empty")
	}
}
