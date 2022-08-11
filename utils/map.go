package utils

import (
	"fmt"
	"strconv"
)

func Dump(d map[string]interface{}) {
	for k, v := range d {
		printInterface(k, v, []bool{}, true)
	}
}

func printInterface(k string, v interface{}, nested []bool, nest bool) {
	switch v.(type) {
	case map[string]interface{}:
		print(nestedToString(nested, nest), k)
		fmt.Println()
		if 0 < len(nested) {
			nested[len(nested)-1] = nest
		}
		nn := append(nested, true)
		subMapDump(v.(map[string]interface{}), nn)
	case []interface{}:
		print(nestedToString(nested, nest), k)
		fmt.Println()
		if 0 < len(nested) {
			nested[len(nested)-1] = nest
		}
		nn := append(nested, true)
		subListDump(v.([]interface{}), nn)
	case string:
		print(nestedToString(nested, nest), k, " string:", v.(string))
		fmt.Println()
	case int, int8, int16, int32, int64:
		print(nestedToString(nested, nest), k, " int:", WidenInt(v))
		fmt.Println()
	case float32, float64:
		print(nestedToString(nested, nest), k, " float:", WidenFloat(v))
		fmt.Println()
	}
}

func nestedToString(nested []bool, nest bool) string {
	r := ""
	for i, n := range nested {
		if i == len(nested)-1 && !nest {
			r += " └"
		} else if n && i == len(nested)-1 {
			r += " ├"
		} else if n {
			r += " │"
		} else {
			r += "  "
		}
	}
	return r
}

func subMapDump(d map[string]interface{}, nested []bool) {
	length := len(d)
	idx := 0
	for k, v := range d {
		printInterface("-"+k, v, nested, idx != length-1)
		idx++
	}
}

func subListDump(d []interface{}, nested []bool) {
	length := len(d)
	idx := 0
	for i, v := range d {
		printInterface("-["+strconv.Itoa(i)+"]", v, nested, idx != length-1)
		idx++
	}
}
