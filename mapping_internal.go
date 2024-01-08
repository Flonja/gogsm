package gogsm

import (
	"errors"
	"strings"
)

func wrap(vs ...any) []any {
	return vs
}

func replace(find string) func(s string) string {
	return func(s string) string {
		return strings.Replace(s, find, "", -1)
	}
}

func split(separator string) func(s string) []string {
	return func(s string) []string {
		return strings.Split(s, separator)
	}
}

func mapArray[T any, V any](f func(s T) V) func([]T) []V {
	return func(collection []T) []V {
		var arr []V
		for _, i := range collection {
			arr = append(arr, f(i))
		}
		return arr
	}
}

func mapped[T any, V any](f func(T) V, variables ...any) (val V, err error) {
	isMapped := false
	for _, variable := range variables {
		if e, ok := variable.(error); ok {
			err = e
		}

		if v, ok := variable.(T); ok {
			isMapped = true
			val = f(v)
		}
	}
	if !isMapped {
		return val, errors.New("not successfully mapped")
	}
	return val, err
}
