package utils

import "fmt"

func RemoveFromArray(source []any, r []any) []any {
	for i := 0; i < len(source); i++ {
		si := source[i]
		for _, rem := range r {
			if fmt.Sprint(si) == fmt.Sprint(rem) {
				source = append(source[:i], source[i+1:]...)
				i--
				break
			}
		}
	}
	return source
}

func ArrToArrAny[T any](source []T) []any {
	result := make([]any, len(source))
	for i, s := range source {
		result[i] = s
	}
	return result
}
