package gostr

import (
	"math"
	"strconv"
	"strings"
)

// Token is a object struct for used token in stack
type Token struct {
	Type   TokenType
	Lexeme string
	Value  Stack
}

// TokenType is used as byte
type TokenType byte

// const declare used constant for tokens
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

// Stack is a collection of Tokens
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
	"and": {0, false, func(x, y float64) interface{} { return true }},
	"OR":  {0, false, func(x, y float64) interface{} { return true }},
	"or":  {0, false, func(x, y float64) interface{} { return true }},
}

var cmpData = map[string]struct {
	fx func(x, y bool) interface{}
}{
	"AND": {func(x, y bool) interface{} { return x && y }},
	"and": {func(x, y bool) interface{} { return x && y }},
	"OR":  {func(x, y bool) interface{} { return x || y }},
	"or":  {func(x, y bool) interface{} { return x || y }},
}

var chkData = map[string]struct {
	fx func(x, y interface{}) interface{}
}{
	"=": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return strings.EqualFold(x.(string), y.(string))
		}
		a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
		b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
		return a == b
	}},
	"!=": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return !strings.EqualFold(x.(string), y.(string))
		}
		a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
		b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
		return a != b
	}},
	"<>": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return !strings.EqualFold(x.(string), y.(string))
		}
		a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
		b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
		return a != b
	}},
	">": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return false
		}
		a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
		b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
		return a > b
	}},
	"<": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return false
		}
		a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
		b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
		return a < b
	}},
	">=": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return false
		}
		a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
		b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
		return a >= b
	}},
	"<=": {func(x, y interface{}) interface{} {
		a, err := strconv.ParseFloat(x.(string), 64)
		b, err := strconv.ParseFloat(y.(string), 64)
		if err != nil {
			return false
		}
		a, err = strconv.ParseFloat(strconv.FormatFloat(a, 'f', 2, 64), 64)
		b, err = strconv.ParseFloat(strconv.FormatFloat(b, 'f', 2, 64), 64)
		return a <= b
	}},
}

// Push insert tokens to stack
func (stack *Stack) Push(i ...Token) {
	stack.Values = append(stack.Values, i...)
}

// Pop remove last token from stack
func (stack *Stack) Pop() Token {
	if len(stack.Values) == 0 {
		return Token{}
	}
	token := stack.Values[len(stack.Values)-1]
	stack.Values = stack.Values[:len(stack.Values)-1]
	return token
}

// Peek return last token of stack
func (stack *Stack) Peek() Token {
	if len(stack.Values) == 0 {
		return Token{}
	}
	return stack.Values[len(stack.Values)-1]
}

// EmptyInto move tokens to another stack
func (stack *Stack) EmptyInto(s *Stack) {
	if !stack.IsEmpty() {
		for i := stack.Length() - 1; i >= 0; i-- {
			s.Push(stack.Pop())
		}
	}
}

// IsEmpty check for empty stack
func (stack *Stack) IsEmpty() bool {
	return len(stack.Values) == 0
}

// Length return count of tokens
func (stack *Stack) Length() int {
	return len(stack.Values)
}
