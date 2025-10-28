package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.14

	v := reflect.ValueOf(x)
	fmt.Println("Type:", v.Type())   // float64
	fmt.Println("Kind:", v.Kind())   // float64
	fmt.Println("Value:", v.Float()) // 3.14
}
