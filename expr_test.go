package expr

import (
	"fmt"

	"github.com/rhysd/actionlint"
)

func ExampleEvaluate() {
	expression := "input.foo <= input.bar"

	// Lex & Parse
	lexer := actionlint.NewExprLexer(expression + "}}")
	parser := actionlint.NewExprParser()
	n, perr := parser.Parse(lexer)
	if perr != nil {
		panic(perr)
	}

	// Evaluate expressions
	result, err := Evaluate(n, ContextData{
		"input": ContextData{
			"foo": float64(1),
			"bar": float64(2),
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Println(result.Value)
	// Output: true
}
