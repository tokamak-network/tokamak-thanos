package enum

import (
	"fmt"
	"strings"
)

// EnumString returns a comma-separated string of the enum values.
// This is primarily used to generate a cli flag.
func EnumString[T ~string](values []T) string {
	var out strings.Builder
	for i, v := range values {
		out.WriteString(string(v))
		if i+1 < len(values) {
			out.WriteString(", ")
		}
	}
	return out.String()
}

// EnumStringer is an alias for EnumString for fmt.Stringer types.
func EnumStringer[T fmt.Stringer](values []T) string {
	strs := make([]string, len(values))
	for i, v := range values {
		strs[i] = v.String()
	}
	return strings.Join(strs, ", ")
}
