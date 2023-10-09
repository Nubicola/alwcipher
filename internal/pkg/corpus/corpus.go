package corpus

import (
	"bufio"
	"strings"
)

type Corpus struct {
	eqs        map[int][]string
	calculator ALWCalculator
}

func NewCorpus(calc ALWCalculator) *Corpus {
	var pC = new(Corpus)
	pC.eqs = make(map[int][]string)
	pC.calculator = calc
	return pC
}

// add this string to the corpus
// this string is assumed to be "clean" and added as a whole string
func (c *Corpus) Add(s string) {
	var val = c.calculator.StringValue(s)
	(c.eqs)[val] = removeDuplicate[string](append((c.eqs)[val], strings.ToUpper(s)))
}

func (c *Corpus) Get(x int) []string {
	return c.eqs[x]
}

func (c *Corpus) Calculate(s string) int {
	// convenience method to return the calculation this Corpus uses
	// also you can just grab c.calculator
	return c.calculator.StringValue(s)
}

func (c *Corpus) Read(r *bufio.Scanner) error {
	for r.Scan() {
		c.Add(r.Text())
	}
	return r.Err()
}
