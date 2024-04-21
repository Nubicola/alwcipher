package calculator

import (
	"strings"
	"unicode"
)

// a processor that calculates based on regular EQ can only keep calculating on that basis

type ALWCalculator interface {
	// Iterates over the array and calculates each element according to its method
	Calculate(ss []string) (int, error)
	// the method on one element only
	StringValue(s string) (int, error)
}

type EQBaseCalculator struct {
}

// sums all the word values together
// for BaseCalculator, Calculating over an array of strings is no different
// than calculating one string of all those combined
func (eqbc *EQBaseCalculator) Calculate(ss []string) (int, error) {
	var value = 0
	for _, s := range ss {
		v, err := eqbc.StringValue(s)
		if err != nil {
			value += v
		} else {
			return v, err
		}
	}
	return value, nil
}

func (eqbc *EQBaseCalculator) StringValue(s string) (int, error) {
	var value = 0
	for _, c := range s {
		if !unicode.IsDigit(c) { // digits in a string are not added to the value
			i := int(unicode.ToLower(c))
			if i >= int('a') && i <= int('z') {
				value += (i-int('a'))*19%26 + 1
			}
		}
	}
	return value, nil
}

type EQFirstCalculator struct {
}

func (eqfc *EQFirstCalculator) Calculate(ss []string) (int, error) {
	var sb strings.Builder
	for _, s := range ss {
		sb.WriteByte(s[0])
	}
	var eqbc = new(EQBaseCalculator)
	return eqbc.StringValue(sb.String())
}

func (eqfc *EQFirstCalculator) StringValue(s string) (int, error) {
	var eqbc = new(EQBaseCalculator)
	// first character of the string only
	return eqbc.StringValue(string(s[0]))
}

type EQLastCalculator struct {
}

func (eqlc *EQLastCalculator) Calculate(ss []string) (int, error) {
	var sb strings.Builder
	for _, s := range ss {
		sb.WriteByte(s[len(s)-1])
	}
	var eqbc = new(EQBaseCalculator)
	return eqbc.StringValue(sb.String())
}

func (eqlc *EQLastCalculator) StringValue(s string) (int, error) {
	var eqbc = new(EQBaseCalculator)
	return eqbc.StringValue(string(s[len(s)-1]))
}
