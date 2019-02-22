package ast

import "github.com/RettyInc/gqlcodegen/ast/asttype"

type Root struct {
	*ast
}

func NewRoot(children []Ast) *Root {
	return &Root{newAst(asttype.Root, "", children)}
}

func (r *Root) Schema() *SchemaDef {
	for _, c := range r.Children() {
		if s, ok := c.(*SchemaDef); ok {
			return s
		}
	}
	return nil
}

func (r *Root) Scalars() []*ScalarDef {
	var ss []*ScalarDef
	for _, c := range r.Children() {
		if s, ok := c.(*ScalarDef); ok {
			ss = append(ss, s)
		}
	}
	return ss
}

func (r *Root) Types() []*TypeDef {
	var ts []*TypeDef
	for _, c := range r.Children() {
		if t, ok := c.(*TypeDef); ok {
			ts = append(ts, t)
		}
	}
	return ts
}

func (r *Root) Enums() []*EnumDef {
	var es []*EnumDef
	for _, c := range r.Children() {
		if e, ok := c.(*EnumDef); ok {
			es = append(es, e)
		}
	}
	return es
}
