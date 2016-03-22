package gostr

import (
	"fmt"
	"math"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func Evaluate(query string, values map[string]interface{}) interface{} {
	tokens := Parser(query)
	rpn := ToPostfix(tokens)
	outs := SolvePostfix(rpn, values)
	return outs
}

func ToPostfix(tokens Stack) Stack {
	ops := Stack{}
	output := Stack{}
	fun := Stack{}
	for _, v := range tokens.Values {
		switch v.Type {
		case Operator:
			if fun.Length() > 0 {
				tok := fun.Pop()
				tok.Value.Push(v)
				fun.Push(tok)
			} else {
				for !ops.IsEmpty() {
					val := v.Lexeme
					top := ops.Peek().Lexeme
					if (oprData[val].prec <= oprData[top].prec && oprData[val].rAsoc == false) || (oprData[val].prec < oprData[top].prec && oprData[val].rAsoc == true) {
						output.Push(ops.Pop())
						continue
					}
					break
				}
				ops.Push(v)
			}
		case LParen:
			if fun.Length() > 0 {
				tok := fun.Pop()
				tok.Value.Push(v)
				fun.Push(tok)
			} else {
				ops.Push(v)
			}
		case RParen:
			if fun.Length() > 0 {
				tok := fun.Pop()
				closeparen := 0
				for i := tok.Value.Length() - 1; i >= 0; i-- {
					if tok.Value.Values[i].Type == RParen {
						closeparen += 1
					}
					if tok.Value.Values[i].Type != LParen {
						continue
					} else if i > 0 && closeparen == 0 {
						tok.Value.Push(v)
						fun.Push(tok)
						break
					} else {
						tok.Value.Push(v)
						output.Push(tok)
						break
					}
				}
			} else {
				for i := ops.Length() - 1; i >= 0; i-- {
					if ops.Values[i].Type != LParen {
						output.Push(ops.Pop())
						continue
					} else {
						ops.Pop()
						break
					}
				}
			}
		case Function:
			if fun.Length() > 0 {
				tok := fun.Pop()
				tok.Value.Push(v)
				fun.Push(tok)
			} else {
				fun.Push(v)
			}
		default:
			if fun.Length() > 0 {
				tok := fun.Pop()
				tok.Value.Push(v)
				fun.Push(tok)
			} else {
				output.Push(v)
			}
		}
	}
	ops.EmptyInto(&output)
	return output
}

// SolvePostfix evaluates and returns the answer of the expression converted to postfix
func SolvePostfix(tokens Stack, vars map[string]interface{}) interface{} {
	stack := Stack{}
	for _, v := range tokens.Values {
		switch v.Type {
		case Number:
			stack.Push(v)
		case Function:
			stack = SolveFunction(v, vars, stack)
		case String:
			stack.Push(Token{String, v.Lexeme, Stack{}})
		case Constant:
			if v.Lexeme == "true" {
				stack.Push(Token{Boolean, strconv.FormatBool(true), Stack{}})
			} else {
				data := vars[v.Lexeme]
				if data == nil {
					stack.Push(Token{Number, "", Stack{}})
				} else {
					stack = PushStringToStack(data, stack)
				}
			}
		case Operator:
			if v.Lexeme == "=" || v.Lexeme == "<>" || v.Lexeme == ">" || v.Lexeme == "<" || v.Lexeme == ">=" || v.Lexeme == "<=" {
				var x, y interface{}
				y = stack.Pop().Lexeme
				x = stack.Pop().Lexeme
				fx := chkData[v.Lexeme].fx
				result := fx(x, y)
				stack.Push(Token{Boolean, strconv.FormatBool(result.(bool)), Stack{}})
			} else if v.Lexeme == "AND" || v.Lexeme == "OR" {
				var x, y bool
				y, _ = strconv.ParseBool(stack.Pop().Lexeme)
				x, _ = strconv.ParseBool(stack.Pop().Lexeme)
				fx := cmpData[v.Lexeme].fx
				result := fx(x, y)
				stack.Push(Token{Boolean, strconv.FormatBool(result.(bool)), Stack{}})
			} else {
				f := oprData[v.Lexeme].fx
				var x, y float64
				y, _ = strconv.ParseFloat(stack.Pop().Lexeme, 64)
				x, _ = strconv.ParseFloat(stack.Pop().Lexeme, 64)
				result := f(x, y)
				stack.Push(Token{Number, strconv.FormatFloat(result.(float64), 'f', -1, 64), Stack{}})
			}
		}
	}
	out := stack.Values[0].Lexeme
	return out
}

// SolveFunction returns the answer of a function found within an expression
func SolveFunction(v Token, vars map[string]interface{}, stack Stack) Stack {
	var value interface{}
	fun_tokens := v.Value
	if fun_tokens.Length() > 1 && v.Lexeme != "PV" {
		toks := ToPostfix(fun_tokens)
		if toks.Length() > 0 {
			value = SolvePostfix(toks, vars)
		}
	}
	if v.Lexeme == "LEN" {
		stack.Push(Token{Number, strconv.Itoa(len(value.(string))), Stack{}})
	} else if v.Lexeme == "ISBLANK" || v.Lexeme == "ISNULL" {
		val := false
		if len(strings.TrimSpace(value.(string))) == 0 {
			val = true
		}
		stack.Push(Token{Boolean, strconv.FormatBool(val), Stack{}})
	} else if v.Lexeme == "MONTH" {
		parsefloat, ok := strconv.ParseFloat(value.(string), 64)
		if ok != nil {
			fmt.Println("Error:", ok)
		}
		days := int(parsefloat)
		months := days * 12 / 365
		stack.Push(Token{Number, strconv.Itoa(months), Stack{}})
	} else if v.Lexeme == "DAY" {
		parsefloat, ok := strconv.ParseFloat(value.(string), 64)
		if ok != nil {
			fmt.Println("Error:", ok)
		}
		days := int(parsefloat)
		basedate, _ := time.Parse("01/02/2006", "01/01/1900")
		date := basedate.AddDate(0, 0, days)
		stack.Push(Token{Number, strconv.Itoa(date.Day()), Stack{}})
	} else if v.Lexeme == "NOT" {
		if value == "true" {
			stack.Push(Token{Boolean, strconv.FormatBool(false), Stack{}})
		} else {
			stack.Push(Token{Boolean, strconv.FormatBool(true), Stack{}})
		}
	} else if v.Lexeme == "PV" {
		var rate, time, pmt float64
		var ok error
		toks := ToPostfix(v.Value)
		for index, item := range toks.Values {
			stk := Stack{}
			stk.Push(item)
			switch index {
			case 0:
				val := SolvePostfix(stk, vars)
				rate, ok = strconv.ParseFloat(val.(string), 64)
				if ok != nil {
					fmt.Println("Error:", ok)
				}
			case 2:
				val := SolvePostfix(stk, vars)
				time, ok = strconv.ParseFloat(val.(string), 64)
				if ok != nil {
					fmt.Println("Error:", ok)
				}
			case 4:
				val := SolvePostfix(stk, vars)
				pmt, ok = strconv.ParseFloat(val.(string), 64)
				if ok != nil {
					fmt.Println("Error:", ok)
				}
			}
		}
		v := math.Pow((1 + rate), -time)
		pv := pmt * ((1 - v) / rate)
		str := strconv.FormatFloat(pv, 'f', 2, 64)
		stack.Push(Token{Number, str, Stack{}})
	}
	return stack
}

func PushStringToStack(data interface{}, stack Stack) Stack {
	switch reflect.TypeOf(data).Kind() {
	case reflect.Int:
		stack.Push(Token{Number, strconv.Itoa(data.(int)), Stack{}})
	case reflect.Float64:
		stack.Push(Token{Number, strconv.FormatFloat(data.(float64), 'f', 2, 64), Stack{}})
	case reflect.String:
		stack.Push(Token{Number, data.(string), Stack{}})
	case reflect.Bool:
		stack.Push(Token{Number, strconv.FormatBool(data.(bool)), Stack{}})
	}
	return stack
}