package gomore

import (
	"bufio"
	"fmt"
	"unicode"
)

const (
	T_UNDECIDED = iota
	T_GUARD
	T_INT
	T_FLOAT
	T_HEX
	T_PERCENT
	T_PIXELS
	T_STRING
	T_COMMENT
	T_LINECOMMENT
	T_PATH
	T_LPAREN
	T_RPAREN
	T_LBRACE
	T_RBRACE
	T_LBRACK
	T_RBRACK
	T_SELECTOR
	T_VARIABLE
	T_COLON
	T_SEMICOLON
	T_DASH
	T_MUL
	T_PLUS
	T_MINUS
	T_COMMA
	T_PERIOD
	T_GT
	T_LT
	T_EQ
	T_DIV
)

type Validator func(s *Scanner, t *Token) bool
type LexicalRule struct {
	t            int // what we will set of_type to
	chars_in     []rune // list for inclusion
	chars_not_in []rune // list for exclusion
	starts_with  []rune // must begin with
	ends_with    []rune // must end with
	must_be_true Validator // A general purpose validator
	reclassifier Validator // A func to possibly reclassify
}

type SyntacticRule struct {
	concerns int // the token type that this rule will apply too
	/*
	   before and after only help us when context can be discerned by
	   an immediate neighbor.
	*/
	before []int // token types that MUST preceed IMMEDIATELY
	after  []int // token types that MUST Follow IMMEDIATELY

	/*
	   What we want to solve is: border: 1px solid #fff; and make sure
	   that #fff is defined as a hex instead of a selector.

	   So context_before is an array of token types that may preceed
	   the current token and therefore influence it's context
	   _reject means if any of these are found BEFORE context_before is
	   satisfied, then this rule cannot apply.

	   in the example above: context_before = []int{T_COLON}
	   and context_before_reject = []int{T_SEMICOLON} tells us
	   to go backwards until we find a semi-colon (newline) if
	   we find a semi-colon before a colon, then a rule that would
	   convert a selector to a hex has failed.

	   context_after and context_after_reject are for lookahead.

	   if both context_before and context_after are non-zero length,
	   then both must be satisfied.
	*/
	context_before        []int
	context_before_reject []int
	context_after         []int
	context_after_reject  []int
	consume_until         []int // Consume all tokens until this one is encountered,
	// ... normally this means consume until T_SEMICOLON or T_BRACETYPE
	// So when we encounter {, (, : etc we are looking at a Block, Tuple, or
	// attribute assignment. : is a special case because : can also be a
	// kind of method, like a:visited

}
type Token struct {
	literal                []rune // everything is kept as a rune array
	of_type                int // the T_XX type
	col                    int // the column, for error reporting
	row                    int // the row, for error reporting
	needs_interpolation    bool // if a string is detected to contain variables
	parent, next, previous *Token // overkill, but heh
	children               []*Token // possibly not used
}
type GomoreError struct {
    msg string
    code int
    row,col int
}
type Scanner struct {
	ch           rune
    nch          rune
    lch          rune
	col          int
	row          int
	tokens       []*Token
	src          *bufio.Reader
	rules        []LexicalRule
	syntax_rules []SyntacticRule
	debug        bool
    warnings, errors []*GomoreError
}

func TokenTypeToString(i int) string {
	if i == T_UNDECIDED {
		return "T_UNDECIDED"
	}
	if i == T_INT {
		return "T_INT"
	}
	if i == T_FLOAT {
		return "T_FLOAT"
	}
	if i == T_HEX {
		return "T_HEX"
	}
	if i == T_PERCENT {
		return "T_PERCENT"
	}
	if i == T_PIXELS {
		return "T_PIXELS"
	}
	if i == T_STRING {
		return "T_STRING"
	}
	if i == T_COMMENT {
		return "T_COMMENT"
	}
	if i == T_LINECOMMENT {
		return "T_LINECOMMENT"
	}
	if i == T_PATH {
		return "T_PATH"
	}
	if i == T_LPAREN {
		return "T_LPAREN"
	}
	if i == T_RPAREN {
		return "T_RPAREN"
	}
	if i == T_LBRACE {
		return "T_LBRACE"
	}
	if i == T_RBRACE {
		return "T_RBRACE"
	}
	if i == T_LBRACK {
		return "T_LBRACK"
	}
	if i == T_RBRACK {
		return "T_RBRACK"
	}
	if i == T_SELECTOR {
		return "T_SELECTOR"
	}
	if i == T_VARIABLE {
		return "T_VARIABLE"
	}
	if i == T_COLON {
		return "T_COLON"
	}
	if i == T_SEMICOLON {
		return "T_SEMICOLON"
	}
	if i == T_DASH {
		return "T_DASH"
	}
	if i == T_MUL {
		return "T_MUL"
	}
	if i == T_PLUS {
		return "T_PLUS"
	}
	if i == T_MINUS {
		return "T_MINUS"
	}
	if i == T_COMMA {
		return "T_COMMA"
	}
	if i == T_PERIOD {
		return "T_PERIOD"
	}
	if i == T_GT {
		return "T_GT"
	}
	if i == T_LT {
		return "T_LT"
	}
	if i == T_EQ {
		return "T_EQ"
	}
	if i == T_GUARD {
		return "T_GUARD"
	}
	return "UNKNOWN_TYPE"
}
func mergeslices(sls ...[]rune) []rune {
	var res []rune
	for _, sl := range sls {
		for _, slr := range sl {
			res = append(res, slr)
		}
	}
	return res
}
func (t *Token) AsString() string {
	var res string
	for _, r := range t.literal {
		res = fmt.Sprintf("%s%c", res, r)
	}
	return res

}
func SetupSyntaxRules(s *Scanner) {
	s.syntax_rules = append(s.syntax_rules, SyntacticRule{
		concerns: T_EQ,
		before:   []int{T_VARIABLE},
		after:    []int{T_SELECTOR, T_HEX, T_PIXELS, T_INT, T_FLOAT},
	})
}
func SetupLexicalRules(s *Scanner) {
	digits := []rune{'1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
	lower_letters := []rune{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z'}
	var alpha []rune
	for _, r := range lower_letters {
		alpha = append(alpha, r)
		alpha = append(alpha, unicode.ToUpper(r))
	}
	var alphanum []rune
	for _, r := range alpha {
		alphanum = append(alphanum, r)
	}
	for _, d := range digits {
		alphanum = append(alphanum, d)
	}
	var identifier []rune

	//empty := []rune{}
	cash := []rune{'$'}
	//at := []rune{'@'}
	dot := []rune{'.'}
	hash := []rune{'#'}
	hyphens := []rune{'-', '_'}
	identifier = mergeslices(alphanum, hyphens)
	s.rules = append(s.rules, LexicalRule{
		t:           T_PERCENT,
		starts_with: digits,
		chars_in:    mergeslices(digits, []rune{'%'}),
		ends_with:   []rune{'%'},
	})
	s.rules = append(s.rules, LexicalRule{
		t:           T_PIXELS,
		starts_with: digits,
		chars_in:    mergeslices(digits, []rune{'p', 'x'}),
		ends_with:   []rune{'x'},
		must_be_true: func(sc *Scanner, tk *Token) bool {
			l := len(tk.literal)
			if tk.literal[l-1] == 'x' && tk.literal[l-2] == 'p' {
				return true
			}
			return false
		},
	})
	s.rules = append(s.rules, LexicalRule{
		t:           T_INT,
		starts_with: digits,
		chars_in:    digits,
		ends_with:   digits,
	})
	s.rules = append(s.rules, LexicalRule{
		t:           T_VARIABLE,
		starts_with: cash,
		chars_in:    alphanum,
		ends_with:   alpha,
	})
	s.rules = append(s.rules, LexicalRule{
		t:           T_SELECTOR, // in PL parlance identifier
		starts_with: mergeslices(dot, hash, alpha),
		chars_in:    identifier,
		ends_with:   alphanum,
	})
	// do solitaries
	solitaries := []struct {
		ival         int
		char         rune
		reclassifier Validator
	}{
		{T_LBRACK, '[', nil},
		{T_RBRACK, ']', nil},
		{T_LBRACE, '{', nil},
		{T_RBRACE, '}', nil},
		{T_LPAREN, '(', nil},
		{T_RPAREN, ')', nil},
		{T_SEMICOLON, ';', nil},
		{T_COLON, ':', nil},
		{T_COMMA, ',', nil},
		{T_EQ, '=', nil},
		{T_PLUS, '+', nil},
		{T_MUL, '*', nil},
		{T_DIV, '/', nil},
	}
	for _, sol := range solitaries {
		s.rules = append(s.rules, LexicalRule{
			t:            sol.ival, // in PL parlance identifier
			starts_with:  []rune{sol.char},
			reclassifier: sol.reclassifier,
		})
	}
}
func InSlice(s []rune, r rune) bool {
	for _, c := range s {
		if r == c {
			return true
		}
	}
	return false
}
func InIntSlice(s []int, r int) bool {
	for _, c := range s {
		if r == c {
			return true
		}
	}
	return false
}
func (rule *LexicalRule) IsSatisfied(s *Scanner, t *Token) bool {
	if len(t.literal) == 0 {
		return false // not possible to have a 0 len token
	}
	if rule.must_be_true != nil {
		if rule.must_be_true(s, t) == false {
			s.Debugs(fmt.Sprintf("%s fails %s because of must_be_true\n", t.AsString(), TokenTypeToString(rule.t)))
			return false
		}
	}
	l_len := len(t.literal) - 1
	for i, cc := range t.literal {
		if i == 0 && len(rule.starts_with) > 0 {
			if InSlice(rule.starts_with, cc) == false {
				s.Debugs(fmt.Sprintf("%s fails %s because of starts_with\n", t.AsString(), TokenTypeToString(rule.t)))
				return false
			} else {
				continue
			}
		}
		if i == l_len && len(rule.ends_with) > 0 {
			if InSlice(rule.ends_with, cc) == false {
				s.Debugs(fmt.Sprintf("%s fails %s because of ends_with\n", t.AsString(), TokenTypeToString(rule.t)))
				return false
			} else {
				continue
			}
		}
		if len(rule.chars_in) > 0 {
			if InSlice(rule.chars_in, cc) == false {
				s.Debugs(fmt.Sprintf("%s fails %s because of chars_in\n", t.AsString(), TokenTypeToString(rule.t)))
				return false
			}
		}
		if len(rule.chars_not_in) > 0 {
			if InSlice(rule.chars_not_in, cc) == true {
				s.Debugs(fmt.Sprintf("%s fails %s because of chars_not_in\n", t.AsString(), TokenTypeToString(rule.t)))
				return false
			}
		}
	}

	return true
}

func (s *Scanner) Classify(t *Token) {
	for _, rule := range s.rules {
		if rule.IsSatisfied(s, t) {
			s.Debugs(fmt.Sprintf("Classify matched %s\n", TokenTypeToString(rule.t)))
			t.of_type = rule.t
			return
		}

	}
}
func (s *Scanner) SetSource(r *bufio.Reader) {
	s.Debugs("Setting source\n")
	s.src = r
}
func (s *Scanner) LastToken() *Token {
	l := len(s.tokens)
	if l > 0 {
		return s.tokens[l-1]
	}
	return nil
}
func (s *Scanner) Append(t *Token) {
	lt := s.LastToken()
	if lt == nil {
		t.previous = lt
		lt.next = t
	}
	s.tokens = append(s.tokens, t)
}
func (s *Scanner) TokenAt(i int) *Token {
	if len(s.tokens) > 0 && len(s.tokens) > i {
		return s.tokens[i]
	}
	return nil
}
func NewScanner() *Scanner {
	s := Scanner{
		tokens: []*Token{{of_type: T_GUARD}},
		debug:  false,
	}
	return &s
}
func (s *Scanner) SetDebug(v bool) {
	s.debug = v
}
func (s *Scanner) Debugs(str string) {
	if s.debug {
		fmt.Printf("%s\n", str)
	}

}
func (s *Scanner) Scan(r *bufio.Reader) error {
	s.SetSource(r)
	s.NextChar() // preload the head
	for {
		t := s.Next()
		if t == nil {
			break
		}
		s.Append(t)
	}
	return nil
}
func (s *Scanner) Next() *Token {
	if s.ch == 0 { // end of input
		s.Debugs("s.ch == 0 so returning nil")
		return nil
	}
	s.ConsumeWhiteSpace()
	t := &Token{}
	t.of_type = T_UNDECIDED
	t.col = s.col
	t.row = s.row
	in_string := false
	was_string := false
	for {
		if s.ch == 0 {
			break // end of input, but we still may have a final token
		}
		if s.ch == '"' {
			if in_string {
				in_string = false
				was_string = true
				s.Debugs("setting in_string to false and was_string to true")
				s.NextChar()
				break
			} else {
				s.Debugs("setting in_string to true")
				in_string = true
			}
			s.NextChar()
		}
		if s.ch == '\\' {
			s.Debugs("found escape char")
			s.NextChar()
		}
		s.Debugs(fmt.Sprintf("s.ch is '%c'\n", s.ch))
		s.Debugs(fmt.Sprintf("Appending '%c'\n", s.ch))
		t.literal = append(t.literal, s.ch)
		if IsSolitary(s.ch) && in_string != true {
			s.NextChar()
			break
		}
		s.NextChar()
		if IsTokenBreak(s.ch) && in_string != true {
			break
		}
	}
	if was_string == true {
		s.Debugs("setting of_type to T_STRING")
		t.of_type = T_STRING
	} else {
		s.Debugs("was_string false, so classifying")
		s.Classify(t)
	}
	return t
}

func IsTokenBreak(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	}
	breakables := []rune{':',
		';',
		'\n',
		'.',
		'(',
		')',
		',',
		'{',
		'}',
		'[',
		']',
	}
	for _, bc := range breakables {
		if bc == r {
			return true
		}
	}
	return false
}
func IsSolitary(r rune) bool {
	if unicode.IsSpace(r) {
		return true
	}
	breakables := []rune{':',
		';',
		'(',
		')',
		',',
		'{',
		'}',
		'[',
		']',
	}
	for _, bc := range breakables {
		if bc == r {
			return true
		}
	}
	return false
}

func (s *Scanner) ConsumeWhiteSpace() {
	s.Debugs("about to consume whitespace")
	for unicode.IsSpace(s.ch) {
		s.NextChar()
	}
}
func (s *Scanner) NextChar() {
	r, _, err := s.src.ReadRune()
	s.Debugs(fmt.Sprintf("r = '%c'\n", r))
	s.col = s.col + 1
	if r == '\n' {
		s.row = s.row + 1
		s.col = 0
	}
    s.lch = s.ch
    nr,_ := s.src.Peek(3)
    s.nch,width = utf8.DecodeRuneInString(string(nr[0:]))
	s.ch = r
	if err != nil {
		s.ch = 0
	}
}
