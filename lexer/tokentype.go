package lexer

type TokenType int

const (
	Schema   TokenType = iota // schema
	Scalar                    // scalar
	Type                      // type
	LParen                    // (
	RParen                    // )
	LBrace                    // {
	RBrace                    // }
	LBracket                  // [
	RBracket                  // ]
	NewLine                   // \n
	Colon                     // :
	Comma                     // ,
	NotNull                   // !
	Eq                        // =
	Number                    // [0-9]+
	String                    // ".*"
	Comment                   // #.*
	Id
)

//go:generate stringer -type=TokenType
