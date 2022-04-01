package expr

import "strings"

type funcDef struct {
	// argsCount is the number of required arguments. Positive values have to be matched exactly,
	// negative values indicate the abs(minimum) number of arguments required
	argsCount int

	call func(args ...interface{}) interface{}
}

var functions map[string]funcDef = map[string]funcDef{
	"startswith": {
		argsCount: 2,
		call: func(args ...interface{}) interface{} {
			// TODO: Check types of parameters

			left := args[0].(string)
			right := args[1].(string)

			// Expression string comparisons are string insensitive
			return strings.HasPrefix(strings.ToLower(left), strings.ToLower(right))
		},
	},

	"join": {
		argsCount: -1,
		call: func(args ...interface{}) interface{} {
			separator := ","

			// if args are an array, join
			ar := args[0].([]interface{})

			if len(args) > 1 {
				separator = args[1].(string)
			}

			v := make([]string, len(ar))
			for i, a := range ar {
				// TODO: Support other types besides string
				v[i] = a.(string)
			}

			return strings.Join(v, separator)
		},
	},
}
