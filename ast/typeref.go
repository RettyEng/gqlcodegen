package ast

import "github.com/RettyInc/gqlcodegen/ast/asttype"

type TypeRef struct {
	*ast
	isNullable bool
}

func NewTypeRef(
	typeName string,
	typeVars []Ast,
	isNullable bool,
) *TypeRef {
	return &TypeRef{
		ast:        newAst(asttype.TypeRef, typeName, typeVars),
		isNullable: isNullable,
	}
}

func (t *TypeRef) TypeVars() []*TypeRef {
	var tvs []*TypeRef
	for _, tv := range t.Children() {
		tv, ok := tv.(*TypeRef)
		if ok {
			tvs = append(tvs, tv)
		}
	}
	return tvs
}

func (t *TypeRef) IsNullable() bool {
	return t.isNullable
}
