package reflection

import (
	"reflect"
)

func walk(x interface{}, fn func(input string)) {
	val := getValue(x)

	numberFields := 0
	var getField func(int) reflect.Value

	switch val.Kind() {
	case reflect.Slice, reflect.Array:
		numberFields = val.Len()
		getField = val.Index
	case reflect.Struct:
		numberFields = val.NumField()
		getField = val.Field
	case reflect.Map:
		for _, key := range val.MapKeys() {
			walk(val.MapIndex(key).Interface(), fn)
		}
	case reflect.Chan:
		for v, ok := val.Recv(); ok; v, ok = val.Recv() {
			walk(v.Interface(), fn)
		}
	case reflect.Func:
		result := val.Call(nil)
		for _, v := range result {
			walk(v.Interface(), fn)
		}
	case reflect.String:
		fn(val.String())
	}

	for i := 0; i < numberFields; i++ {
		walk(getField(i).Interface(), fn)
	}

	
}

func getValue(x interface{}) reflect.Value {
	val := reflect.ValueOf(x)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	return val
}