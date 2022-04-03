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
		return &EvaluationResult{Value: float64(tn.Value), Type: &actionlint.NumberType{}}, nil

	case *actionlint.FloatNode:
		return &EvaluationResult{Value: float64(tn.Value), Type: &actionlint.NumberType{}}, nil

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

		vt := getExprType(v)

		return &EvaluationResult{Value: v, Type: vt}, nil

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

		vt := getExprType(v)

		return &EvaluationResult{Value: v, Type: vt}, nil

	case *actionlint.IndexAccessNode:
		arrayResult, err := Evaluate(tn.Operand, context)
		if err != nil {
			return nil, errs.Wrap(err, "could not get operand for index access")
		}

		idxResult, err := Evaluate(tn.Index, context)
		if err != nil {
			return nil, errs.Wrap(err, "could not evalute index for index access")
		}

		if _, ok := arrayResult.Type.(*actionlint.ArrayType); ok {
			return arrayAccess(arrayResult, idxResult)
		}

	//
	// Function call
	//
	case *actionlint.FuncCallNode:
		// Evaluate arguments
		args := make([]*EvaluationResult, len(tn.Args))
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

		return &EvaluationResult{r.Falsy(), &actionlint.BoolType{}}, nil

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

		switch tn.Kind {
		case actionlint.CompareOpNodeKindEq:
			return &EvaluationResult{left.Equals(right), &actionlint.BoolType{}}, nil

		case actionlint.CompareOpNodeKindNotEq:
			return &EvaluationResult{!left.Equals(right), &actionlint.BoolType{}}, nil

		case actionlint.CompareOpNodeKindGreater:
			return &EvaluationResult{left.GreaterThan(right), &actionlint.BoolType{}}, nil

		case actionlint.CompareOpNodeKindGreaterEq:
			return &EvaluationResult{
				left.Equals(right) || left.GreaterThan(right),
				&actionlint.BoolType{},
			}, nil

		case actionlint.CompareOpNodeKindLess:
			return &EvaluationResult{left.LessThan(right), &actionlint.BoolType{}}, nil

		case actionlint.CompareOpNodeKindLessEq:
			return &EvaluationResult{
				left.Equals(right) || left.LessThan(right),
				&actionlint.BoolType{},
			}, nil
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
			left, err := Evaluate(tn.Left, context)
			if err != nil {
				return nil, err
			}

			right, err := Evaluate(tn.Right, context)
			if err != nil {
				return nil, err
			}

			return &EvaluationResult{left.Truthy() && right.Truthy(), &actionlint.BoolType{}}, nil

		case actionlint.LogicalOpNodeKindOr:
			left, err := Evaluate(tn.Left, context)
			if err != nil {
				return nil, err
			}

			if left.Truthy() {
				// No need to evaluate rhs
				return &EvaluationResult{true, &actionlint.BoolType{}}, nil
			}

			right, err := Evaluate(tn.Right, context)
			if err != nil {
				return nil, err
			}

			return &EvaluationResult{right.Truthy(), &actionlint.BoolType{}}, nil
		}
	}

	panic("unknown node")
}

func fcall(name string, args []*EvaluationResult) (*EvaluationResult, error) {
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

func arrayAccess(array *EvaluationResult, idx *EvaluationResult) (*EvaluationResult, error) {
	// TODO: Handle wildcard

	// TODO: Assert type of array
	arrayT := array.Value.([]interface{})

	// Check for number index
	numberIdx := convertToNumber(idx.Value)
	if !math.IsNaN(numberIdx) && numberIdx >= 0.0 {
		idxInt := int(numberIdx)

		if idxInt < 0 || idxInt >= len(arrayT) {
			return nil, errors.New("index out of range")
		}

		v := arrayT[idxInt]
		return &EvaluationResult{v, getExprType(v)}, nil
	}

	return &EvaluationResult{nil, &actionlint.AnyType{}}, nil
}
