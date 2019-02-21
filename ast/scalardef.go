package ast

import "github.com/RettyInc/gqlcodegen/ast/asttype"

type ScalarDef struct {
	*ast
}

func NewScalarDef(name string) *ScalarDef {
	return &ScalarDef{newAst(asttype.ScalarDef, name, nil)}
}
