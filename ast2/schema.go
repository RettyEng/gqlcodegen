package ast2

import "github.com/RettyInc/gqlcodegen/gql"

type DefineSchemaFieldExpression interface {
	Eval(schema *gql.Schema)
}
