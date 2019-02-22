package ast

import "github.com/RettyInc/gqlcodegen/ast/asttype"

type Ast interface {
	Name() string
	AstType() asttype.AstType
	Children() []Ast
}

type ast struct {
	name     string
	astType  asttype.AstType
	children []Ast
}

func newAst(
	t asttype.AstType, name string, children []Ast,
) *ast {
	return &ast{name: name, astType: t, children: children}
}

func (a *ast) AstType() asttype.AstType {
	return a.astType
}

func (a *ast) Children() []Ast {
	return a.children
}

func (a *ast) Name() string {
	return a.name
}
