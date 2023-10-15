package calculator

import (
	"strings"
	"unicode"
)

// a processor that calculates based on regular EQ can only keep calculating on that basis

type ALWCalculator interface {
	// Iterates over the array and calculates each element according to its method
	Calculate(ss []string) int
	// the method on one element only
	StringValue(s string) int
}

type EQBaseCalculator struct {
}

// sums all the word values together
// for BaseCalculator, Calculating over an array of strings is no different
// than calculating one string of all those combined
func (eqbc *EQBaseCalculator) Calculate(ss []string) int {
	var value = 0
	for _, s := range ss {
		value += eqbc.StringValue(s)
	}
	return value
}

func (eqbc *EQBaseCalculator) StringValue(s string) int {
	var value = 0
	for _, c := range s {
		i := int(unicode.ToLower(c))
		if i >= int('a') {
			value += (i-int('a'))*19%26 + 1
		}
	}
	return value
}

type EQFirstCalculator struct {
}

func (eqfc *EQFirstCalculator) Calculate(ss []string) int {
	var sb strings.Builder
	for _, s := range ss {
		sb.WriteByte(s[0])
	}
	var eqbc = new(EQBaseCalculator)
	return eqbc.StringValue(sb.String())
}

func (eqfc *EQFirstCalculator) StringValue(s string) int {
	var eqbc = new(EQBaseCalculator)
	// first character of the string only
	return eqbc.StringValue(string(s[0]))
}

type EQLastCalculator struct {
}

func (eqlc *EQLastCalculator) Calculate(ss []string) int {
	var sb strings.Builder
	for _, s := range ss {
		sb.WriteByte(s[len(s)-1])
	}
	var eqbc = new(EQBaseCalculator)
	return eqbc.StringValue(sb.String())
}

func (eqlc *EQLastCalculator) StringValue(s string) int {
	var eqbc = new(EQBaseCalculator)
	return eqbc.StringValue(string(s[len(s)-1]))
}
