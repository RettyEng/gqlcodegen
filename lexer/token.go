package lexer

type Token struct {
	Type   TokenType
	Value  string
	Line   int
	Column int
}

func schemaT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: Schema}
}

func scalarT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: Scalar}
}

func typeT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: Type}
}

func enumT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: Enum}
}

func lParenT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: LParen}
}

func rParenT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: RParen}
}

func lBraceT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: LBrace}
}

func rBraceT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: RBrace}
}

func lBracketT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: LBracket}
}

func rBracketT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: RBracket}

}

func newLineT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: NewLine}
}

func colonT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: Colon}
}

func commaT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: Comma}
}

func notNullT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: NotNull}
}

func eqT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: Eq}
}

func numberT(line, column int, value string) *Token {
	return &Token{Line: line, Column: column, Type: Number, Value: value}
}

func stringT(line, column int, value string) *Token {
	return &Token{Line: line, Column: column, Type: String, Value: value}
}

func boolT(line, column int, value string) *Token {
	return &Token{Line: line, Column: column, Type: Bool, Value: value}
}

func nullT(line, column int) *Token {
	return &Token{Line: line, Column: column, Type: Null}
}

func commentT(line, column int, value string) *Token {
	return &Token{Line: line, Column: column, Type: Comment, Value: value}
}

func idT(line, column int, value string) *Token {
	return &Token{Line: line, Column: column, Type: Id, Value: value}
}
