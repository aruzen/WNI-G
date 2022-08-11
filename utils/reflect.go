package utils

import "reflect"

func Unwrap(v reflect.Value) interface{} {
	switch v.Kind() {
	case reflect.Invalid:
		return nil
	case reflect.Bool:
		return v.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int()
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint()
	case reflect.Float32, reflect.Float64:
		return v.Float()
	case reflect.Complex64, reflect.Complex128:
		return v.Complex()
	case reflect.Array:
		return v.Interface()
	case reflect.Chan:
		return v.Interface()
	case reflect.Func:
		return v.Interface()
	case reflect.Interface:
		return v.Interface()
	case reflect.Map:
		return v.Interface()
	case reflect.Slice:
		return v.Interface()
	case reflect.String:
		return v.String()
	case reflect.Struct:
		tmp := reflect.New(v.Type()) // create zero value of same type as v
		tmp.Elem().Set(v)
		return tmp.Interface()
	case reflect.Pointer:
		return v.Pointer()
	case reflect.Uintptr:
		return v.Uint()
	case reflect.UnsafePointer:
		return v.UnsafePointer()
	}
	return nil
}

func Disassemble(data interface{}, options ...Pair[string, interface{}]) map[string]interface{} {
	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if value.Kind() != reflect.Struct {
		return nil
	}

	result := make(map[string]interface{})

	// options
	nest := 0
	for _, o := range options {
		switch o.Key {
		case "nest":
			if v, ok := o.Data.(int); ok {
				nest = v
				if 0 < v {
					o.Data = v - 1
				}
			}
		}
	}

	size := value.NumField()
	for idx := 0; idx < size; idx++ { // idx: 0 はflg
		name := value.Type().Field(idx).Name
		field := value.FieldByName(name)
		result[name] = Unwrap(field)
		if 0 < nest {
			disassembled := Disassemble(result[name], options...)
			if disassembled != nil {
				result[name] = disassembled
			}
		}
	}

	return result
}
