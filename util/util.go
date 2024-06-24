package util

import (
	"fmt"
	"kat/value"
)

func TypeOf(v any) string {
	return fmt.Sprintf("%T", v)
}

func IsTruthy(v value.Value) bool {
	switch v.(type) {
	case value.ValueInt:
		return v.(value.ValueInt).Value != 0

	case value.ValueFloat:
		return v.(value.ValueFloat).Value != 0

	case value.ValueBool:
		return v.(value.ValueBool).Value
	}

	return false
}
