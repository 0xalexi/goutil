package goutil

import "reflect"

func NewOfType(value interface{}) interface{} {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return reflect.Indirect(reflect.New(reflect.TypeOf(value))).Interface()
	}
	return reflect.Indirect(reflect.New(reflect.TypeOf(value).Elem())).Interface()
}

func NewOfReflectType(v reflect.Type) interface{} {
	if v.Kind() != reflect.Ptr {
		return reflect.Indirect(reflect.New(v)).Interface()
	}
	return reflect.Indirect(reflect.New(v.Elem())).Interface()
}

func Pointer_NewOfType(value interface{}) interface{} {
	v := reflect.ValueOf(value)
	if v.Kind() != reflect.Ptr {
		return reflect.New(reflect.TypeOf(value)).Interface()
	}
	return reflect.New(reflect.TypeOf(value).Elem()).Interface()
}

func Pointer_NewOfReflectType(v reflect.Type) interface{} {
	if v.Kind() != reflect.Ptr {
		return reflect.New(v).Interface()
	}
	return reflect.New(v.Elem()).Interface()
}

// Ignores whether a and b are pointers
func TypesMatch(a, b interface{}) bool {
	ta := reflect.TypeOf(a)
	if ta.Kind() == reflect.Ptr {
		ta = ta.Elem()
	}
	tb := reflect.TypeOf(b)
	if tb.Kind() == reflect.Ptr {
		tb = tb.Elem()
	}
	return ta == tb
}

func GetReflectValue(value interface{}) (bool, reflect.Value) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		el := v.Elem()
		if !el.IsValid() {
			return false, v
		}
		v = reflect.ValueOf(el.Interface())
	}
	return true, v
}

func GetReflectPointerValue(value interface{}) reflect.Value {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		return v
	}
	return reflect.Indirect(v)
}

func RefInterface(value interface{}) interface{} {
	v := GetReflectPointerValue(value)
	return v.Interface()
}

func DerefInterface(value interface{}) interface{} {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		el := v.Elem()
		if !el.IsValid() {
			return nil
		}
		return el.Interface()
	}
	return v
}

func IsSlice(value interface{}) bool {
	isvalid, v := getReflectValue(value)
	if !isvalid {
		return false
	}
	return v.Type().Kind() == reflect.Slice
}

func getReflectValue(value interface{}) (bool, reflect.Value) {
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Ptr {
		el := v.Elem()
		if !el.IsValid() {
			return false, v
		}
		v = reflect.ValueOf(el.Interface())
	}
	return true, v
}
