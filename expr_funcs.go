package expr

import (
	"strings"

	"github.com/rhysd/actionlint"
)

type funcDef struct {
	// argsCount is the number of required arguments. Positive values have to be matched exactly,
	// negative values indicate the abs(minimum) number of arguments required
	argsCount int

	call func(args ...*EvaluationResult) *EvaluationResult
}

var functions map[string]funcDef = map[string]funcDef{
	"startswith": {
		argsCount: 2,
		call: func(args ...*EvaluationResult) *EvaluationResult {
			// TODO: Check types of parameters
			left := args[0]
			if !left.Primitive() {
				return &EvaluationResult{false, &actionlint.BoolType{}}
			}

			right := args[1]
			if !left.Primitive() {
				return &EvaluationResult{false, &actionlint.BoolType{}}
			}

			ls := left.CoerceString()
			rs := right.CoerceString()

			// Expression string comparisons are string insensitive
			return &EvaluationResult{strings.HasPrefix(strings.ToLower(ls), strings.ToLower(rs)), &actionlint.BoolType{}}
		},
	},

	"join": {
		argsCount: -1,
		call: func(args ...*EvaluationResult) *EvaluationResult {
			separator := ","

			// String
			if args[0].Primitive() {
				return args[0]
			}

			if len(args) > 1 {
				separator = args[1].CoerceString()
			}

			ar := args[0].Value.([]interface{})

			v := make([]string, len(ar))
			for i, a := range ar {
				ar := &EvaluationResult{a, getExprType(a)}
				v[i] = ar.CoerceString()
			}

			return &EvaluationResult{strings.Join(v, separator), &actionlint.StringType{}}
		},
	},
}
