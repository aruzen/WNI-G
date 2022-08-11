package wni

import (
	"github.com/aruzen/wni-g/utils"
	"testing"
)

func TestTokenizer(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want []token
	}{
		{"struct", "a : {b : {}}", []token{"a", ":", "{", "b", ":", "{", "}", "}"}},
		{"array", "a : []", []token{"a", ":", "[", "]"}},
		{"modifier", "mod a : {}", []token{"mod", "a", ":", "{", "}"}},
		{"int", "a : 1", []token{"a", ":", "1"}},
		{"float", "a : 1.02", []token{"a", ":", "1.02"}},
		{"string", "a : \"aaa\"", []token{"a", ":", "\"", "aaa", "\""}},
		{"st", "st : {b : 10}", []token{"st", ":", "{", "b", ":", "10", "}"}},
		{"instance", "red : std::Color(255, 0, 0)", []token{"red", ":", "std", "::", "Color", "(", "255", ",", "0", ",", "0", ")"}},
		{"named&nest", "redCircle : std::Circle(color=std::Color(255, 0, 0), r=10, 10, 20)",
			[]token{"redCircle", ":", "std", "::", "Circle", "(", "color", "=", "std", "::", "Color", "(", "255", ",", "0", ",", "0", ")", ",", "r", "=", "10", ",", "10", ",", "20", ")"}},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ps := parseSession{}
			got := ps.tokenizer(test.arg)
			if len(got) != len(test.want) {
				t.Error("tokenizer got size ", len(got), ", want size", len(test.want))
			}
			for i := 0; i < len(got) && i < len(test.want); i++ {
				if got[i] != test.want[i] {
					t.Error("tokenizer got ", got[i], ", want", test.want[i], "(index ", i, ")")
				}
			}
		})
	}
}

func TestParse(t *testing.T) {
	nns := Parse(`
std :{
    class Color :{
        r : 0
        g : 0
        b : 0
        a : 1.0
    }
    class Point :{
        x : 0
        y : 0
    }
}

class Circle :{
    center : ..::std::Point()
    r      : 10
    color  : ..::std::Color()
}

redCircle : Circle(std::Point(100, 100), 100, std::Color(r=255))
`)

	got := nns.ToMap()
	utils.Dump(got)
}
