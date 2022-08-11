package wni

import (
	"github.com/aruzen/wni-g/utils"
	"reflect"
	"unicode"
)

func Unmarshal(text string, obj *interface{}) {
	parsed := Parse(text).ToMap()
	AssignByMap(parsed, obj)
}

func AssignByMap(a map[string]interface{}, obj *interface{}) {
	v := reflect.ValueOf(obj).Elem()
	c := reflect.New(v.Elem().Type()).Elem()
	assignStruct(c, v.Elem(), a)
	v.Set(c)
}

func toLower(s string) string {
	runes := []rune(s)
	for idx, r := range runes {
		if unicode.IsUpper(r) {
			if idx == 0 {
				runes[idx] = unicode.ToLower(r)
			} else {
				runes = append(runes[:idx], runes[idx-1:]...)
				runes[idx] = '_'
				runes[idx+1] = unicode.ToLower(r)
			}
		}
	}
	return string(runes)
}

func assignStruct(accepter, original reflect.Value, assigned map[string]interface{}) {
	nameMap := map[string]string{}
	for i := 0; i < original.NumField(); i++ {
		nameMap[original.Type().Field(i).Name] = original.Type().Field(i).Name
		nameMap[toLower(original.Type().Field(i).Name)] = original.Type().Field(i).Name
	}
	for i := 0; i < original.NumField(); i++ {
		nameMap[original.Type().Field(i).Tag.Get("json")] = original.Type().Field(i).Name
	}
	for i := 0; i < original.NumField(); i++ {
		nameMap[original.Type().Field(i).Tag.Get("wni")] = original.Type().Field(i).Name
	}
	for k, _v := range assigned {
		if nameMap[k] == "" {
			continue
		}
		o := original.FieldByName(nameMap[k])
		a := accepter.FieldByName(nameMap[k])
		oc := o
		var kindNest []reflect.Kind
		for oc.Kind() == reflect.Pointer || oc.Kind() == reflect.Interface {
			kindNest = append(kindNest, oc.Kind())
			oc = oc.Elem()
		}
		c := reflect.New(oc.Type()).Elem()
		assign(c, oc, _v)
		for i := len(kindNest) - 1; 0 <= i; i-- {
			switch kindNest[i] {
			case reflect.Pointer:
				c = c.Addr()
			case reflect.Interface:
				//					var in interface{} = c.Interface()
				//					c = reflect.ValueOf(in)
			}
		}
		a.Set(c)
	}
}

func assign(c, originalC reflect.Value, _v interface{}) {
	switch _v.(type) {
	case map[string]interface{}:
		assignStruct(c, originalC, _v.(map[string]interface{}))
	case []interface{}:
		if !c.CanSet() {
			return
		}
		if c.Kind() != reflect.Slice {
			return
		}
		cc := reflect.MakeSlice(c.Type(), len(_v.([]interface{})), len(_v.([]interface{})))
		for idx, i := range _v.([]interface{}) {
			if originalC.Len() == 0 {
				assign(cc.Index(idx), cc.Index(0), i)
			} else if idx < originalC.Len() {
				assign(cc.Index(idx), originalC.Index(idx), i)
			} else {
				assign(cc.Index(idx), originalC.Index(originalC.Len()-1), i)
			}
		}
		c.Set(cc)
	default:
		if !c.CanSet() {
			return
		}
		v := reflect.ValueOf(_v)
		if !c.Type().AssignableTo(v.Type()) {
			if utils.AssignIntValue(_v, c) {
				return
			}
			if utils.AssignFloatValue(_v, c) {
				return
			}
			return
		}
		c.Set(v)
	}
}
