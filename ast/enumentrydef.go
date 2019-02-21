package ast

import "github.com/RettyInc/gqlcodegen/ast/asttype"

type EnumEntryDef struct {
	*ast
}

func NewEnumEntryDef(name string) *EnumEntryDef {
	return &EnumEntryDef{newAst(asttype.EnumEntryDef, name, nil)}
}
