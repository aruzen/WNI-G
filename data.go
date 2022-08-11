package wni

import (
	"strconv"
	"strings"
)

type Instance struct {
	Class        *Struct
	classPointer string
	Data         Struct
	args         []interface{}
	namedArgs    map[string]interface{}
}

func (receiver *Instance) Normalize(p *Element) {
	receiver.Data.Node().Parent = p
	for _, a := range receiver.args {
		if i, ok := a.(*Instance); ok {
			i.Normalize(p)
		}
	}
	for _, a := range receiver.namedArgs {
		if i, ok := a.(*Instance); ok {
			i.Normalize(p)
		}
	}
}

type Unknown string

type PointerText string

type Object interface {
	Normalize()
	Get(pt string) Element // Regget と違い値を返せない
	GetByPointerTexts(pt []PointerText) Element
	Regget(pt string) interface{} // regular get 仕様上のポインタと同じ挙動
	ReggetByPointerTexts(pt []PointerText) interface{}
}
type Element interface {
	Object
	Node() *Node
	Copy() Element
}

func string2PointerTexts(s string) []PointerText {
	ss := strings.Split(s, "::")
	r := make([]PointerText, len(ss))
	for i := range ss {
		r[i] = PointerText(strings.TrimSpace(ss[i]))
	}
	return r
}

type Node struct {
	Parent    *Element
	Name      string
	Modifiers []string
}

func (receiver Node) HaveModifier(mod string) bool {
	for _, m := range receiver.Modifiers {
		if m == mod {
			return true
		}
	}
	return false
}

func (receiver Node) Copy() Node {
	r := Node{}
	r.Parent = receiver.Parent
	r.Name = receiver.Name
	r.Modifiers = make([]string, len(receiver.Modifiers))
	for i := range receiver.Modifiers {
		r.Modifiers[i] = receiver.Modifiers[i]
	}
	return r
}

type Value struct {
	node *Node
	Data interface{}
}

func (receiver Value) Node() *Node {
	return receiver.node
}

func (receiver Value) Normalize() {
	if arr, ok := receiver.Data.([]interface{}); ok {
		for _, a := range arr {
			if obj, ook := a.(*Instance); ook {
				var e Element = receiver
				obj.Data.Node().Parent = &e
				obj.Normalize(&e)
			}
			if obj, ook := a.(Object); ook {
				obj.Normalize()
			}
		}
		return
	}
	if obj, ook := receiver.Data.(*Instance); ook {
		var e Element = receiver
		obj.Data.Node().Parent = &e
		obj.Normalize(&e)
	}
	if obj, ook := receiver.Data.(Object); ook {
		obj.Normalize()
	}
}

func (receiver Value) Get(pt string) Element {
	return receiver.GetByPointerTexts(string2PointerTexts(pt))
}

func (receiver Value) GetByPointerTexts(pt []PointerText) Element {
	if pt[0] == ".." {
		return (*receiver.Node().Parent).GetByPointerTexts(pt[1:])
	}
	return receiver
}

func (receiver Value) Regget(pt string) interface{} {
	return receiver.ReggetByPointerTexts(string2PointerTexts(pt))
}

func (receiver Value) ReggetByPointerTexts(pt []PointerText) interface{} {
	if len(pt) == 0 {
		return receiver
	}
	if pt[0] == ".." {
		return (*receiver.Node().Parent).GetByPointerTexts(pt[1:])
	}
	if pt[0] == "*" {
		if len(pt) == 1 {
			return receiver.Data
		}
		if obj, ok := receiver.Data.(Object); ok {
			return obj.ReggetByPointerTexts(pt[1:])
		}
		return nil
	}
	if pt[0][0] == '[' {
		if arr, ok := receiver.Data.([]interface{}); ok {
			n, err := strconv.ParseInt(string(pt[0][1:len(pt[0])]), 0, 64)
			if err != nil || int(n) >= len(arr) {
				return nil
			}
			if len(pt) == 1 {
				return arr[n]
			}
			if obj, ook := arr[n].(Object); ook {
				return obj.ReggetByPointerTexts(pt[1:])
			}
		}
	}
	return nil
}

func (receiver Value) Copy() Element {
	r := Value{}
	n := receiver.Node().Copy()
	r.node = &n
	switch receiver.Data.(type) {
	case *Instance:
		s := receiver.Data.(*Instance)
		i := Instance{}
		i.Class = s.Class
		s.Data = s.Data.Copy().(Struct)
		r.Data = &s
	default:
		r.Data = receiver.Data
	}
	return r
}

type Struct struct {
	node    *Node
	Members []Element
	isClass bool
}

func (receiver Struct) IsClass() bool {
	return receiver.isClass
}

func (receiver Struct) Node() *Node {
	return receiver.node
}

func (receiver Struct) Normalize() {
	for _, m := range receiver.Members {
		var e Element = receiver
		m.Node().Parent = &e
		m.Normalize()
	}
}

func (receiver Struct) Get(pt string) Element {
	return receiver.GetByPointerTexts(string2PointerTexts(pt))
}

func (receiver Struct) GetByPointerTexts(pt []PointerText) Element {
	if len(pt) == 0 {
		return receiver
	}
	if pt[0] == ".." {
		return (*receiver.Node().Parent).GetByPointerTexts(pt[1:])
	}
	for _, m := range receiver.Members {
		if m.Node().Name == string(pt[0]) {
			return m.GetByPointerTexts(pt[1:])
		}
	}
	return receiver
}

func (receiver Struct) Regget(pt string) interface{} {
	return receiver.ReggetByPointerTexts(string2PointerTexts(pt))
}

func (receiver Struct) ReggetByPointerTexts(pt []PointerText) interface{} {
	if len(pt) == 0 {
		return receiver
	}
	if pt[0] == ".." {
		return (*receiver.Node().Parent).GetByPointerTexts(pt[1:])
	}
	for _, m := range receiver.Members {
		if m.Node().Name == string(pt[0]) {
			return m.GetByPointerTexts(pt[1:])
		}
	}
	return nil
}

func valueToMap(v interface{}) interface{} {
	switch v.(type) {
	case *Instance:
		return v.(*Instance).Data.ToMap()
	case []interface{}:
		t := make([]interface{}, len(v.([]interface{})))
		for idx, n := range v.([]interface{}) {
			t[idx] = valueToMap(n)
		}
		return t
	default:
		return v
	}
}

func (receiver Struct) ToMap() map[string]interface{} {
	result := map[string]interface{}{}
	for _, m := range receiver.Members {
		switch m.(type) {
		case Struct:
			result[m.Node().Name] = m.(Struct).ToMap()
		case Value:
			result[m.Node().Name] = valueToMap(m.(Value).Data)
		}
	}
	return result
}

func (receiver Struct) Copy() Element {
	r := Struct{}
	n := receiver.Node().Copy()
	r.node = &n
	r.isClass = receiver.isClass
	r.Members = make([]Element, len(receiver.Members))
	for i := range receiver.Members {
		r.Members[i] = receiver.Members[i]
	}
	return r
}

type Valuable interface {
	int | int8 | int16 | int32 | int64 | float32 | float64 | string | *Instance | Struct | Unknown
}
