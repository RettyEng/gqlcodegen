package lexer

const (
	linInit = 1
	colInit = 1
)

type Scanner struct {
	line int
	col int
	runes []rune
}

func NewScanner(rs []rune) *Scanner {
	return &Scanner{
		line:  linInit,
		col:   colInit,
		runes: rs,
	}
}

func (s *Scanner) startsWith(str string) bool {
	for i, r := range []rune(str) {
		if s.runes[i] != r {
			return false
		}
	}
	return true
}

func (s *Scanner) startsWithCharset(cs *Charset) bool {
	return cs.Contains(s.runes[0])
}

func (s *Scanner) HasNext() bool {
	return len(s.runes) != 0
}

func (s *Scanner) updateLineCol(popped rune) {
	switch popped {
	case '\n':
		s.line++
		s.col = colInit
		return
	case '\r':
		if !s.startsWith("\n") {
			s.line++
			s.col = colInit
			return
		}
	}
	s.col++
}

func (s *Scanner) LineCol() (int, int) {
	return s.line, s.col
}
func (s *Scanner) Pop() (rune, int, int) {
	line, col := s.LineCol()

	r := s.runes[0]
	s.runes = s.runes[1:]
	s.updateLineCol(r)
	return r, line, col
}

func (s *Scanner) PopN(n int) ([]rune, int, int) {
	l, c := s.LineCol()
	rs := make([]rune, n)
	for i:=0; i<n; i++ {
		r, _, _ := s.Pop()
		rs[i] = r
	}
	return rs, l, c
}
