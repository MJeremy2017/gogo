package main

import (
	"fmt"
	"reflect"
)

type Car struct {
	Name string
	Wheels int
	Color Color
}

type Color struct {
	color1 string
	color2 string
}

func main() {
	fmt.Println("Hello")
	car := Car {
		"Ford",
		3,
		Color {
			"green",
			"yellow",
		},
	}

	val := reflect.ValueOf(car)
	fmt.Printf("values %v\n", val)

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		xx := reflect.ValueOf(field)
		
		fmt.Printf("value of interface: %v\n", reflect.ValueOf(field.Interface()))
		fmt.Printf("xx %v\n", xx)
		fmt.Printf("field %v | type %v | after type %v\n", field, field.Kind(), field.Interface())
	}
	

}