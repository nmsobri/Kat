package token

import (
	"fmt"
)

const (
	// Single character
	PLUS      TokenType = iota // +
	MINUS                      // -
	NEGATE                     // -
	BANG                       // !
	QUESTION                   // ?
	MULTIPLY                   // *
	DIVIDE                     // /
	MODULO                     // %
	EQUAL                      // =
	LESS                       // <
	GREATER                    // >
	LBRACKET                   // [
	RBRACKET                   // ]
	LBRACE                     // {
	RBRACE                     // }
	COLON                      // :
	LPAREN                     // (
	RPAREN                     // )
	COMMA                      // ,
	SEMICOLON                  // ;
	DOT                        // .

	// Double character
	PLUSPLUS     // ++
	MINUSMINUS   // --
	EQUALEQUAL   // ==
	NOTEQUAL     // ==
	GREATEREQUAL // >=
	LESSEQUAL    // <=
	COMMENT      // //

	// Literal
	STRING
	INTEGER
	DOUBLE

	// Keyword
	TRUE       // true
	FALSE      // false
	LET        // let
	CONST      // const
	IF         // if
	ELSE       // else
	FOR        // for
	SELF       // self
	IMPORT     // import
	STRUCT     // struct
	FUNCTION   // fn
	RETURN     // return
	IDENTIFIER // any

	// Special
	EOL                 // End of line
	EOF                 // End of file
	INVALID             // Invalid
	ANNOTATIONPRIMITIVE // Primitive annotation
	ANNOTATIONARRAY     // Array annotation
	ANNOTATIONMAP       // Map annotation
	ANNOTATIONSTRUCT    // Struct annotation
)

var (
	TokenString = map[TokenType]string{
		PLUS:                "+",
		MINUS:               "-",
		NEGATE:              "-",
		BANG:                "!",
		QUESTION:            "?",
		MULTIPLY:            "*",
		DIVIDE:              "/",
		MODULO:              "%",
		EQUAL:               "=",
		LESS:                "<",
		GREATER:             ">",
		LBRACKET:            "[",
		RBRACKET:            "]",
		LBRACE:              "{",
		RBRACE:              "}",
		COLON:               ":",
		LPAREN:              "(",
		RPAREN:              ")",
		COMMA:               ",",
		SEMICOLON:           ";",
		DOT:                 ".",
		PLUSPLUS:            "++",
		MINUSMINUS:          "--",
		EQUALEQUAL:          "==",
		GREATEREQUAL:        ">=",
		LESSEQUAL:           "<=",
		STRING:              "string",
		INTEGER:             "integer",
		DOUBLE:              "double",
		TRUE:                "true",
		FALSE:               "false",
		LET:                 "let",
		CONST:               "const",
		IF:                  "if",
		ELSE:                "else",
		FOR:                 "for",
		SELF:                "self",
		IMPORT:              "import",
		STRUCT:              "struct",
		FUNCTION:            "function",
		IDENTIFIER:          "identifier",
		EOL:                 "eol",
		EOF:                 "eof",
		INVALID:             "invalid",
		ANNOTATIONPRIMITIVE: "annotationprimitive",
		ANNOTATIONARRAY:     "annotationarray",
		ANNOTATIONMAP:       "annotationmap",
		ANNOTATIONSTRUCT:    "annotationstruct",
	}

	Keywords = map[string]TokenType{
		"true":   TRUE,
		"false":  FALSE,
		"let":    LET,
		"const":  CONST,
		"if":     IF,
		"else":   ELSE,
		"for":    FOR,
		"self":   SELF,
		"import": IMPORT,
		"struct": STRUCT,
		"fn":     FUNCTION,
		"return": RETURN,
	}

	Precedence = struct {
		LOWEST      int
		ASSIGNMENT  int
		CONDITIONAL int
		COMPARISON  int
		SUM         int
		PRODUCT     int
		EXPONENT    int
		PREFIX      int
		POSTFIX     int
		EXPR        int
	}{
		LOWEST:      0,
		ASSIGNMENT:  1,
		CONDITIONAL: 2,
		COMPARISON:  3,
		SUM:         4,
		PRODUCT:     5,
		EXPONENT:    6,
		PREFIX:      7,
		POSTFIX:     8,
		EXPR:        9,
	}
)

type (
	TokenType int

	Token struct {
		Row   int
		Col   int
		Value string
		Type  TokenType
	}
)

func (t Token) String() string {
	return fmt.Sprintf(
		"Token{ Line: %d, Col: %d, TokenString: %s, Value: `%s` }",
		t.Row, t.Col, t.Type, t.Value,
	)
}

func (tt TokenType) String() string {
	if tokenString, ok := TokenString[tt]; ok {
		return tokenString
	}

	return TokenString[INVALID]
}

func Keyword(key string) TokenType {
	keyword, ok := Keywords[key]

	if ok {
		return keyword
	}

	return IDENTIFIER
}

// Precedence is only for infix expression i guess
func GetPrecedence(tok Token) int {
	precedences := map[TokenType]int{
		LBRACE: Precedence.ASSIGNMENT,
		EQUAL:  Precedence.ASSIGNMENT,

		LESS:         Precedence.COMPARISON,
		GREATER:      Precedence.COMPARISON,
		LESSEQUAL:    Precedence.COMPARISON,
		GREATEREQUAL: Precedence.COMPARISON,
		EQUALEQUAL:   Precedence.COMPARISON,
		NOTEQUAL:     Precedence.COMPARISON,

		PLUS:  Precedence.SUM,
		MINUS: Precedence.SUM,

		MULTIPLY: Precedence.PRODUCT,
		DIVIDE:   Precedence.PRODUCT,
		MODULO:   Precedence.PRODUCT,

		NEGATE: Precedence.PREFIX,
		BANG:   Precedence.PREFIX,

		QUESTION:   Precedence.CONDITIONAL,
		LPAREN:     Precedence.EXPR,
		LBRACKET:   Precedence.EXPR,
		DOT:        Precedence.EXPR,
		MINUSMINUS: Precedence.PREFIX,
		PLUSPLUS:   Precedence.PREFIX,
	}

	precedence, ok := precedences[tok.Type]

	if ok {
		return precedence
	}

	return 0
}

/*
Language features that need to implement type checking
constant / variable declaration
	primitive, array, map

function declaration
	function arguement, function return

struct declaration
*/
