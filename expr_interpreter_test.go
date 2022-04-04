package expr

import (
	"reflect"
	"testing"

	"github.com/rhysd/actionlint"
)

func Test_Evaluate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    interface{}
		context ContextData
	}{
		{
			name:  "string literal",
			input: "'test'",
			want:  &EvaluationResult{Value: "test", Type: &actionlint.StringType{}},
		},
		{
			name:  "int",
			input: "12",
			want:  &EvaluationResult{Value: float64(12), Type: &actionlint.NumberType{}},
		},
		{
			name:  "negative int",
			input: "-12",
			want:  &EvaluationResult{Value: float64(-12), Type: &actionlint.NumberType{}},
		},
		{
			name:  "int 0",
			input: "0",
			want:  &EvaluationResult{Value: float64(0), Type: &actionlint.NumberType{}},
		},
		{
			name:  "float",
			input: "12.5",
			want:  &EvaluationResult{Value: float64(12.5), Type: &actionlint.NumberType{}},
		},
		{
			name:  "negative float",
			input: "-12.3",
			want:  &EvaluationResult{Value: float64(-12.3), Type: &actionlint.NumberType{}},
		},
		{
			name:  "float 0",
			input: "0.0",
			want:  &EvaluationResult{Value: float64(0.0), Type: &actionlint.NumberType{}},
		},
		{
			name:  "bool true",
			input: "true",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "bool false",
			input: "false",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:  "neg operator",
			input: "!false",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "neg operator",
			input: "!true",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:    "context access",
			input:   "input",
			context: map[string]interface{}{"input": float64(42)},
			want:    &EvaluationResult{Value: float64(42), Type: &actionlint.NumberType{}},
		},
		{
			name:    "context access - . object access",
			input:   "input.test2.test",
			context: map[string]interface{}{"input": map[string]interface{}{"test2": map[string]interface{}{"test": float64(42)}}},
			want:    &EvaluationResult{Value: float64(42), Type: &actionlint.NumberType{}},
		},
		{
			name:    "context access - [] object access",
			input:   "input.test2['test']",
			context: map[string]interface{}{"input": map[string]interface{}{"test2": map[string]interface{}{"test": float64(42)}}},
			want:    &EvaluationResult{Value: float64(42), Type: &actionlint.NumberType{}},
		},
		{
			name:    "context access - array access",
			input:   "input.test[1]",
			context: map[string]interface{}{"input": map[string]interface{}{"test": []interface{}{float64(23), float64(42)}}},
			want:    &EvaluationResult{Value: float64(42), Type: &actionlint.NumberType{}},
		},
		// {
		// 	name:  "context access - wildcard",
		// 	input: "input.*.foo",
		// 	context: map[string]interface{}{"input": map[string]interface{}{
		// 		"test":  map[string]interface{}{"foo": float64(32)},
		// 		"test2": map[string]interface{}{"foo": float64(42)},
		// 	}},
		// 	want: &EvaluationResult{Value: []float64{32, 42}, Type: &actionlint.ArrayType{}},
		// },
		{
			name:  "comparison eq - equal strings",
			input: "'test' == 'test'",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison eq - not equal strings",
			input: "'test' == 'test2'",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison eq - not equal numbers",
			input: "12 == 13",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison eq - equal numbers",
			input: "12 == 12",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison neq",
			input: "12 != 14",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison gt - false",
			input: "12 > 14",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison gt - true",
			input: "14 > 12",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison gteq - true",
			input: "12 >= 12",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison lt - false",
			input: "14 < 12",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison lt - true",
			input: "12 < 14",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison lteq - true",
			input: "12 <= 12",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison eq - bool not equal",
			input: "true == false",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison eq - bool equal",
			input: "false == false",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison eq - bool equal",
			input: "true == true",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison eq - equal string number",
			input: "'2' == 2",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "comparison eq - equal number string",
			input: "2 == '2'",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "logical or - true",
			input: "true || false",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "logical or - false",
			input: "false || false",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:  "logical and - true",
			input: "true && true",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "logical and - false",
			input: "true && false",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:  "logical and complex - true",
			input: "true && (1 == 2)",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:  "logical or complex - true",
			input: "(1 == 2) || (1 == 1)",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "fcall - startsWith",
			input: "startsWith('test', 'tE')",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "fcall - startsWith - false",
			input: "startsWith('test', 'xe')",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:  "fcall - endsWith",
			input: "endsWith('test', 'sT')",
			want:  &EvaluationResult{Value: true, Type: &actionlint.BoolType{}},
		},
		{
			name:  "fcall - startsWith - false",
			input: "endsWith('test', 'xe')",
			want:  &EvaluationResult{Value: false, Type: &actionlint.BoolType{}},
		},
		{
			name:    "fcall - join",
			input:   "join(inputs.values)",
			context: map[string]interface{}{"inputs": map[string]interface{}{"values": []interface{}{"42", "1"}}},
			want:    &EvaluationResult{Value: "42,1", Type: &actionlint.StringType{}},
		},
		{
			name:  "fcall - string",
			input: "join('foo')",
			want:  &EvaluationResult{Value: "foo", Type: &actionlint.StringType{}},
		},
		{
			name:    "fcall - join - custom seperator",
			input:   "join(inputs.values, ':')",
			context: map[string]interface{}{"inputs": map[string]interface{}{"values": []interface{}{"42", "1"}}},
			want:    &EvaluationResult{Value: "42:1", Type: &actionlint.StringType{}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := actionlint.NewExprLexer(tt.input + "}}")
			parser := actionlint.NewExprParser()
			n, perr := parser.Parse(lexer)
			if perr != nil {
				t.Fatal(perr.Error())
			}

			got, err := Evaluate(n, tt.context)
			if err != nil {
				t.Errorf("Evaluate() error = %v", err)
			} else if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
