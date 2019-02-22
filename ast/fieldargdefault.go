package ast

import "github.com/RettyInc/gqlcodegen/ast/asttype"

type FieldArgDefault struct {
	*ast
}

func NewFieldArgDefault(value string) *FieldArgDefault {
	return &FieldArgDefault{
		newAst(asttype.FieldArgDefault, value, nil),
	}
}

func (ad *FieldArgDefault) Value() string {
	return ad.Name()
}
