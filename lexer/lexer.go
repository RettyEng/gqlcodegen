package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
)

type Lexer_ struct {
	buff *bufio.Reader
}

func NewLexer_(r *bufio.Reader) *Lexer_ {
	return &Lexer_{buff: r}
}

func (t *Lexer_) takeRunes() []rune {
	runes := make([]rune, 0)
	for {
		var r rune
		var e error

		if r, _, e = t.buff.ReadRune(); e != nil {
			break
		}
		runes = append(runes, r)
	}
	return runes
}

func (t *Lexer_) Tokenize() []*Token {
	var tokens []*Token
	runes := t.takeRunes()
	cur := 0
	line := 1
	column := 1

	for cur < len(runes) {
		token, consumed := takeTokenValue(runes[cur:])
		tokens = append(
			tokens, convertValueToToken(string(token), line, column),
		)
		column += consumed
		if r := runes[cur]; r == '\n' {
			line++
			column = 1
		}
		cur += consumed
	}
	return tokens
}

func convertValueToToken(value string, line, column int) *Token {
	switch value {
	case "schema":
		return schemaT(line, column)
	case "scalar":
		return scalarT(line, column)
	case "type":
		return typeT(line, column)
	case "enum":
		return enumT(line, column)
	case "true", "false":
		return boolT(line, column, value)
	case "null":
		return nullT(line, column)
	case "(":
		return lParenT(line, column)
	case ")":
		return rParenT(line, column)
	case "{":
		return lBraceT(line, column)
	case "}":
		return rBraceT(line, column)
	case "[":
		return lBracketT(line, column)
	case "]":
		return rBracketT(line, column)
	case "\n":
		return newLineT(line, column)
	case ":":
		return colonT(line, column)
	case ",":
		return commaT(line, column)
	case "!":
		return notNullT(line, column)
	case "=":
		return eqT(line, column)
	}
	r := regexp.MustCompile(`^\d+(\.\d+)?$`)
	if r.MatchString(value) {
		return numberT(line, column, value)
	}
	r = regexp.MustCompile(`^".*"$`)
	if r.MatchString(value) {
		return stringT(line, column, value)
	}
	r = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	if r.MatchString(value) {
		return idT(line, column, value)
	}
	r = regexp.MustCompile(`^#.+$`)
	if r.MatchString(value) {
		return commentT(line, column, value)
	}
	panic(errors.New(fmt.Sprintf(`failed to parse "%s"`, value)))
}

func takeTokenValue(rs []rune) (string, int) {
	consumed := 0
	for cur := 0; rs[cur] == ' ' || rs[cur] == '\t'; cur++ {
		consumed++
	}
	var token []rune
	switch r := rs[consumed]; r {
	case '(', ')', '{', '}', '[', ']', '\n', ':', '!', '=', ',':
		return string([]rune{r}), consumed + 1
	case '"':
		token = takeString(rs[consumed:])
	case '#':
		token = takeComment(rs[consumed:])
	default:
		token = takeWhileSeparator(rs[consumed:])
	}
	return string(token), consumed + len(token)
}

func takeString(rs []rune) []rune {
	str := make([]rune, 0)
	cur := 0
	isEscaping := true
	r := rs[cur]
	for {
		str = append(str, r)
		if r == '"' && !isEscaping {
			break
		}
		if !isEscaping && r == '\\' {
			isEscaping = true
		} else {
			isEscaping = false
		}
		cur++
		r = rs[cur]
	}
	return str
}

func takeComment(rs []rune) []rune {
	comment := make([]rune, 0)
	for cur := 0; rs[cur] != '\n'; cur++ {
		comment = append(comment, rs[cur])
	}
	return comment
}

func takeWhileSeparator(rs []rune) []rune {
	token := make([]rune, 0)
	separators := map[rune]struct{}{
		' ':  {},
		'\t': {},
		'(':  {},
		')':  {},
		'{':  {},
		'}':  {},
		'[':  {},
		']':  {},
		'\n': {},
		':':  {},
		',':  {},
		'!':  {},
		'=':  {},
	}

	cur := 0
	for {
		r := rs[cur]
		if _, contains := separators[r]; contains {
			break
		}
		token = append(token, r)
		cur++
	}
	return token
}
