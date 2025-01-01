package test

import (
	"bufio"
	"bytes"
	"os"
	"strings"
	"testing"

	calculator "github.com/Nubicola/alwcipher/pkg/calculator"
	corpus "github.com/Nubicola/alwcipher/pkg/corpus"
)

func TestAddWholeString(t *testing.T) {
	var calculator = new(calculator.EQBaseCalculator)
	var corpus = corpus.NewCorpusMap(calculator)

	var s = "hello"
	corpus.Add(s)

	var v, _ = calculator.StringValue(s)

	// there is exactly one string in there; check that it's the right one
	if corpus.Get(v)[0] != strings.ToUpper(s) {
		t.Errorf("expected %v but found %v", s, corpus.Get(v)[0])
	}
}

func TestAddStringWithInvalidChars(t *testing.T) {
	var calculator = new(calculator.EQBaseCalculator)
	var corpus = corpus.NewCorpusMap(calculator)

	var s = "hello?!!"
	var p = "hello"
	corpus.Add(s)
	v1, _ := calculator.StringValue(p)
	v2, _ := calculator.StringValue(s)
	if v1 != v2 {
		t.Error("expected value of equivalent strings to be the same, but they're not")
		return
	}
	var xs = corpus.Get(v1)
	if xs[0] != strings.ToUpper(p) {
		t.Error("expected retrieved string to be the same as trimmed string")
	}
}

func TestGetOnce(t *testing.T) {
	var calculator = new(calculator.EQBaseCalculator)
	var corpus = corpus.NewCorpusMap(calculator)

	var s = "hello"

	corpus.Add(s)

	var v, _ = calculator.StringValue(s)
	var xs = corpus.Get(v)
	if xs[0] != strings.ToUpper(s) {
		t.Errorf("expected %v but found %v", s, xs[0])
	}
}

func TestGetMany(t *testing.T) {
	var calculator = new(calculator.EQBaseCalculator)
	var corpus = corpus.NewCorpusMap(calculator)

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
		var v, _ = calculator.StringValue(l)
		ns = corpus.Get(v)
		if ns[0] != strings.ToUpper(l) {
			t.Errorf("expected %v but found %v", l, ns[0])
		}
	}
}

func TestFirst(t *testing.T) {
	var calculator = new(calculator.EQFirstCalculator)
	//var corpus = corpus.NewCorpus(calculator)
	var xs = []string{"TODAY", "ALL", "NEW", "TOYS"} // TANT = 63
	var f, _ = calculator.Calculate(xs)
	if f != 63 {
		t.Errorf("expected 63 but got %v", f)
	}
}

func TestLast(t *testing.T) {
	var calculator = new(calculator.EQLastCalculator)
	//var corpus = corpus.NewCorpus(calculator)
	var xs = []string{"TODAY", "ALL", "NEW", "TOYS"} // YLWS = 25
	var f, _ = calculator.Calculate(xs)
	if f != 25 {
		t.Errorf("expected 25 but got %v", f)
	}
}

func TestCalculate(t *testing.T) {
	var calculator = new(calculator.EQBaseCalculator)
	var corpus = corpus.NewCorpusMap(calculator)

	var s = "hello"
	var p = "foo"
	var g = "gah"
	var k = "kah"
	xs := []string{}
	xs = append(xs, s, p, g, k)
	for _, l := range xs {
		var cs, _ = calculator.StringValue(l)
		var cc, _ = corpus.Calculate(l)
		if cc != cs {
			t.Errorf("expected %v but got %v", cs, cc)
		}
	}
}

func TestReadOneValidCorpusFile(t *testing.T) {
	var calculator = new(calculator.EQBaseCalculator)
	var corpus = corpus.NewCorpusMap(calculator)

	f, err := os.Open("data/corpus_1.md") // values of these are 22
	if err != nil {
		t.Error("failed to open test corpus file")
		return
	}
	defer f.Close()
	scanner := bufio.NewScanner(bufio.NewReader(f))
	scanner.Split(bufio.ScanLines)
	_, err = corpus.ReadFromScanner(scanner)
	if err != nil {
		t.Errorf("failed to read corpus file, %v", err)
		return
	}
	// not sure how to validate that it read it right?
	// could make a scanner here in code?
}

func TestLoad(t *testing.T) {
	var calculator = new(calculator.EQBaseCalculator)
	var corpus = corpus.NewCorpusMap(calculator)
	data := []byte(`{"22":["DRAW","CLAD"]}`)
	// convert byte slice to io.Reader
	reader := bytes.NewReader(data)
	err := corpus.LoadNative(reader)
	if err != nil {
		t.Errorf("expected to read string, received error %v", err)
		return
	}
	var v = corpus.Get(22)
	if !(v[0] == "CLAD" || v[0] == "DRAW") {
		t.Errorf("expected clad or draw, got %v", v[0])
		return
	}
}
