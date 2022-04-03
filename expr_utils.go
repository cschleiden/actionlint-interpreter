package expr

import (
	"math"
	"strconv"
	"strings"

	"github.com/rhysd/actionlint"
)

// parseNumber attempts to follow Javascript rules for coercing a string into a number
// for comparison. That is, the Number() function in Javascript.
func parseNumber(str string) float64 {
	str = strings.TrimSpace(str)

	if str == "" {
		return 0.0
	}

	if v, err := strconv.ParseFloat(str, 64); err == nil {
		return v
	}

	if str[0] == '0' && len(str) > 2 && str[1] == 'x' && strAll(str[2:], func(x rune) bool { return (x >= '0' && x <= '9') || (x >= 'a' && x <= 'f') || (x >= 'A' && x <= 'F') }) {
		if v, err := strconv.ParseInt(str[2:], 16, 64); err == nil {
			return float64(v)
		}

		// Exceeds range
	}

	if str[0] == '0' && len(str) > 2 && str[1] == 'o' && strAll(str[2:], func(x rune) bool { return x >= '0' && x <= '7' }) {
		if v, err := strconv.ParseInt(str[2:], 8, 32); err == nil {
			return float64(v)
		}
	}

	if str == "Infinity" {
		return math.Inf(1)
	}

	if str == "-Infinity" {
		return math.Inf(-1)
	}

	return math.NaN()
}

func strAll(str string, f func(r rune) bool) bool {
	for _, r := range str {
		if !f(r) {
			return false
		}
	}

	return true
}

func convertToNumber(v interface{}) float64 {
	ltype := getExprType(v)

	switch ltype.(type) {
	case *actionlint.NullType:
		return float64(0)

	case *actionlint.BoolType:
		if v.(bool) {
			return float64(1)
		} else {
			return float64(0)
		}

	case *actionlint.NumberType:
		return v.(float64)

	case *actionlint.StringType:
		return parseNumber(v.(string))
	}

	return math.NaN()
}
