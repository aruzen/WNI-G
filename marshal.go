package wni

import (
	"fmt"
	"github.com/aruzen/wni-g/utils"
	"strconv"
)

func MapToWNI(d map[string]interface{}) {
	idx := 0
	length := len(d)
	for k, v := range d {
		printInterface(k, v, []bool{}, idx != length-1)
	}
}

func printInterface(k string, v interface{}, nested []bool, nest bool) {
	switch v.(type) {
	case map[string]interface{}:
		print(nestedToString(nested), k)
		fmt.Println()
		nn := append(nested, nest)
		subMapDump(v.(map[string]interface{}), nn)
	case []interface{}:
		print(nestedToString(nested), k)
		fmt.Println()
		nn := append(nested, nest)
		subListDump(v.([]interface{}), nn)
	case string:
		print(nestedToString(nested), k, " string:", v.(string))
		fmt.Println()
	case int, int8, int16, int32, int64:
		print(nestedToString(nested), k, " int:", utils.WidenInt(v))
		fmt.Println()
	case float32, float64:
		print(nestedToString(nested), k, " float:", utils.WidenFloat(v))
		fmt.Println()
	}
}

func nestedToString(nested []bool) string {
	r := ""
	for _, n := range nested {
		if n {
			r += " |"
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
	}
}

func subListDump(d []interface{}, nested []bool) {
	length := len(d)
	idx := 0
	for i, v := range d {
		printInterface("-["+strconv.Itoa(i)+"]", v, nested, idx != length-1)
	}
}
