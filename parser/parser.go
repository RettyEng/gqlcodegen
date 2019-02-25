package parser

import (
	"errors"
	"io"
	"log"
	"strings"

	"github.com/RettyInc/gqlcodegen/ast/directive"

	"github.com/RettyInc/gqlcodegen/ast"
	"github.com/RettyInc/gqlcodegen/gql"
	"github.com/RettyInc/gqlcodegen/lexer"
	"github.com/RettyInc/gqlcodegen/lexer/token"
)

type Parser struct {
	lexer *lexer.Lexer
	ast   *ast.TopLevel
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		lexer: lexer.NewLexer(reader),
	}
}

func (p *Parser) ParseSchema() *ast.TopLevel {
	if p.ast != nil {
		return p.ast
	}
	var exp []ast.DefinitionExpression
	for {
		if !p.hasNext() {
			break
		}

		t := p.prefetch(0)
		if t.Type() == token.TypeStrVal {
			t = p.prefetch(1)
		}

		switch t.Value() {
		case "schema":
			exp = append(exp, p.parseSchema())
			continue
		case "scalar":
			exp = append(exp, p.parseScalar())
			continue
		case "enum":
			exp = append(exp, p.parseEnum())
			continue
		case "type":
			exp = append(exp, p.parseType())
		case "interface":
			exp = append(exp, p.parseInterface())
		case "union":
			exp = append(exp, p.parseUnion())
		case "directive":
			exp = append(exp, p.parseDirective())
		case "input":
			exp = append(exp, p.parseInput())
		case "extend":
			exp = append(exp, p.parseExtend())
		default:
			unexpectedToken(t)
		}
	}
	p.ast = &ast.TopLevel{exp}
	return p.ast
}

func (p *Parser) ParseAndEvalSchema() *gql.TypeSystem {
	return p.ParseSchema().Eval()
}

func (p *Parser) parseExtend() ast.DefinitionExpression {
	validateTokenValue(p.prefetch(0), "extend")
	switch p.prefetch(1).Value() {
	case "schema":
		return p.parseExtendSchema()
	case "scalar":
		return p.parseExtendScalar()
	case "enum":
		return p.parseExtendEnum()
	case "type":
		return p.parseExtendObject()
	case "interface":
		return p.parseExtendInterface()
	case "union":
		return p.parseExtendUnion()
	case "input":
		return p.parseExtendInput()
	default:
		unexpectedToken(p.prefetch(1))
	}
	return nil
}

func (p *Parser) parseExtendInput() ast.DefinitionExpression {
	validateTokenValue(p.pop(), "extend")
	validateTokenValue(p.pop(), "input")
	n := p.parseName()
	direc := p.parseDirectivesOrEmpty()
	var body []ast.InputValueExpression
	if p.preValueCheck(0, "{") {
		_ = p.pop()
		for !p.preValueCheck(0, "}") {
			body = append(body, p.parseInputValue())
		}
	}
	return &ast.ExtendInputObjectExpression{
		n, direc, body,
	}
}

func (p *Parser) parseExtendUnion() ast.DefinitionExpression {
	validateTokenValue(p.pop(), "extend")
	validateTokenValue(p.pop(), "union")
	n := p.parseName()
	direc := p.parseDirectivesOrEmpty()
	var body []ast.UnionInternalExpression
	if p.preValueCheck(0, "=") {
		body = p.parseUnionBody()
	}
	return &ast.ExtendUnionExpression{n, direc, body}
}

func (p *Parser) parseExtendInterface() ast.DefinitionExpression {
	validateTokenValue(p.pop(), "extend")
	validateTokenValue(p.pop(), "interface")
	n := p.parseName()
	direc := p.parseDirectivesOrEmpty()
	var body []ast.InterfaceInternalExpression
	if p.preValueCheck(0, "{") {
		body = p.parseInterfaceBody()
	}
	return &ast.ExtendInterfaceExpression{
		n, direc, body,
	}
}

func (p *Parser) parseExtendSchema() ast.DefinitionExpression {
	validateTokenValue(p.pop(), "extend")
	validateTokenValue(p.pop(), "schema")
	direc := p.parseDirectivesOrEmpty()
	var body []ast.SchemaInternalExpression
	if p.preValueCheck(0, "{") {
		body = p.parseSchemaBody()
	}
	return &ast.ExtendSchemaExpression{direc, body}
}

func (p *Parser) parseExtendScalar() ast.DefinitionExpression {
	validateTokenValue(p.pop(), "extend")
	validateTokenValue(p.pop(), "scalar")
	n := p.parseName()
	d := p.parseDirectives()
	return &ast.ExtendScalarExpression{n, d}
}

func (p *Parser) parseExtendEnum() ast.DefinitionExpression {
	validateTokenValue(p.pop(), "extend")
	validateTokenValue(p.pop(), "enum")
	n := p.parseName()
	d := p.parseDirectivesOrEmpty()
	var body []ast.EnumInternalExpression
	if p.preValueCheck(0, "{") {
		body = p.parseEnumBody()
	}
	return &ast.ExtendEnumExpression{
		n, d, body,
	}
}

func (p *Parser) parseExtendObject() ast.DefinitionExpression {
	validateTokenValue(p.pop(), "extend")
	validateTokenValue(p.pop(), "type")
	n := p.parseName()
	exp := p.parseImplementsOrEmpty()
	d := p.parseDirectivesOrEmpty()
	if p.preValueCheck(0, "{") {
		exp = append(exp, p.parseObjectBody()...)
	}
	return &ast.ExtendObjectExpression{n, d, exp}
}

func (p *Parser) parseInput() ast.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "input")
	n := p.parseName()
	direc := p.parseDirectivesOrEmpty()
	t = p.pop()
	validateTokenValue(t, "{")
	var args []ast.InputValueExpression
	for !p.preValueCheck(0, "}") {
		args = append(args, p.parseInputValue())
	}
	_ = p.pop()
	return &ast.DefineInputObjectExpression{
		DescriptionExpression:             desc,
		NameExpression:                    n,
		DirectiveExpressions:              direc,
		DefineInputObjectFieldExpressions: args,
	}
}

func (p *Parser) parseDirective() ast.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "directive")
	t = p.pop()
	validateTokenValue(t, "@")
	n := p.parseName()
	args := p.parseDirectiveArgsDefinition()
	t = p.pop()
	validateTokenValue(t, "on")
	if p.preValueCheck(0, "|") {
		_ = p.pop()
	}
	var locs []ast.DirectiveInternalExpression
	locs = append(locs, p.parseDirectiveLocation())
	for p.preValueCheck(0, "|") {
		_ = p.pop()
		locs = append(locs, p.parseDirectiveLocation())
	}
	return &ast.DirectiveDefinition{
		DescriptionExpression: desc,
		NameExpression:        n,
		ArgsExpression:        args,
		Expressions:           locs,
	}
}

func (p *Parser) parseDirectiveLocation() ast.DirectiveInternalExpression {
	t := p.pop()
	validateTokenType(t, token.TypeName)
	for i := 0; !strings.HasPrefix(directive.Location(i).String(), "Location("); i++ {
		if t.Value() == directive.Location(i).String() {
			return &ast.DefineDirectiveLocationExpression{directive.Location(i)}
		}
	}
	l, c := t.LineCol()
	log.Fatalf("unknown directive location %s at line %d, col %d", t.Value(), l, c)
	return nil
}

func (p *Parser) parseDirectiveArgsDefinition() []ast.InputValueExpression {
	t := p.pop()
	validateTokenValue(t, "(")
	var args []ast.InputValueExpression
	for !p.preValueCheck(0, ")") {
		args = append(args, p.parseInputValue())
	}
	_ = p.pop()
	return args
}

func (p *Parser) parseUnion() ast.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "union")
	n := p.parseName()
	directives := p.parseDirectivesOrEmpty()
	body := p.parseUnionBody()
	return &ast.DefineUnionExpression{
		DescriptionExpression: desc,
		NameExpression:        n,
		DirectiveExpressions:  directives,
		UnionExpression:       body,
	}
}

func (p *Parser) parseUnionBody() []ast.UnionInternalExpression {
	t := p.pop()
	validateTokenValue(t, "=")
	if p.preValueCheck(0, "|") {
		_ = p.pop()
	}
	var ts []ast.UnionInternalExpression
	ts = append(ts, &ast.DefineUnionMemberExpression{p.parseTypeRef()})
	for p.preValueCheck(0, "|") {
		_ = p.pop()
		ts = append(ts, &ast.DefineUnionMemberExpression{p.parseTypeRef()})
	}
	return ts
}

func (p *Parser) parseInterface() ast.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "interface")
	n := p.parseName()
	directives := p.parseDirectivesOrEmpty()
	body := p.parseInterfaceBody()
	return &ast.DefineInterfaceExpression{
		DescriptionExpression: desc,
		NameExpression:        n,
		DirectiveExpressions:  directives,
		InterfaceExpression:   body,
	}
}

func (p *Parser) parseInterfaceBody() []ast.InterfaceInternalExpression {
	t := p.pop()
	validateTokenValue(t, "{")
	var exps []ast.InterfaceInternalExpression
	for !p.preValueCheck(0, "}") {
		exps = append(exps, p.parseInterfaceField())
	}
	_ = p.pop()
	return exps
}

func (p *Parser) parseType() ast.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "type")
	n := p.parseName()
	exps := p.parseImplementsOrEmpty()
	directives := p.parseDirectivesOrEmpty()
	exps = append(exps, p.parseObjectBody()...)

	return &ast.DefineObjectExpression{
		DescriptionExpression: desc,
		NameExpression:        n,
		DirectiveExpressions:  directives,
		ObjectExpression:      exps,
	}
}

func (p *Parser) parseObjectBody() []ast.ObjectInternalExpression {
	var exps []ast.ObjectInternalExpression
	t := p.pop()
	validateTokenValue(t, "{")
	for !p.preValueCheck(0, "}") {
		exps = append(exps, p.parseObjectField())
	}
	_ = p.pop()
	return exps
}

func (p *Parser) parseInterfaceField() ast.InterfaceInternalExpression {
	desc := p.parseDescriptionOrEmpty()
	n := p.parseName()
	args := p.parseFieldArgsOrEmpty()
	t := p.pop()
	validateTokenValue(t, ":")
	typ := p.parseTypeRef()
	directives := p.parseDirectivesOrEmpty()
	return &ast.DefineInterfaceFieldExpression{
		NameExp:        n,
		TypeExp:        typ,
		DescriptionExp: desc,
		DirectivesExp:  directives,
		ArgsExp:        args,
	}
}

func (p *Parser) parseObjectField() ast.ObjectInternalExpression {
	desc := p.parseDescriptionOrEmpty()
	n := p.parseName()
	args := p.parseFieldArgsOrEmpty()
	t := p.pop()
	validateTokenValue(t, ":")
	typ := p.parseTypeRef()
	directives := p.parseDirectivesOrEmpty()
	return &ast.DefineFieldExpression{
		Name:        n,
		TypeRef:     typ,
		Description: desc,
		Directives:  directives,
		Args:        args,
	}
}

func (p *Parser) parseFieldArgsOrEmpty() []ast.InputValueExpression {
	if !p.preValueCheck(0, "(") {
		return nil
	}
	var args []ast.InputValueExpression
	_ = p.pop()
	for !p.preValueCheck(0, ")") {
		args = append(args, p.parseInputValue())
	}
	_ = p.pop()
	return args
}

func (p *Parser) parseInputValue() ast.InputValueExpression {
	desc := p.parseDescriptionOrEmpty()
	n := p.parseName()
	t := p.pop()
	validateTokenValue(t, ":")
	typ := p.parseTypeRef()
	var def ast.ValueExpression
	if p.preValueCheck(0, "=") {
		_ = p.pop()
		def = p.parseValue()
	}
	directives := p.parseDirectivesOrEmpty()
	return &ast.InputValueExpressionImpl{
		Description:  desc,
		Name:         n,
		Type:         typ,
		DefaultValue: def,
		Directives:   directives,
	}
}

func (p *Parser) parseImplementsOrEmpty() []ast.ObjectInternalExpression {
	if !p.preValueCheck(0, "implements") {
		return nil
	}
	var exps []ast.ObjectInternalExpression
	_ = p.pop()
	if p.preValueCheck(0, "&") {
		_ = p.pop()
	}
	exps = append(exps, &ast.ImplementExpression{p.parseTypeRef()})
	for p.preValueCheck(0, "&") {
		_ = p.pop()
		exps = append(exps, &ast.ImplementExpression{p.parseTypeRef()})
	}
	return exps
}

func (p *Parser) parseEnum() ast.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "enum")
	name := p.parseName()
	directives := p.parseDirectivesOrEmpty()
	values := p.parseEnumBody()

	return &ast.DefineEnumExpression{
		DescriptionExpression: desc,
		NameExpression:        name,
		DirectiveExpressions:  directives,
		EnumExpression:        values,
	}
}

func (p *Parser) parseEnumBody() []ast.EnumInternalExpression {
	t := p.pop()
	validateTokenValue(t, "{")
	var values []ast.EnumInternalExpression
	for !p.preValueCheck(0, "}") {
		desc := p.parseDescriptionOrEmpty()
		name := p.parseName()
		directives := p.parseDirectivesOrEmpty()
		values = append(values, &ast.DefineEnumValueExpression{
			Directives:  directives,
			Name:        name,
			Description: desc,
		})
	}
	_ = p.pop()
	return values
}

func (p *Parser) parseScalar() ast.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "scalar")
	name := p.parseName()
	directives := p.parseDirectivesOrEmpty()
	return &ast.DefineScalarExpression{
		DescriptionExpression: desc,
		NameExpression:        name,
		DirectiveExpressions:  directives,
	}
}
func (p *Parser) parseDescriptionOrEmpty() ast.DescriptionExpression {
	t := p.lexer.Pop()
	if t.Type() != token.TypeStrVal {
		p.lexer.Push(t)
		return &ast.EmptyDescription{}
	}
	return &ast.DescriptionExpressionImpl{t.Value()}
}

func (p *Parser) parseSchema() ast.DefinitionExpression {
	t := p.pop()
	validateTokenValue(t, "schema")
	directives := p.parseDirectivesOrEmpty()
	exp := p.parseSchemaBody()
	return &ast.DefineSchemaExpression{
		Expressions: exp, DirectiveExpressions: directives,
	}
}

func (p *Parser) parseSchemaBody() []ast.SchemaInternalExpression {
	t := p.pop()
	validateTokenValue(t, "{")
	t = p.pop()
	var exp []ast.SchemaInternalExpression
	for t.Value() != "}" {
		switch t.Value() {
		case "query":
			t = p.pop()
			validateTokenValue(t, ":")
			texp := p.parseTypeRef()
			exp = append(exp, &ast.DefineQueryExpression{texp})
		case "mutation":
			t = p.pop()
			validateTokenValue(t, ":")
			texp := p.parseTypeRef()
			exp = append(exp, &ast.DefineMutationExpression{texp})
		case "subscription":
			t = p.pop()
			validateTokenValue(t, ":")
			texp := p.parseTypeRef()
			exp = append(exp, &ast.DefineSubscriptionExpression{texp})
		}
		t = p.pop()
	}
	return exp
}

func (p *Parser) parseDirectivesOrEmpty() []ast.DirectiveExpression {
	var directves []ast.DirectiveExpression
	for p.preValueCheck(0, "@") {
		_ = p.pop()
		name := p.parseName()
		args := map[string]ast.ValueExpression{}
		if p.preValueCheck(0, "(") {
			args = p.parseDirectiveArgs()
		}
		directves = append(
			directves,
			&ast.DirectiveExpressionImpl{name.Eval(), args},
		)
	}
	return directves
}

func (p *Parser) parseDirectives() []ast.DirectiveExpression {
	var directves []ast.DirectiveExpression
	validateTokenValue(p.prefetch(0), "@")
	for p.preValueCheck(0, "@") {
		_ = p.pop()
		name := p.parseName()
		args := p.parseDirectiveArgs()
		directves = append(
			directves,
			&ast.DirectiveExpressionImpl{name.Eval(), args},
		)
	}
	return directves
}

func (p *Parser) parseDirectiveArgs() map[string]ast.ValueExpression {
	args := map[string]ast.ValueExpression{}
	t := p.pop()
	validateTokenValue(t, "(")
	t = p.pop()
	for t.Value() != ")" {
		validateTokenType(t, token.TypeName)
		name := t.Value()
		t = p.pop()
		validateTokenValue(t, ":")
		v := p.parseValue()
		args[name] = v
		t = p.pop()
	}
	return args
}

func (p *Parser) parseValue() ast.ValueExpression {
	t := p.pop()
	switch t.Type() {
	case token.TypeStrVal, token.TypeFloatVal, token.TypeIntVal, token.TypeName:
		return &ast.ValueExpressionImpl{
			Value: t.Value(),
		}
	}
	if t.Value() == "[" {
		p.push(t)
		return p.parseListValue()
	}
	unexpectedToken(t)
	return nil
}

func (p *Parser) parseListValue() ast.ValueExpression {
	t := p.pop()
	validateTokenValue(t, "[")
	var children []ast.ValueExpression
	for !p.preValueCheck(0, "]") {
		children = append(children, p.parseValue())
	}
	p.pop()
	return &ast.ListValueExpressionImpl{Children: children}
}

func (p *Parser) parseName() ast.NameExpression {
	t := p.lexer.Pop()
	validateTokenType(t, token.TypeName)
	return &ast.NameExpressionImpl{t.Value()}
}

func (p *Parser) parseTypeRef() ast.TypeRefExpression {
	t := p.prefetch(0)
	validateToken(t, func(t *token.Token) bool {
		return t.Value() == "[" || t.Type() == token.TypeName
	})
	if t.Value() == "[" {
		return p.parseList()
	}
	name := p.parseName()
	isNullable := true
	if p.preValueCheck(0, "!") {
		_ = p.pop()
		isNullable = false
	}
	return &ast.TypeRefExpressionImpl{nil, isNullable, name}
}

func (p *Parser) parseList() ast.TypeRefExpression {
	t := p.pop()
	validateTokenValue(t, "[")
	inner := p.parseTypeRef()
	t = p.pop()
	validateTokenValue(t, "]")
	isNullable := true
	if p.preValueCheck(0, "!") {
		_ = p.pop()
		isNullable = false
	}
	return &ast.TypeRefExpressionImpl{inner, isNullable, &ast.NameExpressionImpl{"[]"}}
}

func validateTokenValue(t *token.Token, value ...string) {
	for _, v := range value {
		if t.Value() == v {
			return
		}
	}
	unexpectedToken(t)
}

func validateTokenType(token *token.Token, ts ...token.Type) {
	for _, t := range ts {
		if token.Type() == t {
			return
		}
	}
	unexpectedToken(token)
}

func assertNotNil(t *token.Token) {
	if t == nil {
		log.Fatalf("unexpected eof")
	}
}

func validateToken(token *token.Token, predicate func(t *token.Token) bool) {
	if !predicate(token) {
		unexpectedToken(token)
	}
}

func unexpectedToken(t *token.Token) {
	l, c := t.LineCol()
	panic(errors.New(""))
	log.Fatalf("unexpected token '%s' at line %d, col %d", t.Value(), l, c)
}

func (p *Parser) pop() *token.Token {
	t := p.lexer.Pop()
	assertNotNil(t)
	return t
}

func (p *Parser) popOrNil() *token.Token {
	t := p.lexer.Pop()
	return t
}

func (p *Parser) push(t *token.Token) {
	p.lexer.Push(t)
}

func (p *Parser) preCheck(index int, predicate func(t *token.Token) bool) bool {
	var store []*token.Token
	defer func() {
		for _, t := range store {
			p.push(t)
		}
	}()
	for i := 0; i <= index; i++ {
		t := p.popOrNil()
		if t == nil {
			return false
		}
		store = append([]*token.Token{t}, store...)
	}
	return predicate(store[0])
}

func (p *Parser) preTypeCheck(index int, t token.Type) bool {
	return p.preCheck(index, func(tok *token.Token) bool {
		return tok.Type() == t
	})
}

func (p *Parser) preValueCheck(index int, v string) bool {
	return p.preCheck(index, func(t *token.Token) bool {
		return t.Value() == v
	})
}

func (p *Parser) prefetch(index int) *token.Token {
	var store []*token.Token
	for i := 0; i <= index; i++ {
		store = append([]*token.Token{p.pop()}, store...)
	}
	result := store[0]
	for _, t := range store {
		p.push(t)
	}
	return result
}

func (p *Parser) hasNext() bool {
	t := p.popOrNil()
	if t == nil {
		return false
	}
	p.push(t)
	return true
}
