package parser

import (
	"bufio"
	"errors"
	"fmt"

	"github.com/RettyInc/gqlcodegen/ast"
	"github.com/RettyInc/gqlcodegen/lexer"
)

type Parser struct {
	lexer *lexer.Lexer
}

func NewParser(reader *bufio.Reader) *Parser {
	return &Parser{lexer.NewLexer(reader)}
}

func (p *Parser) tokens() []*lexer.Token {
	var tokens []*lexer.Token
	for _, t := range p.lexer.Tokenize() {
		if t.Type != lexer.Comment && t.Type != lexer.NewLine {
			tokens = append(tokens, t)
		}
	}
	return tokens
}

func validateTokenTypeOrErr(t *lexer.Token, expected ...lexer.TokenType) {
	for _, tt := range expected {
		if t.Type == tt {
			return
		}
	}
	panic(errors.New(
		fmt.Sprintf(
			"unexpected token %s(%s) at line %d, column %d. expected %v.",
			t.Type.String(), t.Value, t.Line, t.Column, expected,
		),
	))
}

func validateTokenValueOrErr(t *lexer.Token, expected ...string) {
	for _, tv := range expected {
		if t.Value == tv {
			return
		}
	}
	panic(errors.New(
		fmt.Sprintf(
			"unexpected token %s(%s) at line %d, column %d. expected %v.",
			t.Type.String(), t.Value, t.Line, t.Column, expected,
		),
	))
}

func (p *Parser) Parse() *ast.Root {
	tokens := p.tokens()
	var children []ast.Ast
	cur := 0
	for cur < len(tokens) {
		t := tokens[cur]
		var ast ast.Ast
		consumes := 0
		switch t.Type {
		case lexer.Schema:
			ast, consumes = parseSchema(tokens[cur:])
		case lexer.Scalar:
			ast, consumes = parseScalar(tokens[cur:])
		case lexer.Enum:
			ast, consumes = parseEnum(tokens[cur:])
		case lexer.Type:
			ast, consumes = parseType(tokens[cur:])
		default:
			validateTokenTypeOrErr(
				t, lexer.Schema, lexer.Scalar, lexer.Type, lexer.Enum,
			)
		}
		cur += consumes
		children = append(children, ast)
	}
	return ast.NewRoot(children)
}

func parseType(ts []*lexer.Token) (*ast.TypeDef, int) {
	cur := 0
	validateTokenTypeOrErr(ts[cur], lexer.Type)
	cur++
	validateTokenTypeOrErr(ts[cur], lexer.Id)

	name := ts[cur].Value

	cur++
	validateTokenTypeOrErr(ts[cur], lexer.LBrace)
	cur++

	fields, consumes := parseTypeBody(ts[cur:])

	cur += consumes
	validateTokenTypeOrErr(ts[cur], lexer.RBrace)
	return ast.NewTypeDef(name, fields), cur + 1
}

func parseTypeBody(ts []*lexer.Token) ([]ast.Ast, int) {
	var fields []ast.Ast
	cur := 0
	for {
		if ts[cur].Type == lexer.RBrace {
			break
		}
		name := ts[cur].Value
		var args []*ast.FieldArgDef
		cur++
		t := ts[cur]
		validateTokenTypeOrErr(t, lexer.LParen, lexer.Colon)
		if t.Type == lexer.LParen {
			cur++
			a, c := parseFieldArgs(ts[cur:])
			args = a
			cur += c
			validateTokenTypeOrErr(ts[cur], lexer.RParen)
			cur++
		}
		validateTokenTypeOrErr(ts[cur], lexer.Colon)
		cur++
		ftype, consumed := parseTypeRef(ts[cur:])
		fields = append(fields, ast.NewTypeFieldDef(name, ftype, args...))
		cur += consumed
	}
	return fields, cur
}

func parseFieldArgs(ts []*lexer.Token) ([]*ast.FieldArgDef, int) {
	var args []*ast.FieldArgDef
	cur := 0
	for {
		if ts[cur].Type == lexer.RParen {
			break
		}
		name := ts[cur].Value
		cur++
		validateTokenTypeOrErr(ts[cur], lexer.Colon)
		cur++
		t, c := parseTypeRef(ts[cur:])
		cur += c
		validateTokenTypeOrErr(ts[cur], lexer.Comma, lexer.RParen, lexer.Eq)
		var defaultValue *ast.FieldArgDefault
		if ts[cur].Type == lexer.Eq {
			cur++
			validateTokenTypeOrErr(
				ts[cur], lexer.Number, lexer.String, lexer.Bool,
				lexer.LBracket, lexer.Id, lexer.Null,
			)
			if ts[cur].Type == lexer.LBracket {
				cur++
				values := ""
				for ts[cur].Type != lexer.RBracket {
					if ts[cur].Type == lexer.Comma {
						values += ","
					} else {
						validateTokenTypeOrErr(
							ts[cur], lexer.Number, lexer.String, lexer.Bool,
							lexer.LBracket, lexer.Id, lexer.Null,
						)
						values += ts[cur].Value
					}
					cur++
				}
				validateTokenTypeOrErr(ts[cur], lexer.RBracket)
				cur++
				defaultValue = ast.NewFieldArgDefault("[" + values + "]")
			} else {
				defaultValue = ast.NewFieldArgDefault(ts[cur].Value)
				cur++
			}
		}
		validateTokenTypeOrErr(ts[cur], lexer.RParen, lexer.Comma)
		if ts[cur].Type == lexer.RParen {
			args = append(args, ast.NewFieldArgDef(name, t, defaultValue))
			continue
		}
		args = append(args, ast.NewFieldArgDef(name, t, defaultValue))
		cur++
	}
	return args, cur
}

func parseSchema(ts []*lexer.Token) (*ast.SchemaDef, int) {
	cur := 0
	validateTokenTypeOrErr(ts[cur], lexer.Schema)
	cur++
	validateTokenTypeOrErr(ts[cur], lexer.LBrace)
	cur++

	query, consumed := parseSchemaQuery(ts[cur:])

	cur += consumed
	validateTokenTypeOrErr(ts[cur], lexer.RBrace)
	cur++
	return ast.NewSchemaDef(query), cur
}

func parseSchemaQuery(ts []*lexer.Token) (*ast.SchemaQueryDef, int) {
	cur := 0
	validateTokenTypeOrErr(ts[cur], lexer.Id)
	validateTokenValueOrErr(ts[cur], "query")
	cur++
	validateTokenTypeOrErr(ts[cur], lexer.Colon)
	cur++

	typeRef, consumed := parseTypeRef(ts[cur:])

	cur += consumed
	return ast.NewSchemaQueryDef(typeRef), cur
}

func parseTypeRef(ts []*lexer.Token) (*ast.TypeRef, int) {
	cur := 0
	t := ts[cur]
	validateTokenTypeOrErr(t, lexer.Id, lexer.LBracket)
	inArray := t.Type == lexer.LBracket
	var typeVars []ast.Ast
	if !inArray {
		name := t.Value
		cur++
		isNotNull := ts[cur].Type == lexer.NotNull
		if isNotNull {
			cur++
		}
		return ast.NewTypeRef(name, typeVars, !isNotNull), cur
	}
	cur++
	typeRef, consumed := parseTypeRef(ts[cur:])
	cur += consumed
	typeVars = append(typeVars, typeRef)
	validateTokenTypeOrErr(ts[cur], lexer.RBracket)
	cur++
	isNotNull := ts[cur].Type == lexer.NotNull
	if isNotNull {
		cur++
	}
	return ast.NewTypeRef("[]", typeVars, !isNotNull), cur
}

func parseEnum(ts []*lexer.Token) (*ast.EnumDef, int) {
	cur := 0
	validateTokenTypeOrErr(ts[cur], lexer.Enum)
	cur++
	validateTokenTypeOrErr(ts[cur], lexer.Id)

	name := ts[cur].Value

	cur++
	validateTokenTypeOrErr(ts[cur], lexer.LBrace)
	cur++
	children, consumes := parseEnumBody(ts[cur:])

	cur += consumes
	validateTokenTypeOrErr(ts[cur], lexer.RBrace)

	return ast.NewEnumDef(name, children), cur + 1
}

func parseEnumBody(ts []*lexer.Token) ([]ast.Ast, int) {
	cur := 0
	var children []ast.Ast
	for {
		t := ts[cur]
		validateTokenTypeOrErr(t, lexer.Id, lexer.RBrace)
		if t.Type == lexer.RBrace {
			return children, cur
		}
		children = append(children, ast.NewEnumEntryDef(t.Value))
		cur++
	}
}

func parseScalar(ts []*lexer.Token) (*ast.ScalarDef, int) {
	cur := 0
	validateTokenTypeOrErr(ts[cur], lexer.Scalar)
	cur++
	validateTokenTypeOrErr(ts[cur], lexer.Id)

	return ast.NewScalarDef(ts[cur].Value), cur + 1
}
