package ast

import (
	"errors"

	"github.com/RettyInc/gqlcodegen/ast/asttype"
)

type FieldArgDef struct {
	*ast
}

func NewFieldArgDef(
	name string, typeRef *TypeRef,
	defaultValue *FieldArgDefault,
) *FieldArgDef {
	childlen := []Ast{typeRef}
	if defaultValue != nil {
		childlen = append(childlen, defaultValue)
	}
	return &FieldArgDef{
		newAst(asttype.FieldArgDef, name, childlen),
	}
}

func (fa *FieldArgDef) Type() *TypeRef {
	t, ok := fa.Children()[0].(*TypeRef)
	if ok {
		return t
	}
	panic(errors.New("type is wrong"))
}

func (fa *FieldArgDef) Default() *FieldArgDefault {
	cs := fa.Children()
	if len(cs) < 2 {
		return nil
	}
	d, ok := cs[1].(*FieldArgDefault)
	if ok {
		return d
	}
	panic(errors.New("default is wrong"))
}
