package util

import "fmt"

func TypeOf(v any) string {
	return fmt.Sprintf("%T", v)
}
