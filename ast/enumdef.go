package ast

import "github.com/RettyInc/gqlcodegen/ast/asttype"

type EnumDef struct {
	*ast
}

func NewEnumDef(name string, entries []Ast) *EnumDef {
	return &EnumDef{newAst(asttype.EnumDef, name, entries)}
}

func (e *EnumDef) Entries() []*EnumEntryDef {
	var es []*EnumEntryDef
	for _, e := range e.Children() {
		e, ok := e.(*EnumEntryDef)
		if ok {
			es = append(es, e)
		}
	}
	return es
}