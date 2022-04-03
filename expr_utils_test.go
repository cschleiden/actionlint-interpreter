package expr

import (
	"math"
	"reflect"
	"testing"

	"github.com/rhysd/actionlint"
)

func Test_coerceTypes(t *testing.T) {
	type args struct {
		li interface{}
		ri interface{}
	}
	tests := []struct {
		name      string
		args      args
		wantLv    interface{}
		wantLtype actionlint.ExprType
		wantRv    interface{}
		wantRtype actionlint.ExprType
	}{
		{"number-bool", args{float64(1), true}, float64(1), &actionlint.NumberType{}, float64(1), &actionlint.NumberType{}},
		{"number-bool-false", args{float64(1), false}, float64(1), &actionlint.NumberType{}, float64(0), &actionlint.NumberType{}},
		{"bool-number-false", args{false, float64(1)}, float64(0), &actionlint.NumberType{}, float64(1), &actionlint.NumberType{}},
		{"number-number", args{float64(1), float64(2)}, float64(1), &actionlint.NumberType{}, float64(2), &actionlint.NumberType{}},
		{"string-string", args{"a", "b"}, "a", &actionlint.StringType{}, "b", &actionlint.StringType{}},
		{"string-number", args{"a", float64(1)}, math.NaN(), &actionlint.NumberType{}, float64(1), &actionlint.NumberType{}},
		{"number-string", args{float64(1), "a"}, float64(1), &actionlint.NumberType{}, math.NaN(), &actionlint.NumberType{}},
		{"bool-bool", args{false, true}, false, &actionlint.BoolType{}, true, &actionlint.BoolType{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLv, gotLtype, gotRv, gotRtype := coerceTypes(tt.args.li, tt.args.ri)
			if flv, ok := tt.wantLv.(float64); ok && math.IsNaN(flv) {
				if !math.IsNaN(gotLv.(float64)) {
					t.Errorf("coerceTypes() gotLv = %v, want %v", gotLv, tt.wantLv)
				}
			} else {
				if !reflect.DeepEqual(gotLv, tt.wantLv) {
					t.Errorf("coerceTypes() gotLv = %v, want %v", gotLv, tt.wantLv)
				}
			}
			if !reflect.DeepEqual(gotLtype, tt.wantLtype) {
				t.Errorf("coerceTypes() gotLtype = %v, want %v", gotLtype, tt.wantLtype)
			}
			if frv, ok := tt.wantRv.(float64); ok && math.IsNaN(frv) {
				if !math.IsNaN(gotRv.(float64)) {
					t.Errorf("coerceTypes() gotRv = %v, want %v", gotRv, tt.wantRv)
				}
			} else {
				if !reflect.DeepEqual(gotRv, tt.wantRv) {
					t.Errorf("coerceTypes() gotRv = %v, want %v", gotRv, tt.wantRv)
				}
			}
			if !reflect.DeepEqual(gotRtype, tt.wantRtype) {
				t.Errorf("coerceTypes() gotRtype = %v, want %v", gotRtype, tt.wantRtype)
			}
		})
	}
}

func Test_parseNumber(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want float64
	}{
		{"int", args{"1"}, 1},
		{"neg int", args{"-1"}, -1},
		{"float", args{"1.5"}, 1.5},
		{"neg float", args{"-1.5"}, -1.5},
		{"hex", args{"0xA"}, 10},
		{"oct", args{"0o10"}, 8},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNumber(tt.args.str); got != tt.want {
				t.Errorf("parseNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
