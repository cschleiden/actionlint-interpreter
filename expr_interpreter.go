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

func Evaluate(n actionlint.ExprNode, context ContextData) (interface{}, error) {
	switch tn := n.(type) {
	//
	// Literals
	//
	case *actionlint.IntNode:
		return tn.Value, nil

	case *actionlint.FloatNode:
		return tn.Value, nil

	case *actionlint.StringNode:
		return tn.Value, nil

	case *actionlint.BoolNode:
		return tn.Value, nil

	//
	// Context access
	//
	case *actionlint.VariableNode:
		name := tn.Name
		v, ok := context[name]
		if !ok {
			return nil, errors.New("unknown variable access: " + name)
		}

		return v, nil

	case *actionlint.ObjectDerefNode:
		result, err := Evaluate(tn.Receiver, context)
		if err != nil {
			return nil, errs.Wrap(err, "could not evaluate receiver")
		}

		receiverContext, ok := result.(ContextData)
		if !ok {
			return nil, errors.New("invalid result received for receiver")
		}

		property := tn.Property
		v, ok := receiverContext[property]
		if !ok {
			return nil, errors.New("unknown context access: " + property)
		}

		return v, nil

	case *actionlint.IndexAccessNode:
		array, err := Evaluate(tn.Operand, context)
		if err != nil {
			return nil, errs.Wrap(err, "could not get operand for index access")
		}

		idx, err := Evaluate(tn.Index, context)
		if err != nil {
			return nil, errs.Wrap(err, "could not evalute index for index access")
		}

		// TODO: Coerce idx to int
		idxInt := idx.(int)

		// TODO: Assert type of array
		arrayT := array.([]interface{})

		if idxInt < 0 || idxInt >= len(arrayT) {
			return nil, errors.New("index out of range")
		}

		return arrayT[idxInt], nil

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

		return fcall(tn.Callee, args)

	//
	// Unary Operators
	//
	case *actionlint.NotOpNode:
		r, err := Evaluate(tn.Operand, context)
		if err != nil {
			return nil, err
		}

		// TODO: Coerce values
		b := r.(bool)
		return !b, nil

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

		// TODO: Support coercion
		switch tn.Kind {
		case actionlint.CompareOpNodeKindEq:
			return left == right, nil

			// TODO: Support other operators
		}

	case *actionlint.LogicalOpNode:
		left, err := Evaluate(tn.Left, context)
		if err != nil {
			return nil, err
		}
		right, err := Evaluate(tn.Right, context)
		if err != nil {
			return nil, err
		}

		switch tn.Kind {
		case actionlint.LogicalOpNodeKindAnd:
			return isTruthy(left) && isTruthy(right), nil

		case actionlint.LogicalOpNodeKindOr:
			return isTruthy(left) || isTruthy(right), nil
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
