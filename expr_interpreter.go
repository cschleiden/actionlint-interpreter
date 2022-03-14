package main

import (
	"github.com/rhysd/actionlint"
)

func Evaluate(n actionlint.ExprNode) interface{} {
	switch tn := n.(type) {
	case *actionlint.IntNode:
		return tn.Value

	case *actionlint.FloatNode:
		return tn.Value

	case *actionlint.StringNode:
		return tn.Value

	case *actionlint.BoolNode:
		return tn.Value

	case *actionlint.NotOpNode:
		return !Evaluate(tn.Operand).(bool)

	case *actionlint.CompareOpNode:
		left := Evaluate(tn.Left)
		right := Evaluate(tn.Right)

		switch tn.Kind {
		case actionlint.CompareOpNodeKindEq:
			return left == right
		}
	}

	panic("unknown node")
}

// func coerce(v interface{}) interface{} {
// 	switch x := v.(type) {
// 		case int:
// 	}
// }
