package expr

import (
	"math"
	"strconv"
	"testing"

	"github.com/rhysd/actionlint"
)

func TestEvaluationResult_CoerceString(t *testing.T) {
	type fields struct {
		Value interface{}
		Type  actionlint.ExprType
	}
	tests := []struct {
		fields fields
		want   string
	}{
		{fields{struct{}{}, &actionlint.NullType{}}, ""},
		{fields{false, &actionlint.BoolType{}}, "false"},
		{fields{true, &actionlint.BoolType{}}, "true"},
		{fields{float64(-0), &actionlint.NumberType{}}, "0"},
		{fields{float64(1.234), &actionlint.NumberType{}}, "1.234"},
		{fields{"test", &actionlint.StringType{}}, "test"},
		{fields{[]interface{}{1, 2}, &actionlint.ArrayType{Elem: &actionlint.NumberType{}}}, "array<number>"},
	}
	for _, tt := range tests {
		name := tt.fields.Type.String() + " " + tt.want

		t.Run(name, func(t *testing.T) {
			ev := &EvaluationResult{
				Value: tt.fields.Value,
				Type:  tt.fields.Type,
			}
			if got := ev.CoerceString(); got != tt.want {
				t.Errorf("EvaluationResult.CoerceString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEvaluationResult_Falsy(t *testing.T) {
	type fields struct {
		Value interface{}
		Type  actionlint.ExprType
	}
	tests := []struct {
		fields fields
		want   bool
	}{
		{fields{struct{}{}, &actionlint.NullType{}}, true},
		{fields{false, &actionlint.BoolType{}}, true},
		{fields{true, &actionlint.BoolType{}}, false},
		{fields{float64(-0), &actionlint.NumberType{}}, true},
		{fields{math.NaN(), &actionlint.NumberType{}}, true},
		{fields{float64(1.234), &actionlint.NumberType{}}, false},
		{fields{"", &actionlint.StringType{}}, true},
		{fields{"123", &actionlint.StringType{}}, false},
		{fields{[]interface{}{1, 2}, &actionlint.ArrayType{Elem: &actionlint.NumberType{}}}, false},
	}
	for _, tt := range tests {
		name := tt.fields.Type.String() + " " + strconv.FormatBool(tt.want)
		t.Run(name, func(t *testing.T) {
			ev := &EvaluationResult{
				Value: tt.fields.Value,
				Type:  tt.fields.Type,
			}
			if got := ev.Falsy(); got != tt.want {
				t.Errorf("EvaluationResult.Falsy() = %v, want %v", got, tt.want)
			}
		})
	}
}
