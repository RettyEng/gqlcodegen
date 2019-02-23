package parser2

import (
	"errors"
	"io"
	"log"
	"strings"

	"github.com/RettyInc/gqlcodegen/ast2/directive"

	"github.com/RettyInc/gqlcodegen/ast2"
	"github.com/RettyInc/gqlcodegen/gql"
	"github.com/RettyInc/gqlcodegen/lexer"
	"github.com/RettyInc/gqlcodegen/lexer/token"
)

type Parser struct {
	lexer *lexer.Lexer
}

func NewParser(reader io.Reader) *Parser {
	return &Parser{
		lexer.NewLexer(reader),
	}
}

func (p *Parser) Parse() *gql.TypeSystem {
	var exp []ast2.DefinitionExpression
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
		}

	}
	return (&ast2.TopLevel{exp}).Eval()
}

func (p *Parser) parseInput() ast2.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "input")
	n := p.parseName()
	direc := p.parseDirectivesOrEmpty()
	t = p.pop()
	validateTokenValue(t, "{")
	var args []ast2.InputValueExpression
	for !p.preValueCheck(0, "}") {
		args = append(args, p.parseInputValue())
	}
	_ = p.pop()
	return &ast2.DefineInputObjectExpression{
		DescriptionExpression:             desc,
		NameExpression:                    n,
		DirectiveExpressions:              direc,
		DefineInputObjectFieldExpressions: args,
	}
}

func (p *Parser) parseDirective() ast2.DefinitionExpression {
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
	var locs []ast2.DirectiveInternalExpression
	locs = append(locs, p.parseDirectiveLocation())
	for p.preValueCheck(0, "|") {
		_ = p.pop()
		locs = append(locs, p.parseDirectiveLocation())
	}
	return &ast2.DirectiveDefinition{
		DescriptionExpression: desc,
		NameExpression:        n,
		ArgsExpression:        args,
		Expressions:           locs,
	}
}

func (p *Parser) parseDirectiveLocation() ast2.DirectiveInternalExpression {
	t := p.pop()
	validateTokenType(t, token.TypeName)
	for i := 0; !strings.HasPrefix(directive.Location(i).String(), "Location("); i++ {
		if t.Value() == directive.Location(i).String() {
			return &ast2.DefineDirectiveLocationExpression{directive.Location(i)}
		}
	}
	l, c := t.LineCol()
	log.Fatalf("unknown directive location %s at line %d, col %d", t.Value(), l, c)
	return nil
}

func (p *Parser) parseDirectiveArgsDefinition() []ast2.InputValueExpression {
	t := p.pop()
	validateTokenValue(t, "(")
	var args []ast2.InputValueExpression
	for !p.preValueCheck(0, ")") {
		args = append(args, p.parseInputValue())
	}
	_ = p.pop()
	return args
}

func (p *Parser) parseUnion() ast2.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "union")
	n := p.parseName()
	directives := p.parseDirectivesOrEmpty()
	t = p.pop()
	validateTokenValue(t, "=")
	if p.preValueCheck(0, "|") {
		_ = p.pop()
	}
	var ts []ast2.UnionInternalExpression
	ts = append(ts, &ast2.DefineUnionMemberExpression{p.parseTypeRef()})
	for p.preValueCheck(0, "|") {
		_ = p.pop()
		ts = append(ts, &ast2.DefineUnionMemberExpression{p.parseTypeRef()})
	}
	return &ast2.DefineUnionExpression{
		DescriptionExpression: desc,
		NameExpression:        n,
		DirectiveExpressions:  directives,
		UnionExpression:       ts,
	}
}

func (p *Parser) parseInterface() ast2.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "interface")
	n := p.parseName()
	directives := p.parseDirectivesOrEmpty()
	t = p.pop()
	validateTokenValue(t, "{")
	var exps []ast2.InterfaceInternalExpression
	for !p.preValueCheck(0, "}") {
		exps = append(exps, p.parseInterfaceField())
	}
	_ = p.pop()
	return &ast2.DefineInterfaceExpression{
		DescriptionExpression: desc,
		NameExpression:        n,
		DirectiveExpressions:  directives,
		InterfaceExpression:   exps,
	}
}

func (p *Parser) parseType() ast2.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "type")
	n := p.parseName()
	exps := p.parseImplementsOrEmpty()
	directives := p.parseDirectivesOrEmpty()
	t = p.pop()
	validateTokenValue(t, "{")
	for !p.preValueCheck(0, "}") {
		exps = append(exps, p.parseObjectField())
	}
	_ = p.pop()
	return &ast2.DefineObjectExpression{
		DescriptionExpression: desc,
		NameExpression:        n,
		DirectiveExpressions:  directives,
		ObjectExpression:      exps,
	}
}

func (p *Parser) parseInterfaceField() ast2.InterfaceInternalExpression {
	desc := p.parseDescriptionOrEmpty()
	n := p.parseName()
	args := p.parseFieldArgsOrEmpty()
	t := p.pop()
	validateTokenValue(t, ":")
	typ := p.parseTypeRef()
	directives := p.parseDirectivesOrEmpty()
	return &ast2.DefineInterfaceFieldExpression{
		NameExp:        n,
		TypeExp:        typ,
		DescriptionExp: desc,
		DirectivesExp:  directives,
		ArgsExp:        args,
	}
}

func (p *Parser) parseObjectField() ast2.ObjectInternalExpression {
	desc := p.parseDescriptionOrEmpty()
	n := p.parseName()
	args := p.parseFieldArgsOrEmpty()
	t := p.pop()
	validateTokenValue(t, ":")
	typ := p.parseTypeRef()
	directives := p.parseDirectivesOrEmpty()
	return &ast2.DefineFieldExpression{
		Name:        n,
		TypeRef:     typ,
		Description: desc,
		Directives:  directives,
		Args:        args,
	}
}

func (p *Parser) parseFieldArgsOrEmpty() []ast2.InputValueExpression {
	if !p.preValueCheck(0, "(") {
		return nil
	}
	var args []ast2.InputValueExpression
	_ = p.pop()
	for !p.preValueCheck(0, ")") {
		args = append(args, p.parseInputValue())
	}
	_ = p.pop()
	return args
}

func (p *Parser) parseInputValue() ast2.InputValueExpression {
	desc := p.parseDescriptionOrEmpty()
	n := p.parseName()
	t := p.pop()
	validateTokenValue(t, ":")
	typ := p.parseTypeRef()
	var def ast2.ValueExpression
	if p.preValueCheck(0, "=") {
		_ = p.pop()
		def = p.parseValue()
	}
	directives := p.parseDirectivesOrEmpty()
	return &ast2.InputValueExpressionImpl{
		Description:  desc,
		Name:         n,
		Type:         typ,
		DefaultValue: def,
		Directives:   directives,
	}
}

func (p *Parser) parseImplementsOrEmpty() []ast2.ObjectInternalExpression {
	if !p.preValueCheck(0, "implements") {
		return nil
	}
	var exps []ast2.ObjectInternalExpression
	_ = p.pop()
	if p.preValueCheck(0, "&") {
		_ = p.pop()
	}
	exps = append(exps, &ast2.ImplementExpression{p.parseTypeRef()})
	for p.preValueCheck(0, "&") {
		_ = p.pop()
		exps = append(exps, &ast2.ImplementExpression{p.parseTypeRef()})
	}
	return exps
}

func (p *Parser) parseEnum() ast2.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "enum")
	name := p.parseName()
	directives := p.parseDirectivesOrEmpty()
	t = p.pop()
	validateTokenValue(t, "{")
	var values []ast2.EnumInternalExpression
	for !p.preValueCheck(0, "}") {
		desc := p.parseDescriptionOrEmpty()
		name := p.parseName()
		directives := p.parseDirectivesOrEmpty()
		values = append(values, &ast2.DefineEnumValueExpression{
			Directives:  directives,
			Name:        name,
			Description: desc,
		})
	}
	_ = p.pop()
	return &ast2.DefineEnumExpression{
		DescriptionExpression: desc,
		NameExpression:        name,
		DirectiveExpressions:  directives,
		EnumExpression:        values,
	}
}
func (p *Parser) parseScalar() ast2.DefinitionExpression {
	desc := p.parseDescriptionOrEmpty()
	t := p.pop()
	validateTokenValue(t, "scalar")
	name := p.parseName()
	directives := p.parseDirectivesOrEmpty()
	return &ast2.DefineScalarExpression{
		DescriptionExpression: desc,
		NameExpression:        name,
		DirectiveExpressions:  directives,
	}
}
func (p *Parser) parseDescriptionOrEmpty() ast2.DescriptionExpression {
	t := p.lexer.Pop()
	if t.Type() != token.TypeStrVal {
		p.lexer.Push(t)
		return &ast2.EmptyDescription{}
	}
	return &ast2.DescriptionExpressionImpl{t.Value()}
}

func (p *Parser) parseSchema() ast2.DefinitionExpression {
	t := p.pop()
	validateTokenValue(t, "schema")
	directives := p.parseDirectivesOrEmpty()
	t = p.pop()
	validateTokenValue(t, "{")
	t = p.pop()
	var exp []ast2.SchemaInternalExpression
	for t.Value() != "}" {
		switch t.Value() {
		case "query":
			t = p.pop()
			validateTokenValue(t, ":")
			texp := p.parseTypeRef()
			exp = append(exp, &ast2.DefineQueryExpression{texp})
		case "mutation":
			t = p.pop()
			validateTokenValue(t, ":")
			texp := p.parseTypeRef()
			exp = append(exp, &ast2.DefineMutationExpression{texp})
		case "subscription":
			t = p.pop()
			validateTokenValue(t, ":")
			texp := p.parseTypeRef()
			exp = append(exp, &ast2.DefineSubscriptionExpression{texp})
		}
		t = p.pop()
	}
	if t.Value() != "}" {
		unexpectedToken(t)
	}
	return &ast2.DefineSchemaExpression{
		Expressions: exp, DirectiveExpressions: directives,
	}
}

func (p *Parser) parseDirectivesOrEmpty() []ast2.DirectiveExpression {
	var directves []ast2.DirectiveExpression
	for p.preValueCheck(0, "@") {
		_ = p.pop()
		name := p.parseName()
		args := p.parseDirectiveArgs()
		directves = append(
			directves,
			&ast2.DirectiveExpressionImpl{name.Eval(), args},
		)
	}
	return directves
}

func (p *Parser) parseDirectiveArgs() map[string]ast2.ValueExpression {
	args := map[string]ast2.ValueExpression{}
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

func (p *Parser) parseValue() ast2.ValueExpression {
	t := p.pop()
	switch t.Type() {
	case token.TypeStrVal, token.TypeFloatVal, token.TypeIntVal, token.TypeName:
		return &ast2.ValueExpressionImpl{
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

func (p *Parser) parseListValue() ast2.ValueExpression {
	t := p.pop()
	validateTokenValue(t, "[")
	var children []ast2.ValueExpression
	for !p.preValueCheck(0, "]") {
		children = append(children, p.parseValue())
	}
	p.pop()
	return &ast2.ListValueExpressionImpl{Children: children}
}

func (p *Parser) parseName() ast2.NameExpression {
	t := p.lexer.Pop()
	validateTokenType(t, token.TypeName)
	return &ast2.NameExpressionImpl{t.Value()}
}

func (p *Parser) parseTypeRef() ast2.TypeRefExpression {
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
	return &ast2.TypeRefExpressionImpl{nil, isNullable, name}
}

func (p *Parser) parseList() ast2.TypeRefExpression {
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
	return &ast2.TypeRefExpressionImpl{inner, isNullable, &ast2.NameExpressionImpl{"[]"}}
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
	log.Printf("unexpected token '%s' at line %d, col %d", t.Value(), l, c)
	panic(errors.New(""))

}

func (p *Parser) pop() *token.Token {
	t := p.lexer.Pop()
	log.Printf("%+v", t)
	assertNotNil(t)
	return t
}

func (p *Parser) popOrNil() *token.Token {
	t := p.lexer.Pop()
	log.Printf("%+v", t)
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
