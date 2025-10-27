package main

import (
	"reflect"
)

func main() {
	var x float64 = 3.14

	v := reflect.ValueOf(x)
	//EVAL - fmt.Println("Type:", v.Type())         // float64
	//EVAL - fmt.Println("Kind:", v.Kind())         // float64
	//EVAL - fmt.Println("Value:", v.Float())       // 3.14
}
