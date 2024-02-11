package token

import "fmt"

type Token struct {
	Row   int
	Col   int
	Value string
	Type  int
}

func (t Token) String() string {
	return fmt.Sprintf(
		"Token{ Row: %d, Col: %d, Value: %s, Type: %d }\n",
		t.Row, t.Col, t.Value, t.Type,
	)
}

const (
	// Single character
	PLUS     = "PLUS"      // +
	MINUS    = "MINUS"     // -
	MULTIPLY = "MULTIPLY " // *
	DIVIDE   = "DIVIDE"    // /
	MODULE   = "MODULE"    // %
	EQUAL    = "EQUAL"     // =
	LESS     = "LESS"      // <
	GREATER  = "GREATER"   // >
	LBRACKET = "LBRACKET"  // [
	RBRACKET = "RBRACKET"  // ]
	LBRACE   = "LBRACE"    // {
	RBRACE   = "RBRACE"    // }
	QUOTE    = "QUOTE"     // "
	COLON    = "COLON"     // :
	LPAREN   = "LPAREN"    // (
	RPAREN   = "RPAREN"    // )
	COMMA    = "COMMA"     // ,
	DOT      = "DOT"       // .

	// Double character
	PLUSPLUS   = "PLUSPLUS"   // ++
	MINUSMINUS = "MINUSMINUS" // --
	EQUALEQUAL = "EQUALEQUAL" // ==

	// Keywords
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
	EOF //
	EOL //
)

func Symbol(key string) int {
	keywords := map[string]int{
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
