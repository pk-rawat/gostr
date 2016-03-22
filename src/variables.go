package gostr

import (
	"math"
	"strconv"
)

type Token struct {
	Type   TokenType
	Lexeme string
	Value  Stack
}

type TokenType byte

const (
	Number TokenType = iota
	Boolean
	Column
	Function
	Operator
	Comparator
	LParen
	RParen
	Quotes
	Constant
	String
	Comma
)

type Stack struct {
	Values []Token
}

var oprData = map[string]struct {
	prec  int
	rAsoc bool
	fx    func(x, y float64) interface{}
}{
	"^":   {4, true, func(x, y float64) interface{} { return math.Pow(x, y) }},
	"*":   {3, false, func(x, y float64) interface{} { return x * y }},
	"/":   {3, false, func(x, y float64) interface{} { return x / y }},
	"+":   {2, false, func(x, y float64) interface{} { return x + y }},
	"-":   {2, false, func(x, y float64) interface{} { return x - y }},
	"=":   {2, false, func(x, y float64) interface{} { return true }},
	"<":   {2, false, func(x, y float64) interface{} { return true }},
	">":   {2, false, func(x, y float64) interface{} { return true }},
	"<=":  {2, false, func(x, y float64) interface{} { return true }},
	">=":  {2, false, func(x, y float64) interface{} { return true }},
	"!>":  {2, false, func(x, y float64) interface{} { return true }},
	"<>":  {2, false, func(x, y float64) interface{} { return true }},
	"AND": {0, false, func(x, y float64) interface{} { return true }},
	"OR":  {0, false, func(x, y float64) interface{} { return true }},
}

var cmpData = map[string]struct {
	fx func(x, y bool) interface{}
}{
	"AND": {func(x, y bool) interface{} { return x && y }},
	"OR":  {func(x, y bool) interface{} { return x || y }},
}

var chkData = map[string]struct {
	fx func(x, y interface{}) interface{}
}{
	"=": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return x == y
		} else {
			a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
			b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
			return a == b
		}
	}},
	"!=": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return x != y
		} else {
			a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
			b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
			return a != b
		}
	}},
	"<>": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return x != y
		} else {
			a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
			b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
			return a != b
		}
	}},
	">": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return false
		} else {
			a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
			b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
			return a > b
		}
	}},
	"<": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return false
		} else {
			a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
			b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
			return a < b
		}
	}},
	">=": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return false
		} else {
			a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
			b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
			return a >= b
		}
	}},
	"<=": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return false
		} else {
			a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
			b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
			return a <= b
		}
	}},
}

func (self *Stack) Push(i ...Token) {
	self.Values = append(self.Values, i...)
}

func (self *Stack) Pop() Token {
	if len(self.Values) == 0 {
		return Token{}
	}
	token := self.Values[len(self.Values)-1]
	self.Values = self.Values[:len(self.Values)-1]
	return token
}

func (self *Stack) Peek() Token {
	if len(self.Values) == 0 {
		return Token{}
	}
	return self.Values[len(self.Values)-1]
}

func (self *Stack) EmptyInto(s *Stack) {
	if !self.IsEmpty() {
		for i := self.Length() - 1; i >= 0; i-- {
			s.Push(self.Pop())
		}
	}
}

func (self *Stack) IsEmpty() bool {
	return len(self.Values) == 0
}

func (self *Stack) Length() int {
	return len(self.Values)
}
