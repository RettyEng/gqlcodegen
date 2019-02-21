package ast

import (
	"errors"
	"github.com/RettyInc/gqlcodegen/ast/asttype"
)

type SchemaDef struct {
	*ast
}

func NewSchemaDef(query *SchemaQueryDef) *SchemaDef {
	return &SchemaDef{
		newAst(asttype.SchemaDef, "schema", []Ast{query}),
	}
}

func (s *SchemaDef) Query() *SchemaQueryDef {
	q, ok := s.Children()[0].(*SchemaQueryDef)
	if ok {
		return q
	}
	panic(errors.New("query is wrong"))
}
