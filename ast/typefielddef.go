package ast

import (
	"errors"
	"github.com/RettyInc/gqlcodegen/ast/asttype"
)

type TypeFieldDef struct {
	*ast
}

func NewTypeFieldDef(
	name string, fieldType *TypeRef, args ...*FieldArgDef,
) *TypeFieldDef {
	children := []Ast{fieldType}
	for _, a := range args {
		children = append(children, a)
	}
	return &TypeFieldDef{
		newAst(asttype.TypeFieldDef, name, children),
	}
}

func (tf *TypeFieldDef) Type() *TypeRef {
	t, ok := tf.Children()[0].(*TypeRef)
	if ok {
		return t
	}
	panic(errors.New("type is wrong"))
}

func (tf *TypeFieldDef) Args() []*FieldArgDef {
	cs := tf.Children()
	if len(cs) < 2 {
		return nil
	}
	var as []*FieldArgDef
	for _, c := range cs[1:] {
		if c, ok := c.(*FieldArgDef); ok {
			as = append(as, c)
		}
	}
	return as
}
