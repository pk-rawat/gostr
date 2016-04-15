package main

import (
	"fmt"

	"github.com/pk-rawat/gostr/src"
)

func main() {
	values := make(map[string]interface{})
	values["a"] = 5
	values["b"] = 3
	values["c"] = 2
	values["e"] = ""
	values["f"] = 1203.62
	fmt.Printf("%v\n", values)
	query := "a * 2"
	result := gostr.Evaluate(query, values)
	fmt.Printf("Query: %s                   Result: %v\n", query, result)
	query = "a + b * c"
	result = gostr.Evaluate(query, values)
	fmt.Printf("Query: %s               Result: %v\n", query, result)
	query = "(a + b) * c"
	result = gostr.Evaluate(query, values)
	fmt.Printf("Query: %s             Result: %v\n\n", query, result)
	query = "a < b"
	result = gostr.Evaluate(query, values)
	fmt.Printf("Query: %s                   Result: %v\n", query, result)
	query = "a =< (b * c)"
	result = gostr.Evaluate(query, values)
	fmt.Printf("Query: %s            Result: %v\n", query, result)
	query = "a <> b"
	result = gostr.Evaluate(query, values)
	fmt.Printf("Query: %s                  Result: %v\n", query, result)
	query = "(a + 1) = (b * c)"
	result = gostr.Evaluate(query, values)
	fmt.Printf("Query: %s       Result: %v\n", query, result)
	query = "(a > b) AND (b < c)"
	result = gostr.Evaluate(query, values)
	fmt.Printf("Query: %s     Result: %v\n", query, result)
	query = "(a > b) OR (b < c)"
	result = gostr.Evaluate(query, values)
	fmt.Printf("Query: %s      Result: %v\n", query, result)
	query = "ISBLANK(d)"
	result = gostr.Evaluate(query, values)
	fmt.Printf("Query: %s              Result: %v\n", query, result)
	query = "NOT(ISBLANK(e))"
	result = gostr.Evaluate(query, values)
	fmt.Printf("Query: %s         Result: %v\n", query, result)
	query = "ROUND(f)"
	result = gostr.Evaluate(query, values)
	fmt.Printf("Query: %s                Result: %v\n", query, result)
}
