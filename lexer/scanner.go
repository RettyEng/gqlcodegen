package lexer

const (
	linInit = 1
	colInit = 1
)

type Scanner struct {
	line  int
	col   int
	runes []rune
}

func NewScanner(rs []rune) *Scanner {
	return &Scanner{
		line:  linInit,
		col:   colInit,
		runes: rs,
	}
}

func (s *Scanner) Take(matcher Matcher) (string, int, int) {
	c := matcher.MatchCount(s.runes)
	r, line, col := s.PopN(c)
	return string(r), line, col
}

func (s *Scanner) TakeWhileMatch(matcher Matcher) (string, int, int) {
	var ret []rune
	line, col := s.LineCol()
	c := matcher.MatchCount(s.runes)
	for c != 0 {
		r, _, _ := s.PopN(c)
		ret = append(ret, r...)
		c = matcher.MatchCount(s.runes)
	}
	return string(ret), line, col
}

func (s *Scanner) StartsWith(m Matcher) bool {
	return m.HeadMatches(s.runes)
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
		if !s.StartsWith(Str("\n")) {
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
	for i := 0; i < n; i++ {
		r, _, _ := s.Pop()
		rs[i] = r
	}
	return rs, l, c
}
