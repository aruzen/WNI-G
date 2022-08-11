package wni

import (
	"github.com/aruzen/wni-g/utils"
	"strconv"
	"strings"
)

type token string

type LogType int8

const (
	Critical LogType = iota
	Unpredicted
	Syntax
	Warning
)

func LogTypeToString(logType LogType) string {
	switch logType {
	case Critical:
		return "Critical"
	case Unpredicted:
		return "Unpredicted"
	case Syntax:
		return "Syntax"
	case Warning:
		return "Warning"
	}
	return ""
}

var Logger = func(t LogType, messages interface{}) {
	print(LogTypeToString(t)+" : ", messages)
}

var (
	specTokens     = []token{":", "::"}
	definedTokens  = []rune{'{', '}', '[', ']', '(', ')', ',', '='}
	voidedTokens   = []rune{';', ' ', '\n', '\t'}
	stringLiterals = []rune{'"', '\''}
)

type parseSession struct {
	tokens        []token
	indexes       []utils.Pair[tokenIndex, indexType]
	postProcesses []func(nns *Struct)
}

func (receiver *parseSession) tokenizer(text string) []token {
	runes := []rune(text)

	front_idx := -1
	var tokens []token
	for idx := 0; idx < len(runes); idx++ {
		if ':' == runes[idx] {
			if ':' == runes[idx+1] {
				if front_idx < idx-1 {
					tokens = append(tokens, token(runes[front_idx+1:idx]))
				}
				tokens = append(tokens, "::")
				idx++
				goto toknizer_runes_for_finded
			}
			if front_idx < idx-1 {
				tokens = append(tokens, token(runes[front_idx+1:idx]))
			}
			tokens = append(tokens, ":")
			goto toknizer_runes_for_finded
		}
		for _, definedToken := range definedTokens {
			if definedToken == runes[idx] {
				if front_idx < idx-1 {
					tokens = append(tokens, token(runes[front_idx+1:idx]))
				}
				tokens = append(tokens, token(runes[idx]))
				goto toknizer_runes_for_finded
			}
		}
		for _, voidedToken := range voidedTokens {
			if voidedToken == runes[idx] {
				if front_idx < idx-1 {
					tokens = append(tokens, token(runes[front_idx+1:idx]))
				}
				goto toknizer_runes_for_finded
			}
		}
		for _, stringLiteral := range stringLiterals {
			if stringLiteral != runes[idx] {
				continue
			}
			find_literal_idx := idx + 1
			for ; find_literal_idx < len(runes); find_literal_idx++ {
				if runes[find_literal_idx] == stringLiteral {
					if runes[find_literal_idx-1] == '\\' { // escape
						continue
					}
					break
				}
			}
			tokens = append(tokens, token(runes[idx]))
			tokens = append(tokens, token(runes[idx+1:find_literal_idx]))
			tokens = append(tokens, token(runes[find_literal_idx]))
			idx = find_literal_idx
			goto toknizer_runes_for_finded
		}
		continue
	toknizer_runes_for_finded:
		front_idx = idx
	}
	if front_idx+1 < len(runes) {
		tokens = append(tokens, token(runes[front_idx+1:]))
	}
	receiver.tokens = tokens
	return tokens
}

type indexType int8
type tokenIndex int
type structuralIndex int

const (
	Colon indexType = iota
	StructBegin
	StructEnd
	ArrayBegin
	ArrayEnd
	InstanceArgBegin
	InstanceArgEnd
)

func (receiver *parseSession) structuralResolution() []utils.Pair[tokenIndex, indexType] {
	var indexes []utils.Pair[tokenIndex, indexType]
	for idx := 0; idx < len(receiver.tokens); idx++ {
		if receiver.tokens[idx] == ":" {
			indexes = append(indexes, utils.Pair[tokenIndex, indexType]{tokenIndex(idx), Colon})
		} else if receiver.tokens[idx] == "{" {
			indexes = append(indexes, utils.Pair[tokenIndex, indexType]{tokenIndex(idx), StructBegin})
		} else if receiver.tokens[idx] == "}" {
			indexes = append(indexes, utils.Pair[tokenIndex, indexType]{tokenIndex(idx), StructEnd})
		} else if receiver.tokens[idx] == "[" {
			indexes = append(indexes, utils.Pair[tokenIndex, indexType]{tokenIndex(idx), ArrayBegin})
		} else if receiver.tokens[idx] == "]" {
			indexes = append(indexes, utils.Pair[tokenIndex, indexType]{tokenIndex(idx), ArrayEnd})
		} else if receiver.tokens[idx] == "(" {
			indexes = append(indexes, utils.Pair[tokenIndex, indexType]{tokenIndex(idx), InstanceArgBegin})
		} else if receiver.tokens[idx] == ")" {
			indexes = append(indexes, utils.Pair[tokenIndex, indexType]{tokenIndex(idx), InstanceArgEnd})
		}
	}
	receiver.indexes = indexes
	return indexes
}

func (receiver *parseSession) block(p *Struct, searchedStructuralIndex structuralIndex) structuralIndex {
	var searched tokenIndex
	if searchedStructuralIndex == -1 {
		searched = -1
	} else {
		searched = receiver.indexes[searchedStructuralIndex].Key
	}
	idx := searchedStructuralIndex + 1
	for ; int(idx) < len(receiver.indexes); idx++ {
		if len(receiver.indexes) <= int(idx+1) {
			if receiver.indexes[idx].Data == Colon {
				n := Node{}
				receiver.node(&n, receiver.indexes[idx].Key, searched)
				v := Value{}
				v.node = &n
				v.Data, searched = receiver.primitive(receiver.indexes[idx].Key, searched)
				p.Members = append(p.Members, v)
			}
			break
		}
		if receiver.indexes[idx].Data == Colon {
			n := Node{}
			receiver.node(&n, receiver.indexes[idx].Key, searched)
			// 次のコロンまでに挟まってるものでNodeの決定
			switch receiver.indexes[idx+1].Data {
			default:
				v := Value{}
				v.node = &n
				v.Data, searched = receiver.primitive(receiver.indexes[idx].Key, searched)
				p.Members = append(p.Members, v)
			case StructBegin:
				s := Struct{}
				s.node = &n
				idx = receiver.block(&s, idx+1)
				searched = receiver.indexes[idx].Key
				s.isClass = n.HaveModifier("class")
				p.Members = append(p.Members, s)
				continue // continueがないと、終端検索に引っかかる
			case ArrayBegin:
				v := Value{}
				v.node = &n
				v.Data, searched = receiver.array(receiver.indexes[idx+1].Key, searched)
				p.Members = append(p.Members, v)
			case InstanceArgBegin:
				v := Value{}
				v.node = &n
				v.Data, searched = receiver.instance(receiver.indexes[idx].Key, searched)
				p.Members = append(p.Members, v)
			}
		}
		// 読んだとこまで飛ばす
		/*for receiver.indexes[idx].Key <= searched {
			idx++
		}*/
		// いまだ読み込まれていない構造体終点記号がきたら処理を終える
		if receiver.indexes[idx].Data == StructEnd {
			searched = receiver.indexes[idx].Key
			break
		}
	}
	return idx
}

func (receiver *parseSession) exp(separatorIndex tokenIndex, searched tokenIndex) (interface{}, tokenIndex) {
	if receiver.tokens[separatorIndex+1] == "[" {
		return receiver.array(separatorIndex+1, separatorIndex+1)
	}
	v, ss := receiver.primitive(separatorIndex, searched)
	switch v.(type) {
	case Unknown:
		v2, ss2 := receiver.instance(separatorIndex, searched)
		if v2.classPointer != "" {
			return v2, ss2
		}
	}
	return v, ss
}

func (receiver *parseSession) node(node *Node, colonIdx tokenIndex, searched tokenIndex) {
	if !(searched+1 < colonIdx) {
		Logger(Syntax, "no name node.")
		return
	}
	node.Name = string(receiver.tokens[colonIdx-1])
	for i := searched + 1; i < colonIdx-1; i++ {
		node.Modifiers = append(node.Modifiers, string(receiver.tokens[i]))
	}
}

func (receiver *parseSession) primitive(separatorIndex tokenIndex, searched tokenIndex) (interface{}, tokenIndex) {
	for _, sr := range stringLiterals {
		if receiver.tokens[separatorIndex+1] == token(sr) {
			return string(receiver.tokens[separatorIndex+2]), separatorIndex + 3
		}
	}
	var err error = nil
	var d interface{} = nil
	if strings.Index(string(receiver.tokens[separatorIndex+1]), ".") == -1 {
		d, err = strconv.ParseInt(string(receiver.tokens[separatorIndex+1]), 0, 64)
		d = utils.FitBit(d)
	} else {
		d, err = strconv.ParseFloat(string(receiver.tokens[separatorIndex+1]), 64)
	}
	if err != nil {
		d, err = strconv.ParseBool(string(receiver.tokens[separatorIndex+1]))
	}
	if err != nil {
		d = Unknown(receiver.tokens[separatorIndex+1])
	}
	return d, separatorIndex + 1
}

func (receiver *parseSession) array(separatorIndex tokenIndex, searched tokenIndex) ([]interface{}, tokenIndex) {
	var arr []interface{}
	idx := separatorIndex
	for ; ; idx++ {
		if receiver.tokens[idx] == "[" || receiver.tokens[idx] == "," {
			td, ti := receiver.exp(idx, idx-1)
			idx = ti
			arr = append(arr, td)
		}
		if receiver.tokens[idx+1] == "]" {
			break
		}
	}
	return arr, searched
}

func (receiver *parseSession) instance(separatorIndex tokenIndex, searched tokenIndex) (*Instance, tokenIndex) {
	i := Instance{namedArgs: map[string]interface{}{}, Data: Struct{node: &Node{}}}
	idx := separatorIndex + 1
	for ; receiver.tokens[idx] != "(" && receiver.tokens[idx] != ":" && receiver.tokens[idx] != "["; idx++ {
		i.classPointer += string(receiver.tokens[idx])
	}
	if receiver.tokens[idx] != "(" {
		return &Instance{}, searched
	}
	for ; int(idx) < len(receiver.tokens); idx++ {
		if receiver.tokens[idx+1] == ")" {
			break
		}
		if receiver.tokens[idx] == "(" || receiver.tokens[idx] == "," {
			if receiver.tokens[idx+2] == "=" {
				td, ti := receiver.exp(idx+2, idx+1)
				i.namedArgs[string(receiver.tokens[idx+1])] = td
				idx = ti
			} else {
				td, ti := receiver.exp(idx, idx-1)
				idx = ti
				i.args = append(i.args, td)
			}
			if receiver.tokens[idx+1] == ")" {
				break
			}
		}
	}
	receiver.postProcesses = append(receiver.postProcesses, func(nns *Struct) {
		base := *(*i.Data.Node().Parent).Node().Parent
		e := base.Get(i.classPointer)
		c, ok := e.(Struct)
		if !ok || !c.IsClass() {
			return
		}
		i.Class = &c
		cc := i.Class.Copy().(Struct)
		cc.node = i.Data.Node()
		for idx, m := range cc.Members {
			if idx < len(i.args) {
				cc.Members[idx] = Value{node: m.Node(), Data: i.args[idx]}
			}
			d, ok := i.namedArgs[m.Node().Name]
			if !ok {
				continue
			}
			cc.Members[idx] = Value{node: m.Node(), Data: d}
		}
		i.Data = cc
	})
	return &i, idx + 1
}

func Parse(text string) *Struct {
	nns := Struct{}
	ps := parseSession{}
	ps.tokenizer(text)
	ps.structuralResolution()
	ps.block(&nns, -1)
	nns.Normalize()
	for _, pp := range ps.postProcesses {
		pp(&nns)
	}
	return &nns
}
