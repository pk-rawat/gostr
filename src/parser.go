package gostr

import (
	"strconv"
	"time"
	"unicode"
)

func Includes(array []rune, value rune) bool {
	for _, item := range array {
		if item == value {
			return true
		}
	}
	return false
}

func Parser(query string) Stack {
	tokens := Stack{}
	matchstr := Stack{}
	var curr rune
	var next rune
	buff := make([]rune, 0)
	for i := 0; i < len(query); i++ {
		curr = rune(query[i])
		if i < len(query)-1 {
			next = rune(query[i+1])
		} else {
			next = 10000
		}
		switch {
		case unicode.IsDigit(curr) || curr == '.':
			if !matchstr.IsEmpty() {
				buff = append(buff, curr)
				break
			}
			buff = append(buff, curr)
			if next != 10000 && (unicode.IsDigit(next) || next == '.') {
				continue
			} else {
				if Includes(buff, '_') {
					tokens.Push(Token{Constant, string(buff), Stack{}})
				} else {
					tokens.Push(Token{Number, string(buff), Stack{}})
				}
				buff = nil
			}
		case unicode.IsLetter(curr) || curr == '_':
			if !matchstr.IsEmpty() {
				buff = append(buff, curr)
				break
			}
			if unicode.IsLetter(curr) || curr == '_' {
				buff = append(buff, curr)
			}
			if next != 10000 && (unicode.IsLetter(next) || next == '_' || curr == '_') {
				continue
			} else {
				if string(buff) == "LEN" || string(buff) == "ISBLANK" || string(buff) == "ISNULL" || string(buff) == "NOT" || string(buff) == "DAY" || string(buff) == "MONTH" || string(buff) == "PV" {
					tokens.Push(Token{Function, string(buff), Stack{}})
					buff = nil
				} else if string(buff) == "AND" || string(buff) == "OR" {
					tokens.Push(Token{Operator, string(buff), Stack{}})
					buff = nil
				} else if string(buff) == "TODAY" {
					basedate, _ := time.Parse("01/02/2006", "01/01/1900")
					date := time.Now()
					duration := date.Sub(basedate)
					tokens.Push(Token{Number, strconv.Itoa(int(duration.Hours() / 24)), Stack{}})
					buff = nil
				} else {
					tokens.Push(Token{Constant, string(buff), Stack{}})
					buff = nil
				}
			}
		case IsOperator(curr):
			if !matchstr.IsEmpty() {
				buff = append(buff, curr)
				break
			}
			tokens.Push(Token{Operator, string(curr), Stack{}})
		case IsComparator(curr):
			if !matchstr.IsEmpty() {
				buff = append(buff, curr)
				break
			}
			if (curr == '<' && next == '>') || (curr == '!' && next == '=') {
				tokens.Push(Token{Operator, "<>", Stack{}})
				i++
			} else if curr == '<' && next == '=' {
				tokens.Push(Token{Operator, "<=", Stack{}})
				i++
			} else if curr == '>' && next == '=' {
				tokens.Push(Token{Operator, ">=", Stack{}})
				i++
			} else {
				tokens.Push(Token{Operator, string(curr), Stack{}})
			}
		case curr == '(':
			if !matchstr.IsEmpty() {
				buff = append(buff, curr)
				break
			}
			tokens.Push(Token{LParen, string(curr), Stack{}})
		case curr == ')':
			if !matchstr.IsEmpty() {
				buff = append(buff, curr)
				break
			}
			tokens.Push(Token{RParen, string(curr), Stack{}})
		case curr == '"':
			if matchstr.IsEmpty() {
				matchstr.Push(Token{Quotes, string(curr), Stack{}})
			} else {
				matchstr.Pop()
				tokens.Push(Token{String, string(buff), Stack{}})
				buff = nil
			}
		case curr == ' ':
			if !matchstr.IsEmpty() {
				buff = append(buff, curr)
				break
			}
		case curr == ',':
			if !matchstr.IsEmpty() {
				buff = append(buff, curr)
				break
			}
			tokens.Push(Token{Comma, string(curr), Stack{}})
		}
	}
	return tokens
}

func IsOperator(r rune) bool {
	if r == '+' || r == '-' || r == '/' || r == '*' {
		return true
	}
	return false
}

func IsComparator(r rune) bool {
	if r == '=' || r == '>' || r == '<' {
		return true
	}
	return false
}
