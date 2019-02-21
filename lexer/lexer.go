package lexer

import (
	"bufio"
	"errors"
	"fmt"
	"regexp"
)

type Lexer struct {
	buff *bufio.Reader
}

func NewLexer(r *bufio.Reader) *Lexer {
	return &Lexer{buff: r}
}

func (t *Lexer) takeRunes() []rune {
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

func (t *Lexer) Tokenize() []*Token {
	var tokens []*Token
	runes := t.takeRunes()
	cur := 0

	for cur < len(runes) {
		token, consumed := takeTokenValue(runes[cur:])
		cur += consumed
		tokens = append(tokens, convertValueToToken(string(token)))
	}
	return tokens
}

func convertValueToToken(value string) *Token {
	switch value {
	case "schema":
		return schemaT()
	case "scalar":
		return scalarT()
	case "type":
		return typeT()
	case "(":
		return lParenT()
	case ")":
		return rParenT()
	case "{":
		return lBraceT()
	case "}":
		return rBraceT()
	case "[":
		return lBracketT()
	case "]":
		return rBracketT()
	case "\n":
		return newLineT()
	case ":":
		return colonT()
	case ",":
		return commaT()
	case "!":
		return notNullT()
	case "=":
		return eqT()
	}
	r := regexp.MustCompile(`^\d+$`)
	if r.MatchString(value) {
		return numberT(value)
	}
	r = regexp.MustCompile(`^".*"$`)
	if r.MatchString(value) {
		return stringT(value)
	}
	r = regexp.MustCompile(`^[a-zA-Z0-9_\-]+$`)
	if r.MatchString(value) {
		return idT(value)
	}
	r = regexp.MustCompile(`^#.+$`)
	if r.MatchString(value) {
		return commentT(value)
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
