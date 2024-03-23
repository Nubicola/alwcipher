package corpus

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"

	calculator "github.com/Nubicola/alwcipher/pkg/calculator"
)

type Corpus struct {
	Eqs        map[int][]string // exported so it can be marshalled
	calculator calculator.ALWCalculator
}

func NewCorpus(calc calculator.ALWCalculator) *Corpus {
	var pC = new(Corpus)
	pC.Eqs = make(map[int][]string)
	pC.calculator = calc
	return pC
}

// add this string to the corpus
// this string is assumed to be "clean" and added as a whole string
func (c *Corpus) Add(s string) error {
	v := strings.ToUpper(strings.TrimRight(s, ":?!;.,'"))
	var val, err = c.calculator.StringValue(v)
	if err == nil {
		(c.Eqs)[val] = removeDuplicate[string](append((c.Eqs)[val], v))
	}
	return err
}

func (c *Corpus) Get(x int) []string {
	return c.Eqs[x]
}

func (c *Corpus) Calculate(s string) (int, error) {
	// convenience method to return the calculation this Corpus uses
	// also you can just grab c.calculator
	return c.calculator.StringValue(s)
}

func (c *Corpus) ReadFromScanner(r *bufio.Scanner) (int, error) {
	i := 0
	for r.Scan() {
		c.Add(r.Text())
		i += 1
	}
	return i, r.Err()
}

func (c *Corpus) SaveNative(w io.Writer) error {
	b, err := json.Marshal(c.Eqs)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	if err != nil {
		return err
	}
	return nil
}

func (c *Corpus) LoadNative(r io.Reader) error {
	b, err := io.ReadAll(r)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &c.Eqs)
}

// helper functions
// from stack overflow: https://stackoverflow.com/questions/66643946/how-to-remove-duplicates-strings-or-int-from-slice-in-go
func removeDuplicate[T string | int](sliceList []T) []T {
	allKeys := make(map[T]bool)
	list := []T{}
	for _, item := range sliceList {
		if _, value := allKeys[item]; !value {
			allKeys[item] = true
			list = append(list, item)
		}
	}
	return list
}
