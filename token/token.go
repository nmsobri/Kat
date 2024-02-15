package token

import "fmt"

type TokenType string

type Token struct {
	Row   int
	Col   int
	Value string
	Type  TokenType
}

func (t Token) String() string {
	return fmt.Sprintf(
		"NextToken{ Line: %d, Col: %d, Value: `%s`, Type: %s }\n",
		t.Row, t.Col, t.Value, t.Type,
	)
}

const (
	// Single character
	PLUS      = "PLUS"      // +
	MINUS     = "MINUS"     // -
	MULTIPLY  = "MULTIPLY " // *
	DIVIDE    = "DIVIDE"    // /
	MODULO    = "MODULO"    // %
	EQUAL     = "EQUAL"     // =
	LESS      = "LESS"      // <
	GREATER   = "GREATER"   // >
	LBRACKET  = "LBRACKET"  // [
	RBRACKET  = "RBRACKET"  // ]
	LBRACE    = "LBRACE"    // {
	RBRACE    = "RBRACE"    // }
	QUOTE     = "QUOTE"     // "
	COLON     = "COLON"     // :
	LPAREN    = "LPAREN"    // (
	RPAREN    = "RPAREN"    // )
	COMMA     = "COMMA"     // ,
	SEMICOLON = "SEMICOLON" // ;
	DOT       = "DOT"       // .

	// Double character
	PLUSPLUS     = "PLUSPLUS"     // ++
	MINUSMINUS   = "MINUSMINUS"   // --
	EQUALEQUAL   = "EQUALEQUAL"   // ==
	GREATEREQUAL = "GREATEREQUAL" // >=
	LESSEQUAL    = "LESSQUAL"     // <=

	// Literal
	STRING = "STRING"
	DIGIT  = "DIGIT"

	// Keyword
	LET        = "LET"        // let
	IF         = "IF"         // if
	ELSE       = "ELSE"       // else
	FOR        = "FOR"        // for
	SELF       = "SELF"       // self
	IMPORT     = "IMPORT"     // import
	STRUCT     = "STRUCT"     // struct
	FUNCTION   = "FUNCTION"   // fn
	IDENTIFIER = "IDENTIFIER" // any

	// Special
	EOL     = "EOL"     // End of line
	EOF     = "EOF"     // End of file
	INVALID = "INVALID" // End of file
)

func Symbol(key string) TokenType {
	keywords := map[string]TokenType{
		"let":    LET,
		"if":     IF,
		"else":   ELSE,
		"for":    FOR,
		"self":   SELF,
		"import": IMPORT,
		"struct": STRUCT,
		"fn":     FUNCTION,
	}

	keyword, ok := keywords[key]

	if ok {
		return keyword
	}

	return IDENTIFIER
}
