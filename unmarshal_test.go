package wni

import (
	"github.com/aruzen/wni-g/utils"
	"testing"
)

func TestMarshal(t *testing.T) {
	type args struct {
		text string
		obj  interface{}
	}

	type Normal struct {
		A string
		B int
	}
	type Nest struct {
		A string
		B Normal
	}
	type Interface struct {
		A string
		B interface{}
	}
	type Pointer struct {
		A string
		B *Normal
	}
	type Slice struct {
		A []string
		B []float64
	}
	type Instance struct {
		A Normal
	}
	type InstanceSlice struct {
		A []Normal
	}
	type Named struct {
		A int `wni:"apple" json:"atugi"`
		B int `wni:"banana" json:"bankoku"`
	}
	type ToLower struct {
		CPlusPlus int
		GoLang    float32
	}

	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			"normal",
			args{"A : 'aaa' B : 10", Normal{}},
			Normal{"aaa", 10},
		},
		{
			"nest",
			args{"A : 'nest' B : { A : 'aaa' B : 10 }", Nest{}},
			Nest{"nest", Normal{"aaa", 10}},
		},
		{
			"interface",
			args{"A : 'nest' B : { A : 'aaa' B : 10 }", Interface{B: Normal{}}},
			Interface{"nest", Normal{"aaa", 10}},
		},
		{
			"pointer",
			args{"A : 'pointer' B : { A : 'aaa' B : 10 }", Pointer{B: &Normal{}}},
			Pointer{"pointer", &Normal{"aaa", 10}},
		},
		{
			"slice",
			args{"A : ['one', 'tow', 'three'] b : [1.0, 2.0, 3.0]", Slice{}},
			Slice{[]string{"one", "tow", "three"}, []float64{1.0, 2.0, 3.0}},
		},
		{
			"instance",
			args{`
class Normal :{
	a : "a"
	b : 1
}
a : Normal(b=2)
`, Instance{}},
			Instance{Normal{"a", 2}},
		},
		{
			"instance & slice",
			args{`
class Normal :{
	a : "a"
	b : 1
}
a : [Normal("one"), Normal("two", 2), Normal("three", 3)]
`, InstanceSlice{}},
			InstanceSlice{[]Normal{{"one", 1}, {"two", 2}, {"three", 3}}},
		},
		{
			"named_wni",
			args{"apple : 10 banana : 20", Named{}},
			Named{10, 20},
		},
		{
			"named_json",
			args{"atugi : 10 bankoku : 20", Named{}},
			Named{10, 20},
		},
		{
			"to_lower_head",
			args{"a : 10 b : 20", Named{}},
			Named{10, 20},
		},
		{
			"to_lower_full",
			args{"c_plus_plus : 23 go_lang : 19.0", ToLower{}},
			ToLower{23, 19.0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Unmarshal(tt.args.text, &tt.args.obj)
			if !utils.Compare(tt.args.obj, tt.want) {
				t.Errorf("not equal.")
			}
		})
	}
}

func Test_toLower(t *testing.T) {
	tests := []struct {
		name string
		arg  string
		want string
	}{
		{"head", "Apple", "apple"},
		{"full", "GoLang", "go_lang"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := toLower(tt.arg); got != tt.want {
				t.Errorf("toLower() = %v, want %v", got, tt.want)
			}
		})
	}
}
