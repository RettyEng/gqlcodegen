package ast2

import (
	"log"

	"github.com/RettyInc/gqlcodegen/gql"
)

type TopLevel struct {
	Expressions []DefinitionExpression
}

func (t *TopLevel) Eval() *gql.TypeSystem {
	sys := gql.NewTypeSystem()
	for _, e := range t.Expressions {
		e.Eval(sys)
	}
	return sys
}

type DefinitionExpression interface {
	Eval(system *gql.TypeSystem)
}

func evalDirectives(exp []DirectiveExpression) []*gql.DirectiveRef {
	var ret []*gql.DirectiveRef
	for _, e := range exp {
		ret = append(ret, e.Eval())
	}
	return ret
}

func evalInputValues(exp []InputValueExpression) []*gql.InputValue {
	var ret []*gql.InputValue
	for _, e := range exp {
		ret = append(ret, e.Eval())
	}
	return ret
}

type DefineSchemaExpression struct {
	DirectiveExpressions []DirectiveExpression
	Expressions          []SchemaInternalExpression
}

func (d *DefineSchemaExpression) Eval(system *gql.TypeSystem) {
	system.Schema.Directives = evalDirectives(d.DirectiveExpressions)
	for _, e := range d.Expressions {
		e.Eval(system.Schema)
	}
}

type ExtendSchemaExpression struct {
	directiveExpressions []DirectiveExpression
	expressions          []SchemaInternalExpression
}

func (e *ExtendSchemaExpression) Eval(system *gql.TypeSystem) {
	directives := evalDirectives(e.directiveExpressions)
	system.Schema.Directives = append(system.Schema.Directives, directives...)
	for _, e := range e.expressions {
		e.Eval(system.Schema)
	}
}

type DefineScalarExpression struct {
	DescriptionExpression DescriptionExpression
	NameExpression        NameExpression
	DirectiveExpressions  []DirectiveExpression
}

func (d *DefineScalarExpression) Eval(system *gql.TypeSystem) {
	system.ScalarTypes[d.NameExpression.Eval()] = &gql.Scalar{
		Description: d.DescriptionExpression.Eval(),
		Name:        d.NameExpression.Eval(),
		Directives:  evalDirectives(d.DirectiveExpressions),
	}
}

type ExtendScalarExpression struct {
	nameExpression       NameExpression
	directiveExpressions []DirectiveExpression
}

func (e *ExtendScalarExpression) Eval(system *gql.TypeSystem) {
	s, ok := system.ScalarTypes[e.nameExpression.Eval()]
	if !ok {
		log.Fatalf("extend target not found")
	}
	s.Directives =
		append(s.Directives, evalDirectives(e.directiveExpressions)...)
}

type DefineObjectExpression struct {
	DescriptionExpression DescriptionExpression
	NameExpression        NameExpression
	DirectiveExpressions  []DirectiveExpression
	ObjectExpression      []ObjectInternalExpression
}

func (d *DefineObjectExpression) Eval(system *gql.TypeSystem) {
	obj := &gql.Object{
		Description: d.DescriptionExpression.Eval(),
		Name:        d.NameExpression.Eval(),
		Directives:  evalDirectives(d.DirectiveExpressions),
	}
	for _, e := range d.ObjectExpression {
		e.Eval(obj)
	}
	system.ObjectTypes[d.NameExpression.Eval()] = obj
}

type ExtendObjectExpression struct {
	nameExpression       NameExpression
	directiveExpressions []DirectiveExpression
	objectExpression     []ObjectInternalExpression
}

func (e *ExtendObjectExpression) Eval(system *gql.TypeSystem) {
	obj, ok := system.ObjectTypes[e.nameExpression.Eval()]
	if !ok {
		log.Fatalf("extend target not found")
	}
	obj.Directives =
		append(obj.Directives, evalDirectives(e.directiveExpressions)...)
	for _, exp := range e.objectExpression {
		exp.Eval(obj)
	}
	system.ObjectTypes[e.nameExpression.Eval()] = obj
}

type DefineInterfaceExpression struct {
	DescriptionExpression DescriptionExpression
	NameExpression        NameExpression
	DirectiveExpressions  []DirectiveExpression
	InterfaceExpression   []InterfaceInternalExpression
}

func (d *DefineInterfaceExpression) Eval(system *gql.TypeSystem) {
	i := &gql.Interface{
		Description: d.DescriptionExpression.Eval(),
		Name:        d.NameExpression.Eval(),
		Directives:  evalDirectives(d.DirectiveExpressions),
	}
	for _, exp := range d.InterfaceExpression {
		exp.Eval(i)
	}
	system.InterfaceTypes[d.NameExpression.Eval()] = i
}

type ExtendInterfaceExpression struct {
	nameExpression       NameExpression
	directiveExpressions []DirectiveExpression
	interfaceExpression  []InterfaceInternalExpression
}

func (e *ExtendInterfaceExpression) Eval(system *gql.TypeSystem) {
	i, ok := system.InterfaceTypes[e.nameExpression.Eval()]
	if !ok {
		log.Fatalf("extend target not found")
	}
	i.Directives = append(i.Directives, evalDirectives(e.directiveExpressions)...)
	for _, exp := range e.interfaceExpression {
		exp.Eval(i)
	}
	system.InterfaceTypes[e.nameExpression.Eval()] = i
}

type DefineUnionExpression struct {
	DescriptionExpression DescriptionExpression
	NameExpression        NameExpression
	DirectiveExpressions  []DirectiveExpression
	UnionExpression       []UnionInternalExpression
}

func (d *DefineUnionExpression) Eval(system *gql.TypeSystem) {
	u := &gql.Union{
		Description: d.DescriptionExpression.Eval(),
		Name:        d.NameExpression.Eval(),
		Directives:  evalDirectives(d.DirectiveExpressions),
	}
	for _, e := range d.UnionExpression {
		e.Eval(u)
	}
	system.UnionTypes[d.NameExpression.Eval()] = u
}

type ExtendUnionExpression struct {
	nameExpression       NameExpression
	directiveExpressions []DirectiveExpression
	unionExpression      []UnionInternalExpression
}

func (e *ExtendUnionExpression) Eval(system *gql.TypeSystem) {
	u, ok := system.UnionTypes[e.nameExpression.Eval()]
	if !ok {
		log.Fatalf("extend target not found")
	}
	u.Directives = append(u.Directives, evalDirectives(e.directiveExpressions)...)
	for _, e := range e.unionExpression {
		e.Eval(u)
	}
	system.UnionTypes[e.nameExpression.Eval()] = u
}

type DefineEnumExpression struct {
	DescriptionExpression DescriptionExpression
	NameExpression        NameExpression
	DirectiveExpressions  []DirectiveExpression
	EnumExpression        []EnumInternalExpression
}

func (d *DefineEnumExpression) Eval(system *gql.TypeSystem) {
	enum := &gql.Enum{
		Description: d.DescriptionExpression.Eval(),
		Name:        d.NameExpression.Eval(),
		Directives:  evalDirectives(d.DirectiveExpressions),
	}
	for _, e := range d.EnumExpression {
		e.Eval(enum)
	}
	system.EnumTypes[d.NameExpression.Eval()] = enum
}

type ExtendEnumExpression struct {
	nameExpression       NameExpression
	directiveExpressions []DirectiveExpression
	enumExpression       []EnumInternalExpression
}

func (e *ExtendEnumExpression) Eval(system *gql.TypeSystem) {
	enum, ok := system.EnumTypes[e.nameExpression.Eval()]
	if !ok {
		log.Fatalf("extend target not found")
	}
	enum.Directives = append(enum.Directives, evalDirectives(e.directiveExpressions)...)
	for _, e := range e.enumExpression {
		e.Eval(enum)
	}
	system.EnumTypes[e.nameExpression.Eval()] = enum
}

type DefineInputObjectExpression struct {
	DescriptionExpression             DescriptionExpression
	NameExpression                    NameExpression
	DirectiveExpressions              []DirectiveExpression
	DefineInputObjectFieldExpressions []InputValueExpression
}

func (d *DefineInputObjectExpression) Eval(system *gql.TypeSystem) {
	obj := &gql.InputObject{
		Description: d.DescriptionExpression.Eval(),
		Name:        d.NameExpression.Eval(),
		Directives:  evalDirectives(d.DirectiveExpressions),
		InputValue:  evalInputValues(d.DefineInputObjectFieldExpressions),
	}
	system.InputObjectTypes[d.NameExpression.Eval()] = obj
}

type ExtendInputObjectExpression struct {
	nameExpression                    NameExpression
	directiveExpressions              []DirectiveExpression
	defineInputObjectFieldExpressions []InputValueExpression
}

func (e *ExtendInputObjectExpression) Eval(system *gql.TypeSystem) {
	obj, ok := system.InputObjectTypes[e.nameExpression.Eval()]
	if !ok {
		log.Fatalf("extend target not found")
	}
	obj.Directives = append(obj.Directives, evalDirectives(e.directiveExpressions)...)
	obj.InputValue = append(obj.InputValue, evalInputValues(e.defineInputObjectFieldExpressions)...)
	system.InputObjectTypes[e.nameExpression.Eval()] = obj
}

type DirectiveDefinition struct {
	DescriptionExpression DescriptionExpression
	NameExpression        NameExpression
	ArgsExpression        []InputValueExpression
	Expressions           []DirectiveInternalExpression
}

func (d *DirectiveDefinition) Eval(system *gql.TypeSystem) {
	directive := &gql.Directive{
		Description: d.DescriptionExpression.Eval(),
		Name:        d.NameExpression.Eval(),
		Arguments:   evalInputValues(d.ArgsExpression),
	}
	for _, e := range d.Expressions {
		e.Eval(directive)
	}
	system.Directives[d.NameExpression.Eval()] = directive
}
