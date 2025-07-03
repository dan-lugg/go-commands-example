package util

import (
	"fmt"
	"github.com/sarulabs/di/v2"
	"reflect"
)

func As[T any](v any) T {
	return v.(T)
}

func TypeNameFor[T any]() (typeName string) {
	t := reflect.TypeFor[T]()
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	typeName = fmt.Sprintf("%s.%s", t.PkgPath(), t.Name())
	return typeName
}

func GetFromContainer[T any](container di.Container) T {
	return container.Get(TypeNameFor[T]()).(T)
}
