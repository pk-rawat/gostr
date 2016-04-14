package main

import (
	"fmt"
	"os"

	"github.com/pk-rawat/gostr/src"
)

func main() {
	values := make(map[string]interface{})
	values["a"] = 5
	values["b"] = -1
	values["c"] = false
	if len(os.Args) < 2 {
		fmt.Printf("USAGE: %s EXPR", os.Args[0])
		os.Exit(1)
	}
	result := gostr.Evaluate(os.Args[1], values)
	fmt.Println("Result:", result)
}
