package main

import (
	"fmt"

	"github.com/pk-rawat/gostr/src"
)

func main() {
	var values map[string]interface{}
	values["a"] = 5
	result := gostr.Evaluate("a * 2", values)
	fmt.Println("Result:", result)
}
