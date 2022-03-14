package main

import (
	"github.com/rhysd/actionlint"
)

func Evaluate(n actionlint.ExprNode) interface{} {
	switch tn := n.(type) {
	case *actionlint.StringNode:
		return tn.Value

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
