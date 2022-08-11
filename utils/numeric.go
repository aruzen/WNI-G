package utils

import (
	"math"
	"reflect"
)

func FitBit(n interface{}) interface{} {
	var wide int64
	switch n.(type) {
	case int8:
		wide = int64(n.(int8))
	case int16:
		wide = int64(n.(int16))
	case int32:
		wide = int64(n.(int32))
	case int:
		wide = int64(n.(int))
	case int64:
		wide = n.(int64)
	}
	if math.MinInt8 <= wide && wide <= math.MaxInt8 {
		return int8(wide)
	}
	if math.MinInt16 <= wide && wide <= math.MaxInt16 {
		return int16(wide)
	}
	if math.MinInt32 <= wide && wide <= math.MaxInt32 {
		return int8(wide)
	}
	return wide
}

func WidenInt(n interface{}) int64 {
	switch n.(type) {
	case int8:
		return int64(n.(int8))
	case int16:
		return int64(n.(int16))
	case int32:
		return int64(n.(int32))
	case int:
		return int64(n.(int))
	case int64:
		return n.(int64)
	default:
		return 0
	}
}

func WidenFloat(n interface{}) float64 {
	switch n.(type) {
	case float32:
		return float64(n.(float32))
	case float64:
		return n.(float64)
	default:
		return 0
	}
}

func AssignIntValue(n interface{}, a reflect.Value) bool {
	var wide int64
	switch n.(type) {
	case int8:
		wide = int64(n.(int8))
	case int16:
		wide = int64(n.(int16))
	case int32:
		wide = int64(n.(int32))
	case int:
		wide = int64(n.(int))
	case int64:
		wide = n.(int64)
	default:
		return false
	}
	switch a.Kind() {
	case reflect.Int:
		a.Set(reflect.ValueOf(int(wide)))
	case reflect.Int8:
		a.Set(reflect.ValueOf(int8(wide)))
	case reflect.Int16:
		a.Set(reflect.ValueOf(int16(wide)))
	case reflect.Int32:
		a.Set(reflect.ValueOf(int32(wide)))
	case reflect.Int64:
		a.Set(reflect.ValueOf(wide))
	default:
		return false
	}
	return true
}

func AssignFloatValue(n interface{}, a reflect.Value) bool {
	var wide float64
	switch n.(type) {
	case float32:
		wide = float64(n.(float32))
	case float64:
		wide = n.(float64)
	default:
		return false
	}
	switch a.Kind() {
	case reflect.Float32:
		a.Set(reflect.ValueOf(float32(wide)))
	case reflect.Float64:
		a.Set(reflect.ValueOf(wide))
	default:
		return false
	}
	return true
}
