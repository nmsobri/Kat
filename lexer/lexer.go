package lexer

import (
	"kat/token"
	"slices"
)

type Lexer struct {
	Col    int
	Line   int
	Offset int
	Input  []byte
}

func New(input []byte) *Lexer {
	return &Lexer{Col: 0, Line: 0, Offset: 0, Input: input}
}

func (l *Lexer) MakeToken(col int, val string, tokenType token.Type, tokenKind ...token.Type) token.Token {
	kind := token.Type(val)

	if len(tokenKind) > 0 {
		kind = tokenKind[0]
	}

	return token.Token{
		Line:     l.Line,
		Column:   col - l.Offset,
		Value:    val,
		Type:     tokenType,
		TypeKind: kind,
	}
}

func (l *Lexer) NextToken() token.Token {
	l.SkipWhitespace()
	ch := l.Char()
	var t token.Token

	switch ch {

	case '+':
		if l.PeekChar() == '+' {
			col := l.Col
			l.AdvanceChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.PLUS_PLUS)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.PLUS)
		}

	case '-':
		if l.PeekChar() == '-' {
			col := l.Col
			l.AdvanceChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.MINUS_MINUS)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.MINUS)
		}

	case '=':
		if l.PeekChar() == '=' {
			col := l.Col
			l.AdvanceChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.EQUAL_EQUAL)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.EQUAL)
		}

	case '!':
		if l.PeekChar() == '=' {
			col := l.Col
			l.AdvanceChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.NOT_EQUAL)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.BANG)
		}

	case '<':
		if l.PeekChar() == '=' {
			col := l.Col
			l.AdvanceChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.LESS_EQUAL)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.LESS)
		}

	case '>':
		if l.PeekChar() == '=' {
			col := l.Col
			l.AdvanceChar()
			t = l.MakeToken(col, string(l.Input[col:col+2]), token.GREATER_EQUAL)
		} else {
			t = l.MakeToken(l.Col, string(ch), token.GREATER)
		}

	case '*':
		t = l.MakeToken(l.Col, string(ch), token.MULTIPLY)

	case '/':
		if l.PeekChar() == '/' {
			l.AdvanceChar()
			l.AdvanceChar()

			for l.PeekChar() != '\n' && l.PeekChar() != 0 {
				l.AdvanceChar()
			}

			l.AdvanceChar()      // consume EOL or EOF
			l.AdvanceChar()      // advance to next character
			return l.NextToken() // advance to next token and return it
		}

		t = l.MakeToken(l.Col, string(ch), token.DIVIDE)

	case '%':
		t = l.MakeToken(l.Col, string(ch), token.MODULO)

	case '[':
		if l.PeekChar() == ']' { // []float, []int, []string
			l.AdvanceChar() // ]
			l.AdvanceChar() // type
			arrayType := l.MakeIdentifier()
			symbol := "array"
			t = l.MakeToken(l.Col, string(arrayType), token.Symbol(symbol), "array")
		} else {
			t = l.MakeToken(l.Col, string(ch), token.LBRACKET)
		}

	case ']':
		t = l.MakeToken(l.Col, string(ch), token.RBRACKET)

	case '{':
		t = l.MakeToken(l.Col, string(ch), token.LBRACE)

	case '}':
		t = l.MakeToken(l.Col, string(ch), token.RBRACE)

	case '(':
		t = l.MakeToken(l.Col, string(ch), token.LPAREN)

	case ')':
		t = l.MakeToken(l.Col, string(ch), token.RPAREN)

	case '?':
		t = l.MakeToken(l.Col, string(ch), token.QUESTION)

	case '"':
		col := l.Col
		str := l.MakeString()
		t = l.MakeToken(col, string(str), token.STRING)

	case ':':
		t = l.MakeToken(l.Col, string(ch), token.COLON)

	case ',':
		t = l.MakeToken(l.Col, string(ch), token.COMMA)

	case ';':
		t = l.MakeToken(l.Col, string(ch), token.SEMICOLON)

	case '.':
		t = l.MakeToken(l.Col, string(ch), token.DOT)

	case '\n':
		t = l.MakeToken(l.Col, "\\n", token.EOL)
		l.Offset = l.Col + 1
		l.Line++

	case 0:
		t = l.MakeToken(l.Col, "EOF", token.EOF)

	default:
		if l.IsDigit(ch) {
			col := l.Col
			dig := l.MakeDigit()

			var tok token.Type = token.INTEGER

			if slices.Contains(dig, 46) {
				tok = token.FLOAT
			}

			t = l.MakeToken(col, string(dig), tok)

		} else if l.IsChar(ch) {
			col := l.Col
			unknown := l.MakeIdentifier()
			symbol := string(unknown)
			t = l.MakeToken(col, symbol, token.Symbol(symbol))
		} else {
			col := l.Col
			invalid := l.MakeInvalid()
			t = l.MakeToken(col, string(invalid), token.INVALID)
		}
	}

	l.AdvanceChar()
	return t
}

func (l *Lexer) PeekToken(count int) token.Token {
	start := l.Col
	var t token.Token

	for i := 0; i < count; i++ {
		t = l.NextToken()
	}

	l.Col = start

	return t
}

func (l *Lexer) Char() byte {
	if l.Col < len(l.Input) {
		return l.Input[l.Col]
	}

	return 0
}

func (l *Lexer) AdvanceChar() {
	l.SkipWhitespace()
	l.Col++
}

func (l *Lexer) PeekChar() byte {
	col := l.Col
	l.Col++

	if l.Col < len(l.Input) {
		l.SkipWhitespace()
		ch := l.Input[l.Col]
		l.Col = col
		return ch
	}

	return 0
}

func (l *Lexer) IsChar(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

func (l *Lexer) IsDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func (l *Lexer) IsDouble(ch byte) bool {
	return ch >= '0' && ch <= '9' || ch == '.'
}

func (l *Lexer) IsAlphaNum(ch byte) bool {
	return l.IsChar(ch) || l.IsDigit(ch) || ch == '_'
}

func (l *Lexer) IsWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\r' || ch == '\n'
}

func (l *Lexer) IsEndOfString() bool {
	return l.Char() == '"' || l.Char() == 0
}

func (l *Lexer) SkipWhitespace() {
	for l.Col < len(l.Input) && l.IsWhitespace(l.Input[l.Col]) {
		l.Col++
	}
}

func (l *Lexer) MakeString() []byte {
	start := l.Col
	l.AdvanceChar() // skip first `"` so when we find next `"` it mark end of string

	for !l.IsEndOfString() {
		l.AdvanceChar()
	}

	end := l.Col
	return l.Input[start : end+1]
}

func (l *Lexer) MakeDigit() []byte {
	start := l.Col

	for l.IsDouble(l.Char()) {
		l.AdvanceChar()
	}

	end := l.Col
	l.Col--

	return l.Input[start:end]
}

func (l *Lexer) MakeIdentifier() []byte {
	start := l.Col

	for l.IsAlphaNum(l.Char()) {
		l.AdvanceChar()
	}

	end := l.Col
	l.Col--
	return l.Input[start:end]
}

func (l *Lexer) MakeInvalid() []byte {
	start := l.Col

	for l.Char() != 0 && !l.IsWhitespace(l.Char()) {
		l.AdvanceChar()
	}

	end := l.Col
	l.Col--
	return l.Input[start:end]
}
