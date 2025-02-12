package utils

import (
	"fmt"
	"reflect"
	"runtime"
	"slices"
	"strconv"

	"analyzer/pkg/logger"
)

func BoolPtr(b bool) *bool {
    return &b
}

var builtInTypes = []string{
	"nil",
	"bool",
	"string",
	"int", "int8", "int16", "int32", "int64",
	"uint", "uint8", "uint16", "uint32", "uint64",
	"byte",
	"rune",
	"float32", "float64",
	"complex64", "complex128",
	"error",
}

var builtInConsts = []string{"true", "false"}
var builtInFuncs = []string{"make", "println", "append", "len", "delete"}

func AddValue(valueType string, s1 string, s2 string) string {
	switch valueType {
	case "string":
		return s1 + s2
	case "int":
		i1, _ := strconv.Atoi(s1)
		i2, _ := strconv.Atoi(s2)
		i1 += i2
		return fmt.Sprintf("%d", i1+i2)
	case "int8":
		var i1, i2 int8
		fmt.Sscan(s1, &i1)
     	fmt.Sscan(s2, &i2)
		i1 += i2
		return fmt.Sprintf("%d", i1+i2)
	case "int16":
		var i1, i2 int8
		fmt.Sscan(s1, &i1)
     	fmt.Sscan(s2, &i2)
		i1 += i2
		return fmt.Sprintf("%d", i1+i2)
	case "int32":
		var i1, i2 int8
		fmt.Sscan(s1, &i1)
     	fmt.Sscan(s2, &i2)
		i1 += i2
		return fmt.Sprintf("%d", i1+i2)
	case "int64":
		var i1, i2 int8
		fmt.Sscan(s1, &i1)
     	fmt.Sscan(s2, &i2)
		return fmt.Sprintf("%d", i1+i2)
	default:
		logger.Logger.Fatalf("[UTILS ADD VALUE] unknown type (%s) for adding values (%s) and (%s)", valueType, s1, s2)
	}
	return s1 + s2
}

func GetType(node interface{}) string {
	if node == nil {
		pc, file, line, ok := runtime.Caller(1)
		if !ok {
			logger.Logger.Warnf("[UTILS TYPE] INVALID TYPE FOR <nil> and unable to retrieve caller information")
		}
		callerFunc := runtime.FuncForPC(pc).Name()
		logger.Logger.Warnf("[UTILS TYPE] INVALID TYPE FOR <nil> (caller: %s) \n\t\t\t\t\t %s:%d", callerFunc, file, line)
		return ""
	}
	return reflect.TypeOf(node).Elem().Name()
}

func IsType[T any](obj interface{}) bool {
	_, ok := obj.(T)
	return ok
}

func IsBuiltInGoTypeOrFunc(name string) bool {
	return IsBuiltInGoFunc(name) || IsBuiltInGoType(name)
}

func IsBuiltInGoFunc(name string) bool {
	return slices.Contains(builtInFuncs, name)
}

func IsBuiltInGoType(name string) bool {
	return slices.Contains(builtInTypes, name)
}

func IsBuiltInConstValue(name string) bool {
	return slices.Contains(builtInConsts, name)
}

func GetBuiltInConstTypeName(name string) string {
	switch name {
	case "true":
		return "bool"
	case "false":
		return "bool"
	default:
		logger.Logger.Fatalf("[UTILS] unknown built in const type/value for name (%s)", name)
	}
	return ""
}
