package utils

import (
	"reflect"
)

func CompareMap[K comparable, T any](a, b map[K]T) bool {
	for k := range a {
		vb, ok := b[k]
		if !ok {
			return false
		}
		if !Compare(a[k], vb) {
			return false
		}
	}
	return true
}

func CompareReflectMap(a, b reflect.Value) bool {
	if len(a.MapKeys()) != len(b.MapKeys()) {
		return false
	}
	ai := a.MapRange()
	bi := b.MapRange()
	for ai.Next() && bi.Next() {
		v := CompareReflect(ai.Value(), bi.Value())
		if !CompareReflect(ai.Key(), bi.Key()) || !v {
			return false
		}
	}
	return true
}

func CompareSlice[T any](a, b T) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	return CompareReflectSlice(va, vb)
}

func CompareReflectSlice(av, bv reflect.Value) bool {
	if av.Len() != bv.Len() {
		return false
	}
	for idx := 0; idx < av.Len(); idx++ {
		if !CompareReflect(av.Index(idx), bv.Index(idx)) {
			return false
		}
	}
	return true
}

func CompareStruct(a, b interface{}) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	return CompareReflectStruct(va, vb)
}

func CompareReflectStruct(va, vb reflect.Value) bool {
	if va.NumField() != vb.NumField() {
		return false
	}
	for i := 0; i < va.NumField(); i++ {
		if !CompareReflect(va.Field(i), vb.Field(i)) {
			return false
		}
	}
	return true
}

func Compare(a, b interface{}) bool {
	va := reflect.ValueOf(a)
	vb := reflect.ValueOf(b)
	return CompareReflect(va, vb)
}

func CompareReflect(va, vb reflect.Value) bool {
	for va.Kind() == reflect.Pointer {
		va = va.Elem()
	}
	for vb.Kind() == reflect.Pointer {
		vb = vb.Elem()
	}
	switch va.Kind() {
	case reflect.Invalid:
		return false
	case reflect.Bool:
		if vb.Kind() != reflect.Bool {
			return false
		}
		return va.Bool() == vb.Bool()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch vb.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			return va.Int() == vb.Int()
		default:
			return false
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		switch vb.Kind() {
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
			return va.Uint() == vb.Uint()
		default:
			return false
		}
	case reflect.Float32, reflect.Float64:
		switch vb.Kind() {
		case reflect.Float32, reflect.Float64:
			return va.Float() == vb.Float()
		default:
			return false
		}
	case reflect.Complex64, reflect.Complex128:
		switch vb.Kind() {
		case reflect.Complex64, reflect.Complex128:
			return va.Complex() == vb.Complex()
		default:
			return false
		}
	case reflect.String:
		if vb.Kind() != reflect.String {
			return false
		}
		return va.String() == vb.String()
	case reflect.Array:
		if vb.Kind() != reflect.Array {
			return false
		}
		return CompareReflectSlice(va, vb)
	case reflect.Chan:
		return va.Interface() == vb.Interface()
	case reflect.Func:
		return va.Interface() == vb.Interface()
	case reflect.Interface:
		return Compare(va.Interface(), vb.Interface())
	case reflect.Map:
		if vb.Kind() != reflect.Map {
			return false
		}
		return CompareReflectMap(va, vb)
	case reflect.Slice:
		if vb.Kind() != reflect.Slice {
			return false
		}
		return CompareReflectSlice(va, vb)
	case reflect.Struct:
		if vb.Kind() != reflect.Struct {
			return false
		}
		return CompareReflectStruct(va, vb)
	case reflect.Pointer:
		if vb.Kind() != reflect.Pointer {
			return false
		}
		return va.Pointer() == vb.Pointer()
	case reflect.UnsafePointer:
		if vb.Kind() != reflect.UnsafePointer {
			return false
		}
		return va.UnsafePointer() == vb.UnsafePointer()
	}
	return va.Interface() == vb.Interface()
}
