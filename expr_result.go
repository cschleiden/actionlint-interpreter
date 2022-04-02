package expr

import (
	"math"
	"strconv"

	"github.com/rhysd/actionlint"
)

type EvaluationResult struct {
	Value interface{}
	Type  actionlint.ExprType
}

const (
	Expression_True  = "true"
	Expression_False = "false"
)

func (ev *EvaluationResult) CoerceString() string {
	switch tt := ev.Type.(type) {
	case *actionlint.NullType:
		return ""

	case *actionlint.BoolType:
		if b := ev.Value.(bool); b {
			return Expression_True
		} else {
			return Expression_False
		}

	case *actionlint.NumberType:
		// Preserve compat with C# implementation
		if d := ev.Value.(float64); d == -0 {
			return strconv.FormatFloat(float64(0), 'G', 15, 64)
		}

		dv := ev.Value.(float64)
		return strconv.FormatFloat(dv, 'G', 15, 64)

	case *actionlint.StringType:
		return ev.Value.(string)

	default:
		return tt.String()
	}
}

func (ev *EvaluationResult) Falsy() bool {
	switch ev.Type.(type) {
	case *actionlint.NullType:
		return true

	case *actionlint.BoolType:
		return !ev.Value.(bool)

	case *actionlint.NumberType:
		dv := ev.Value.(float64)
		return dv == float64(0) || math.IsNaN(dv)

	case *actionlint.StringType:
		str := ev.Value.(string)
		return str == ""
	default:
		return false
	}
}
func (ev *EvaluationResult) Truthy() bool {
	return !ev.Falsy()
}
