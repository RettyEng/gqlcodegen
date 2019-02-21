package ast

import (
	"errors"
	"github.com/RettyInc/gqlcodegen/ast/asttype"
)

type TypeDef struct {
	*ast
}

func NewTypeDef(name string, fields []Ast) *TypeDef {
	return &TypeDef{
		newAst(asttype.TypeDef, name, fields),
	}
}

func (t *TypeDef) Fields() []*TypeFieldDef {
	cs := t.Children()
	ret := make([]*TypeFieldDef, len(cs))
	for i, c := range cs {
		f, ok := c.(*TypeFieldDef)
		if !ok {
			panic(errors.New("field is wrong"))
		}
		ret[i] = f
	}
	return ret
}
