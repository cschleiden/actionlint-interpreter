package expr

import (
	"fmt"

	"github.com/rhysd/actionlint"
)

func ExampleEvaluate() {
	expression := "1 <= 2"

	// Lex & Parse
	lexer := actionlint.NewExprLexer(expression + "}}")
	parser := actionlint.NewExprParser()
	n, perr := parser.Parse(lexer)
	if perr != nil {
		panic(perr)
	}

	// Evaluate expressions
	result, err := Evaluate(n, map[string]interface{}{})
	if err != nil {
		panic(err)
	}

	fmt.Println(result.Value)
	// Output: true
}
