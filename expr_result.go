package expr

import (
	"math"
	"reflect"
	"strconv"
	"strings"

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

func (ev *EvaluationResult) Primitive() bool {
	switch ev.Type.(type) {
	case
		*actionlint.NullType,
		*actionlint.BoolType,
		*actionlint.NumberType,
		*actionlint.StringType:
		return true

	default:
		return false
	}
}

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

func coerceTypes(li interface{}, ri interface{}) (lv interface{}, ltype actionlint.ExprType, rv interface{}, rtype actionlint.ExprType) {
	lv = li
	rv = ri

	ltype = getExprType(li)
	rtype = getExprType(ri)

	lt := reflect.TypeOf(ltype)
	rt := reflect.TypeOf(rtype)

	// Do nothing, same kind
	if lt == rt {
		return
	}

	switch ltype.(type) {
	// Number, String
	case *actionlint.NumberType:
		if _, ok := rtype.(*actionlint.StringType); ok {
			rv = convertToNumber(ri)
			rtype = &actionlint.NumberType{}

			return
		}

	// String, Number
	case *actionlint.StringType:
		if _, ok := rtype.(*actionlint.NumberType); ok {
			lv = convertToNumber(li)
			ltype = &actionlint.NumberType{}

			return
		}

		// Boolean|Null, Any
	case *actionlint.NullType, *actionlint.BoolType:
		lv = convertToNumber(li)
		lv, ltype, rv, rtype = coerceTypes(lv, rv)

		return
	}

	// Any, Boolean|Null
	switch rtype.(type) {
	case *actionlint.NullType, *actionlint.BoolType:
		rv = convertToNumber(ri)
		lv, ltype, rv, rtype = coerceTypes(lv, rv)
		return
	}

	return
}

func getExprType(value interface{}) actionlint.ExprType {
	if value == nil {
		return &actionlint.NullType{}
	}

	switch value.(type) {
	case bool:
		return &actionlint.BoolType{}
	case float64:
		return &actionlint.NumberType{}
	case string:
		return &actionlint.StringType{}
	}

	t := reflect.TypeOf(value)
	if t.Kind() == reflect.Array || t.Kind() == reflect.Slice {
		// et := t.Elem()

		return &actionlint.ArrayType{
			Elem: &actionlint.AnyType{}, // TODO: getExprTypeForType(et),
		}
	}

	return &actionlint.ObjectType{
		Props:  map[string]actionlint.ExprType{}, // TODO: Set types?
		Mapped: &actionlint.AnyType{},            // TODO: Can we make this strict?
	}
}

// func getExprTypeForType(t reflect.Type) actionlint.ExprType {
// 	switch t.Kind() {
// 	case reflect.Bool:
// 		return &actionlint.BoolType{}

// 	case reflect.String:
// 		return &actionlint.StringType{}

// 	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
// 		return &actionlint.NumberType{}

// 	}
// }

func (ev *EvaluationResult) Equals(rhs *EvaluationResult) bool {
	// CoerceTypes(ref canonicalLeftValue, ref canonicalRightValue, out var leftKind, out var rightKind);
	lv, ltype, rv, rtype := coerceTypes(ev.Value, rhs.Value)

	if ltype.String() != rtype.String() {
		return false
	}

	switch ltype.(type) {
	// Null, Null
	case *actionlint.NullType:
		return true

		// Number, Number
	case *actionlint.NumberType:
		ld := lv.(float64)
		rd := rv.(float64)
		if math.IsNaN(ld) || math.IsNaN(rd) {
			return false
		}

		return ld == rd

		// String, String
	case *actionlint.StringType:
		ls := lv.(string)
		rs := rv.(string)
		return strings.EqualFold(ls, rs)

		// Boolean, Boolean
	case *actionlint.BoolType:
		lb := lv.(bool)
		rb := rv.(bool)
		return lb == rb

		// Object, Object
	case *actionlint.ObjectType, *actionlint.ArrayType:
		// Check reference equality
		return lv == rv
	}

	return false
}
