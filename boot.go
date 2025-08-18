package boot

import (
	"fmt"
	"reflect"
)

type wrapper func([]reflect.Value) []reflect.Value

func Boot[G, N any](fn any) any {
	fnValue := reflect.ValueOf(fn)

	if fnValue.Kind() != reflect.Func {
		panic(fmt.Sprintf("gosuit: Boot expects a constructor function as input, but a '%s' is passed", fnValue.Type().String()))
	}

	if fnValue.Type().IsVariadic() {
		panic(fmt.Sprintf("gosuit: '%s' func mustn`t be variadic", fnValue.Type().String()))
	}

	var zeroG G
	var zeroN N

	gType := reflect.TypeOf(zeroG)
	nType := reflect.TypeOf(zeroN)

	if gType.Kind() != reflect.Struct {
		panic("gosuit: G type for Boot must be struct")
	}

	if nType.Kind() != reflect.Struct {
		panic("gosuit: N type for Boot must be struct")
	}

	bootType, gTypeIndex := getBootType(fnValue.Type(), gType, nType)
	bootWrapper := getBootWrap(fnValue, nType, gTypeIndex)

	boot := reflect.MakeFunc(bootType, bootWrapper)

	return boot.Interface()
}

func getBootType(fnType, gType, nType reflect.Type) (reflect.Type, int) {
	if fnType.NumIn() == 0 {
		panic(fmt.Sprintf("gosuit: '%s' must contain '%s' as argument", fnType.String(), nType.String()))
	}

	in := make([]reflect.Type, fnType.NumIn())
	flag := false
	gTypeIndex := 0

	for i := range fnType.NumIn() {
		argType := fnType.In(i)

		if argType == nType || (argType.Kind() == reflect.Pointer && argType.Elem() == nType) {
			in[i] = reflect.PointerTo(gType)
			flag = true
			gTypeIndex = i
		} else {
			in[i] = argType
		}
	}

	if !flag {
		panic(fmt.Sprintf("gosuit: '%s' must contain '%s' as argument", fnType.String(), nType.String()))
	}

	out := make([]reflect.Type, fnType.NumOut())

	for i := range fnType.NumOut() {
		out[i] = fnType.Out(i)
	}

	return reflect.FuncOf(in, out, false), gTypeIndex
}

func getBootWrap(fn reflect.Value, nType reflect.Type, gTypeIndex int) wrapper {
	return func(args []reflect.Value) []reflect.Value {
		g := args[gTypeIndex].Elem()
		n, ok := extractN(g, nType)
		if !ok {
			panic("gosuit: can`t extract N from G")
		}

		isPointer := fn.Type().In(gTypeIndex).Kind() == reflect.Pointer

		if isPointer {
			args[gTypeIndex] = n.Addr()
		} else {
			args[gTypeIndex] = n
		}

		return fn.Call(args)
	}
}

func extractN(g reflect.Value, nType reflect.Type) (reflect.Value, bool) {
	for i := range g.NumField() {
		field := g.Field(i)

		if field.Type() == nType {
			return field, true
		} else if field.Kind() == reflect.Pointer && field.Type().Elem() == nType {
			return field.Elem(), true
		}

		if field.Kind() == reflect.Struct {
			n, ok := extractN(field, nType)
			if ok {
				return n, true
			}
		}
	}

	return reflect.Value{}, false
}
