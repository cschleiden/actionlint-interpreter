package main

import (
	"reflect"
	"testing"

	"github.com/rhysd/actionlint"
)

func Test_Evaluate(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  interface{}
	}{
		{
			name:  "string literal",
			input: "'test'",
			want:  "test",
		},
		{
			name:  "number",
			input: "12",
			want:  12,
		},
		{
			name:  "negative number",
			input: "-12",
			want:  -12,
		},
		{
			name:  "number 0",
			input: "0",
			want:  0,
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
			n, err := parser.Parse(lexer)
			if err != nil {
				t.Fatal(err.Error())
			}

			if got := Evaluate(n); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("evaluate() = %v, want %v", got, tt.want)
			}
		})
	}
}
