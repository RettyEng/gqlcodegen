package ast2

import (
	"log"

	"github.com/RettyInc/gqlcodegen/gql"
)


type TopLevel struct {
	expressions []DefinitionExpression
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


type DefineSchemaExpression struct {
	directiveExpressions []DirectiveExpression
	expressions []DefineSchemaFieldExpression
}
func (d *DefineSchemaExpression) Eval(system *gql.TypeSystem) {
	system.Schema.Directives = evalDirectives(d.directiveExpressions)
	for _, e := range d.expressions {
		e.Eval(system.Schema)
	}
}


type ExtendSchemaExpression struct {
	directiveExpressions []DirectiveExpression
	expressions []DefineSchemaFieldExpression
}
func (e *ExtendSchemaExpression) Eval(system *gql.TypeSystem) {
	directives := evalDirectives(e.directiveExpressions)
	system.Schema.Directives = append(system.Schema.Directives, directives...)
	for _, e := range e.expressions {
		e.Eval(system.Schema)
	}
}


type DefineScalarExpression struct {
	descriptionExpression DescriptionExpression
	nameExpression        NameExpression
	directiveExpressions  []DirectiveExpression
}
func (d *DefineScalarExpression) Eval(system *gql.TypeSystem) {
	system.ScalarTypes[d.nameExpression.Eval()] = &gql.Scalar{
		Description: d.descriptionExpression.Eval(),
		Name: d.nameExpression.Eval(),
		Directives: evalDirectives(d.directiveExpressions),
	}
}


type ExtendScalarExpression struct {
	nameExpression        NameExpression
	directiveExpressions  []DirectiveExpression
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
	descriptionExpression  DescriptionExpression
	nameExpression         NameExpression
	implementsExpressions  []ImplementsExpression
	directiveExpressions   []DirectiveExpression
	defineFieldExpressions []DefineFieldExpression
}
func (d *DefineObjectExpression) Eval(system *gql.TypeSystem) {
	obj := &gql.Object{
		Description: d.descriptionExpression.Eval(),
		Name: d.nameExpression.Eval(),
		Directives:evalDirectives(d.directiveExpressions),
	}
	for _, e := range d.implementsExpressions {
		e.Eval(obj)
	}
	for _, e := range d.defineFieldExpressions {
		e.Eval(obj)
	}
	system.ObjectTypes[d.nameExpression.Eval()] = obj
}


type ExtendObjectExpression struct {
	nameExpression         NameExpression
	implementsExpressions  []ImplementsExpression
	directiveExpressions   []DirectiveExpression
	defineFieldExpressions []DefineFieldExpression
}
func (e *ExtendObjectExpression) Eval(system *gql.TypeSystem) {
	obj, ok := system.ObjectTypes[e.nameExpression.Eval()]
	if !ok {
		log.Fatalf("extend target not found")
	}
	obj.Directives =
	 	append(obj.Directives, evalDirectives(e.directiveExpressions)...)
	for _, exp := range e.implementsExpressions {
		exp.Eval(obj)
	}
	for _, exp := range e.defineFieldExpressions {
		exp.Eval(obj)
	}
	system.ObjectTypes[e.nameExpression.Eval()] = obj
}


type DefineInterfaceExpression struct {
	descriptionExpression  DescriptionExpression
	nameExpression         NameExpression
	directiveExpressions   []DirectiveExpression
	defineFieldExpressions []DefineInterfaceFieldExpression
}
func (d *DefineInterfaceExpression) Eval(system *gql.TypeSystem) {
	i := &gql.Interface{
		Description: d.descriptionExpression.Eval(),
		Name: d.nameExpression.Eval(),
		Directives: evalDirectives(d.directiveExpressions),
	}
	for _, exp := range d.defineFieldExpressions {
		exp.Eval(i)
	}
	system.InterfaceTypes[d.nameExpression.Eval()] = i
}


type ExtendInterfaceExpression struct {
	nameExpression         NameExpression
	directiveExpressions   []DirectiveExpression
	defineFieldExpressions []DefineInterfaceFieldExpression
}
func (e *ExtendInterfaceExpression) Eval(system *gql.TypeSystem) {
	i, ok := system.InterfaceTypes[e.nameExpression.Eval()]
	if !ok {
		log.Fatalf("extend target not found")
	}
	i.Directives = append(i.Directives, evalDirectives(e.directiveExpressions)...)
	for _, exp := range e.defineFieldExpressions {
		exp.Eval(i)
	}
	system.InterfaceTypes[e.nameExpression.Eval()] = i
}


type DefineUnionExpression struct {
	descriptionExpression       DescriptionExpression
	nameExpression              NameExpression
	directiveExpressions        []DirectiveExpression
	defineUnionMemberExpression []DefineUnionMemberExpression
}
func (d *DefineUnionExpression) Eval(system *gql.TypeSystem) {
	u := &gql.Union{
		Description:d.descriptionExpression.Eval(),
		Name: d.nameExpression.Eval(),
		Directives: evalDirectives(d.directiveExpressions),
	}
	for _, e := range d.defineUnionMemberExpression {
		e.Eval(u)
	}
	system.UnionTypes[d.nameExpression.Eval()] = u
}


type ExtendUnionExpression struct {
	nameExpression              NameExpression
	directiveExpressions        []DirectiveExpression
	defineUnionMemberExpression []DefineUnionMemberExpression
}
func (e *ExtendUnionExpression) Eval(system *gql.TypeSystem) {
	u, ok := system.UnionTypes[e.nameExpression.Eval()]
	if !ok {
		log.Fatalf("extend target not found")
	}
	u.Directives = append(u.Directives, evalDirectives(e.directiveExpressions)...)
	for _, e := range e.defineUnionMemberExpression {
		e.Eval(u)
	}
	system.UnionTypes[e.nameExpression.Eval()] = u
}


type DefineEnumExpression struct {
	descriptionExpression       DescriptionExpression
	nameExpression              NameExpression
	directiveExpressions        []DirectiveExpression
	defineEnumValuesExpressions []DefineEnumValuesExpression
}
func (d *DefineEnumExpression) Eval(system *gql.TypeSystem) {
	enum := &gql.Enum{
		Description:d.descriptionExpression.Eval(),
		Name: d.nameExpression.Eval(),
		Directives: evalDirectives(d.directiveExpressions),
	}
	for _, e := range d.defineEnumValuesExpressions {
		e.Eval(enum)
	}
	system.EnumTypes[d.nameExpression.Eval()] = enum
}


type ExtendEnumExpression struct {
	nameExpression              NameExpression
	directiveExpressions        []DirectiveExpression
	defineEnumValuesExpressions []DefineEnumValuesExpression
}
func (e *ExtendEnumExpression) Eval(system *gql.TypeSystem) {
	enum, ok := system.EnumTypes[e.nameExpression.Eval()]
	if !ok {
		log.Fatalf("extend target not found")
	}
	enum.Directives = append(enum.Directives, evalDirectives(e.directiveExpressions)...)
	for _, e := range e.defineEnumValuesExpressions {
		e.Eval(enum)
	}
	system.EnumTypes[e.nameExpression.Eval()] = enum
}


type DefineInputObjectExpression struct {
	descriptionExpression             DescriptionExpression
	nameExpression                    NameExpression
	directiveExpressions              []DirectiveExpression
	defineInputObjectFieldExpressions []DefineInputObjectFieldExpression
}
func (d *DefineInputObjectExpression) Eval(system *gql.TypeSystem) {
	obj := &gql.InputObject{
		Description:d.descriptionExpression.Eval(),
		Name: d.nameExpression.Eval(),
		Directives: evalDirectives(d.directiveExpressions),
	}
	for _, e := range d.defineInputObjectFieldExpressions {
		e.Eval(obj)
	}
	system.InputObjectTypes[d.nameExpression.Eval()] = obj
}


type ExtendInputObjectExpression struct {
	nameExpression                    NameExpression
	directiveExpressions              []DirectiveExpression
	defineInputObjectFieldExpressions []DefineInputObjectFieldExpression
}
func (e *ExtendInputObjectExpression) Eval(system *gql.TypeSystem) {
	obj, ok := system.InputObjectTypes[e.nameExpression.Eval()]
	if !ok {
		log.Fatalf("extend target not found")
	}
	obj.Directives = append(obj.Directives, evalDirectives(e.directiveExpressions)...)
	for _, e := range e.defineInputObjectFieldExpressions {
		e.Eval(obj)
	}
	system.InputObjectTypes[e.nameExpression.Eval()] = obj
}


type DirectiveDefinition struct {
	descriptionExpression             DescriptionExpression
	nameExpression                    NameExpression
	defineArgsExpressions []DefineDirectiveArgsExpression
	locationExpressions []DefineDirectiveLocationExpression
}
func (d *DirectiveDefinition) Eval(system *gql.TypeSystem) {
	directive := &gql.Directive{
		Description:d.descriptionExpression.Eval(),
		Name: d.nameExpression.Eval(),
	}
	for _, e := range d.defineArgsExpressions {
		e.Eval(directive)
	}
	for _, e := range d.locationExpressions {
		e.Eval(directive)
	}
	system.Directives[d.nameExpression.Eval()] = directive
}
