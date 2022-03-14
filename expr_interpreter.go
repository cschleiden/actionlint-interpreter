package main

import (
	"strings"

	"github.com/rhysd/actionlint"
)

func Evaluate(n actionlint.ExprNode) (interface{}, error) {
	switch tn := n.(type) {
	case *actionlint.IntNode:
		return tn.Value, nil

	case *actionlint.FloatNode:
		return tn.Value, nil

	case *actionlint.StringNode:
		return tn.Value, nil

	case *actionlint.BoolNode:
		return tn.Value, nil

	case *actionlint.NotOpNode:
		r, err := Evaluate(tn.Operand)
		if err != nil {
			return nil, err
		}

		// TODO: Coerce
		b := r.(bool)
		return !b, nil

	case *actionlint.FuncCallNode:
		args := make([]interface{}, len(tn.Args))
		for i, arg := range tn.Args {
			a, err := Evaluate(arg)
			if err != nil {
				return nil, err
			}

			args[i] = a
		}

		return fcall(tn.Callee, args)

	case *actionlint.CompareOpNode:
		left, err := Evaluate(tn.Left)
		if err != nil {
			return nil, err
		}
		right, err := Evaluate(tn.Right)
		if err != nil {
			return nil, err
		}

		// TODO: Support coercion
		switch tn.Kind {
		case actionlint.CompareOpNodeKindEq:
			return left == right, nil

			// TODO: Support other operators
		}
	}

	panic("unknown node")
}

func fcall(name string, args []interface{}) (interface{}, error) {
	// Expression function names are case-insensitive.
	switch strings.ToLower(name) {
	case "startswith":
		// TODO: Verify args, count and type
		return strings.HasPrefix(args[0].(string), args[1].(string)), nil
	}

	panic("unknown function")
}

// func coerce(v interface{}) interface{} {
// 	switch x := v.(type) {
// 		case int:
// 	}
// }
