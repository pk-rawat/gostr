Gostr
=======

DESCRIPTION
-----------

Gostr is a evaluator for a mathematical and logical expressions that allows run-time binding of values to variables referenced in the formulas.

EXAMPLE
-------

This is probably simplest to illustrate in code:

```go
values := map[string]interface{}
values["a"] = 5
result := gostr.Evaluate("a * 2", values)
fmt.Println("Result:", result)
#=> Result: 10

result := gostr.Evaluate("a + 3 * 2", values)
fmt.Println("Result:", result)
#=> Result: 11

result := gostr.Evaluate("(a + 3) * 2", values)
fmt.Println("Result:", result)
#=> Result: 16
```

BUILT-IN OPERATORS AND FUNCTIONS
---------------------------------

Math: `+ - * /`

Logic: `< > <= >= <> != = AND OR`

Functions: `ISBLANK LENGTH NOT`
