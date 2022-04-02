package expr

import (
	"errors"
	"fmt"
	"math"
	"strings"

	errs "github.com/pkg/errors"

	"github.com/rhysd/actionlint"
)

type ContextData = map[string]interface{}

func Evaluate(n actionlint.ExprNode, context ContextData) (*EvaluationResult, error) {
	switch tn := n.(type) {
	//
	// Literals
	//
	case *actionlint.IntNode:
		return &EvaluationResult{Value: tn.Value, Type: &actionlint.NumberType{}}, nil

	case *actionlint.FloatNode:
		return &EvaluationResult{Value: tn.Value, Type: &actionlint.NumberType{}}, nil

	case *actionlint.StringNode:
		return &EvaluationResult{Value: tn.Value, Type: &actionlint.StringType{}}, nil

	case *actionlint.BoolNode:
		return &EvaluationResult{Value: tn.Value, Type: &actionlint.BoolType{}}, nil

	//
	// Context access
	//
	case *actionlint.VariableNode:
		name := tn.Name
		v, ok := context[name]
		if !ok {
			return nil, errors.New("unknown variable access: " + name)
		}

		return &EvaluationResult{Value: v, Type: &actionlint.AnyType{}}, nil

	case *actionlint.ObjectDerefNode:
		result, err := Evaluate(tn.Receiver, context)
		if err != nil {
			return nil, errs.Wrap(err, "could not evaluate receiver")
		}

		value := result.Value

		// TODO: Is this always ContextData?
		receiverContext, ok := value.(ContextData)
		if !ok {
			return nil, errors.New("invalid result received for receiver")
		}

		property := tn.Property
		v, ok := receiverContext[property]
		if !ok {
			return nil, errors.New("unknown context access: " + property)
		}

		// TODO: Should we try to determine the type here?
		return &EvaluationResult{Value: v, Type: &actionlint.AnyType{}}, nil

	case *actionlint.IndexAccessNode:
		arrayResult, err := Evaluate(tn.Operand, context)
		if err != nil {
			return nil, errs.Wrap(err, "could not get operand for index access")
		}

		if _, ok := arrayResult.Type.(*actionlint.ArrayType); !ok {
			return nil, errors.New("index access is not supported for non-array type")
		}

		idxResult, err := Evaluate(tn.Index, context)
		if err != nil {
			return nil, errs.Wrap(err, "could not evalute index for index access")
		}

		// TODO: Coerce other types?
		if _, ok := idxResult.Type.(*actionlint.NumberType); !ok {
			return nil, errors.New("index has to be a number type")
		}

		// TODO: Coerce idx to int
		idxInt := idxResult.Value.(int)

		// TODO: Assert type of array
		arrayT := arrayResult.Value.([]interface{})

		if idxInt < 0 || idxInt >= len(arrayT) {
			return nil, errors.New("index out of range")
		}

		// TODO: assert type?
		return &EvaluationResult{Value: arrayT[idxInt], Type: &actionlint.AnyType{}}, nil

	//
	// Function call
	//
	case *actionlint.FuncCallNode:
		// Evaluate arguments
		args := make([]interface{}, len(tn.Args))
		for i, arg := range tn.Args {
			a, err := Evaluate(arg, context)
			if err != nil {
				return nil, err
			}

			args[i] = a
		}

		// return fcall(tn.Callee, args)

		return nil, errors.New("not implemented")

	//
	// Unary Operators
	//
	case *actionlint.NotOpNode:
		_, err := Evaluate(tn.Operand, context)
		if err != nil {
			return nil, err
		}

		// TODO: Coerce values
		// b := r.(bool)
		//return !b, nil
		return nil, errors.New("not implemented")

	//
	// Binary Operators
	//
	case *actionlint.CompareOpNode:
		left, err := Evaluate(tn.Left, context)
		if err != nil {
			return nil, err
		}
		right, err := Evaluate(tn.Right, context)
		if err != nil {
			return nil, err
		}

		if !left.Type.Assignable(right.Type) {
			// TODO: Coerce values
			return nil, errors.New("incompatible types for comparison")
		}

		// TODO: Support coercion
		switch tn.Kind {
		case actionlint.CompareOpNodeKindEq:
			return &EvaluationResult{left.Value == right.Value, &actionlint.BoolType{}}, nil

		case actionlint.CompareOpNodeKindNotEq:
			return &EvaluationResult{left.Value != right.Value, &actionlint.BoolType{}}, nil

			// TODO: Support other operators
		}

	case *actionlint.LogicalOpNode:
		_, err := Evaluate(tn.Left, context)
		if err != nil {
			return nil, err
		}
		_, err = Evaluate(tn.Right, context)
		if err != nil {
			return nil, err
		}

		switch tn.Kind {
		case actionlint.LogicalOpNodeKindAnd:
			// return isTruthy(left) && isTruthy(right), nil
			return nil, nil

		case actionlint.LogicalOpNodeKindOr:
			// return isTruthy(left) || isTruthy(right), nil
			return nil, nil
		}
	}

	panic("unknown node")
}

func fcall(name string, args []interface{}) (interface{}, error) {
	// Expression function names are case-insensitive.
	funcDef, ok := functions[strings.ToLower(name)]
	if !ok {
		return nil, errors.New("unknown function: " + name)
	}

	if funcDef.argsCount >= 0 {
		if funcDef.argsCount != len(args) {
			return nil, errors.New(fmt.Sprintf("invalid number of arguments. expected %d, got %d", funcDef.argsCount, len(args)))
		}
	} else {
		if int(math.Abs(float64(funcDef.argsCount))) > len(args) {
			return nil, errors.New(fmt.Sprintf("invalid number of arguments. expected at least %d, got %d", funcDef.argsCount, len(args)))
		}
	}

	return funcDef.call(args...), nil
}
