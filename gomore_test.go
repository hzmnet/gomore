package gomore

import (
	"bufio"
	"strings"
	"testing"
	//"fmt"
)

func TestIsTokenBreak(t *testing.T) {
	tdata := []struct {
		c        rune
		expected bool
	}{
		{'\n', true},
		{'a', false},
		{' ', true},
		{';', true},
		{':', true},
		{'.', true},
	}
	for _, tt := range tdata {
		actual := IsTokenBreak(tt.c)
		if actual != tt.expected {
			t.Errorf("IsTokenBreak(%s): expected %s, actual %s", tt.c, tt.expected, actual)
		}
	}
	// Output: Moo
}
func TestConsumeWhiteSpace(t *testing.T) {
	data := []struct {
		s string
		e rune
	}{
		{" x", 'x'},
		{" ", 0},
		{"\n\n\n$", '$'},
	}
	var s Scanner
	for _, tt := range data {
		b := bufio.NewReader(strings.NewReader(tt.s))
		s.SetSource(b)
		s.NextChar()
		s.ConsumeWhiteSpace()

		if s.ch != tt.e {
			t.Errorf("TestConsumeWhiteSpace expected %c got '%c'", tt.e, s.ch)
		}

	}
}
func TestNext(t *testing.T) {
	data := []struct {
		s string
		e []rune
	}{
		{" xyz;", []rune{'x', 'y', 'z'}},
		{" .y24 ;", []rune{'.', 'y', '2'}},
	}
	var s Scanner
	for _, tt := range data {
		b := bufio.NewReader(strings.NewReader(tt.s))
		s.SetSource(b)
		s.NextChar()
		tok := s.Next()

		i := 0
		if tok.literal[i] != tt.e[i] {
			t.Errorf("TestNext expected %c got '%c'", tt.e[i], tok.literal[i])
		}
		i++
		if tok.literal[i] != tt.e[i] {
			t.Errorf("TestNext expected %c got '%c'", tt.e[i], tok.literal[i])
		}
		i++
		if tok.literal[i] != tt.e[i] {
			t.Errorf("TestNext expected %c got '%c'", tt.e[i], tok.literal[i])
		}
		tok = s.Next()
		if tok.literal[0] != ';' {
			t.Errorf("TestNext expected ; got '%c'", tok.literal[i])
		}

	}
}
func TestTokenClassifications(t *testing.T) {
	var s Scanner
	example_text := `#ident1 { font-size: 100%; color: blue; }
    body {
        background: #fff;
        margin-bottom: 2px;
        #ident1(font-size, color)
    }
    `

	table := []struct {
		exp int
	}{
		{T_SELECTOR},
		{T_LBRACE},
		{T_SELECTOR},
		{T_COLON},
		{T_PERCENT},
		{T_SEMICOLON},
		{T_SELECTOR},
		{T_COLON},
		{T_SELECTOR},
		{T_SEMICOLON},
		{T_RBRACE},

		{T_SELECTOR}, {T_LBRACE}, {T_SELECTOR}, {T_COLON}, {T_SELECTOR}, {T_SEMICOLON},
		{T_SELECTOR}, {T_COLON}, {T_PIXELS}, {T_SEMICOLON},
		{T_SELECTOR}, {T_LPAREN}, {T_SELECTOR}, {T_COMMA}, {T_SELECTOR}, {T_RPAREN},
		{T_RBRACE},
	}
	b := bufio.NewReader(strings.NewReader(example_text))
	s.SetSource(b)
	SetupLexicalRules(&s)
	s.NextChar()
	for i, tk_int := range table {
		tok := s.Next()
		if tk_int.exp != tok.of_type {
			t.Errorf("TestTokenClassifications expected %s got '%s' from %s on toke %d", TokenTypeToString(tk_int.exp), TokenTypeToString(tok.of_type), tok.AsString(), i)
		}
	}

}
func TestStringCollection(t *testing.T) {
	s := NewScanner()
	txt := `"$this = 1;"`
	s.SetDebug(true)
	b := bufio.NewReader(strings.NewReader(txt))
	SetupLexicalRules(s)
	//SetupSyntaxRules(s)
	err := s.Scan(b)
	if err != nil {
		t.Errorf("Err does not = nil")
	}
	tok := s.TokenAt(1) // 0 is a T_GUARD
	if tok.of_type != T_STRING {
		t.Errorf("Expected %s to have type T_STRING but got %s\n", tok.AsString(), TokenTypeToString(tok.of_type))
	}
}
func TestStatementSyntaxRule(t *testing.T) {
	s := NewScanner()

	txt := `$this = 1;`
	b := bufio.NewReader(strings.NewReader(txt))
	SetupLexicalRules(s)
	SetupSyntaxRules(s)
	err := s.Scan(b)
	if err != nil {
		t.Errorf("Err does not = nil")
	}
}
