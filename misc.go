package bencode

import (
	"fmt"
	"reflect"
	"unsafe"
)

// Wow Go is retarded.
var (
	marshalerType   = reflect.TypeOf((*Marshaler)(nil)).Elem()
	unmarshalerType = reflect.TypeOf((*Unmarshaler)(nil)).Elem()
)

func bytesAsString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func expectNil(x interface{}) {
	if x != nil {
		panic(fmt.Sprintf("expected nil; got %v", x))
	}
}

func isZeroValue(i interface{}) bool {
	return isEmptyValue(reflect.ValueOf(i))
}

// Returns whether the value represents the empty value for its type. Used for
// example to determine if complex types satisfy the common "omitempty" tag
// option for marshalling. Taken from
// http://stackoverflow.com/a/23555352/149482.
func isEmptyValue(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Func, reflect.Map, reflect.Slice:
		return v.IsNil()
	case reflect.Array:
		z := true
		for i := 0; i < v.Len(); i++ {
			z = z && isEmptyValue(v.Index(i))
		}
		return z
	case reflect.Struct:
		z := true
		vType := v.Type()
		for i := 0; i < v.NumField(); i++ {
			// ignore unexported fields to avoid reflection panics
			// just use PkgPath == "" to replace IsExported reports whether the struct field is exported.
			if vType.Field(i).PkgPath == "" {
				continue
			}
			z = z && isEmptyValue(v.Field(i))
		}
		return z
	}
	// Compare other types directly:
	z := reflect.Zero(v.Type())
	return v.Interface() == z.Interface()
}
