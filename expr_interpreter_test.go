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
		context map[string]interface{}
	}{
		{
			name:  "string literal",
			input: "'test'",
			want:  "test",
		},
		{
			name:  "int",
			input: "12",
			want:  12,
		},
		{
			name:  "negative int",
			input: "-12",
			want:  -12,
		},
		{
			name:  "int 0",
			input: "0",
			want:  0,
		},
		{
			name:  "float",
			input: "12.5",
			want:  float64(12.5),
		},
		{
			name:  "negative float",
			input: "-12.3",
			want:  float64(-12.3),
		},
		{
			name:  "float 0",
			input: "0.0",
			want:  float64(0.0),
		},
		{
			name:  "bool true",
			input: "true",
			want:  true,
		},
		{
			name:  "bool false",
			input: "false",
			want:  false,
		},
		{
			name:  "neg operator",
			input: "!false",
			want:  true,
		},
		{
			name:  "neg operator",
			input: "!true",
			want:  false,
		},
		{
			name:    "context access - one level",
			input:   "input",
			context: map[string]interface{}{"input": 42},
			want:    42,
		},
		{
			name:    "context access - object dereferece",
			input:   "input.test2.test",
			context: map[string]interface{}{"input": map[string]interface{}{"test2": map[string]interface{}{"test": 42}}},
			want:    42,
		},
		{
			name:    "context access - mixed dereferece",
			input:   "input.test[1]",
			context: map[string]interface{}{"input": map[string]interface{}{"test": []interface{}{23, 42}}},
			want:    42,
		},
		{
			name:  "comparison eq - equal strings",
			input: "'test' == 'test'",
			want:  true,
		},
		{
			name:  "comparison eq - not equal strings",
			input: "'test' == 'test2'",
			want:  false,
		},
		{
			name:  "comparison eq - not equal numbers",
			input: "12 == 13",
			want:  false,
		},
		{
			name:  "comparison eq - equal numbers",
			input: "12 == 12",
			want:  true,
		},
		{
			name:  "comparison eq - bool not equal",
			input: "true == false",
			want:  false,
		},
		{
			name:  "comparison eq - bool equal",
			input: "false == false",
			want:  true,
		},
		{
			name:  "comparison eq - bool equal",
			input: "true == true",
			want:  true,
		},
		{
			name:  "fcall - startsWith",
			input: "startsWith('test', 'tE')",
			want:  true,
		},
		{
			name:  "fcall - startsWith - false",
			input: "startsWith('test', 'xe')",
			want:  false,
		},
		{
			name:    "fcall - join",
			input:   "join(inputs.values)",
			context: map[string]interface{}{"inputs": map[string]interface{}{"values": []interface{}{"42", "1"}}},
			want:    "42,1",
		},
		{
			name:    "fcall - join - custom seperator",
			input:   "join(inputs.values, ':')",
			context: map[string]interface{}{"inputs": map[string]interface{}{"values": []interface{}{"42", "1"}}},
			want:    "42:1",
		},
		// {
		// 	name:  "comparison eq - equal string number",
		// 	input: "'2' == 2",
		// 	want:  true,
		// },
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
