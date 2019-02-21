package ast

import (
	"errors"
	"github.com/RettyInc/gqlcodegen/ast/asttype"
)

type SchemaQueryDef struct {
	*ast
}

func NewSchemaQueryDef(typeRef *TypeRef) *SchemaQueryDef {
	return &SchemaQueryDef{
		newAst(asttype.SchemaQueryDef, "query", []Ast{typeRef}),
	}
}

func (sq *SchemaQueryDef) Type() *TypeRef {
	t, ok := sq.Children()[0].(*TypeRef)
	if ok {
		return t
	}
	panic(errors.New("type is wrong"))
}
