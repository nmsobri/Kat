package parser

import (
	"kat/ast"
	"kat/lexer"
	"kat/token"
	"log"
	"strconv"
)

type PrefixParselet func() ast.Node
type InfixParselet func(left ast.Node) ast.Node

type Parser struct {
	Lex             *lexer.Lexer
	Token           token.Token
	NextToken       token.Token
	PrefixFunctions map[token.TokenType]PrefixParselet
	InfixFunctions  map[token.TokenType]InfixParselet
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{
		Lex:             lex,
		PrefixFunctions: make(map[token.TokenType]PrefixParselet),
		InfixFunctions:  make(map[token.TokenType]InfixParselet),
	}

	// Register Prefixes
	p.PrefixFunctions[token.SELF] = p.ParseIdentifier
	p.PrefixFunctions[token.INTEGER] = p.ParseNodeDigit
	p.PrefixFunctions[token.DOUBLE] = p.ParseNodeDouble
	p.PrefixFunctions[token.TRUE] = p.ParseNodeBoolean
	p.PrefixFunctions[token.FALSE] = p.ParseNodeBoolean
	p.PrefixFunctions[token.MINUS] = p.ParsePrefixExpr
	p.PrefixFunctions[token.BANG] = p.ParsePrefixExpr
	p.PrefixFunctions[token.IDENTIFIER] = p.ParseIdentifier
	p.PrefixFunctions[token.CONST] = p.ParseConstDecl
	p.PrefixFunctions[token.IMPORT] = p.ParseImportDecl
	p.PrefixFunctions[token.STRING] = p.ParseNodeString
	p.PrefixFunctions[token.STRUCT] = p.ParseNodeStruct
	p.PrefixFunctions[token.FUNCTION] = p.ParseNodeFunction

	// Register Infixes
	p.InfixFunctions[token.PLUS] = p.ParseBinaryExpr
	p.InfixFunctions[token.MINUS] = p.ParseBinaryExpr
	p.InfixFunctions[token.MULTIPLY] = p.ParseBinaryExpr
	p.InfixFunctions[token.DIVIDE] = p.ParseBinaryExpr
	p.InfixFunctions[token.MODULO] = p.ParseBinaryExpr
	p.InfixFunctions[token.QUESTION] = p.ParseConditionExpr
	p.InfixFunctions[token.LESS] = p.ParseBinaryExpr
	p.InfixFunctions[token.LESSEQUAL] = p.ParseBinaryExpr
	p.InfixFunctions[token.GREATER] = p.ParseBinaryExpr
	p.InfixFunctions[token.GREATEREQUAL] = p.ParseBinaryExpr
	p.InfixFunctions[token.EQUAL] = p.ParseBinaryExpr

	p.NextToken = p.Lex.NextToken()

	return p
}

func (p *Parser) ConsumeToken() token.Token {
	p.Token = p.NextToken
	p.NextToken = p.Lex.NextToken()
	return p.Token
}

func (p *Parser) ExpectToken(tok token.TokenType) token.Token {
	if p.NextToken.Type != tok {
		log.Fatalf("Expect next token of type: %s, got: %s\n", tok, p.NextToken)
		return token.Token{}
	}

	return p.ConsumeToken()
}

func (p *Parser) CurrentToken() token.Token {
	return p.Token
}

func (p *Parser) PeekToken() token.Token {
	return p.NextToken
}

func (*Parser) GetOperatorPrecedence(tok token.Token) int {
	return token.GetPrecedence(tok)
}

func (p *Parser) ParseProgram() ast.NodeProgram {
	program := ast.NodeProgram{}

	for p.CurrentToken().Type != token.EOF {
		program.Body = append(program.Body, p.ParseExpression(0))

		for p.PeekToken().Type == token.EOL {
			p.ConsumeToken() // consume EOL
		}

		if p.PeekToken().Type == token.EOF {
			p.ConsumeToken()
		}
	}

	return program
}

func (p *Parser) ParseExpression(currentPrecedence int) ast.Node {
	p.ConsumeToken()
	prefixFunction, ok := p.PrefixFunctions[p.CurrentToken().Type]

	if !ok {
		log.Fatalf("Could not parse token: %s, value: %s", p.CurrentToken().Type, p.CurrentToken().Value)
	}

	left := prefixFunction()

	for p.PeekToken().Type != token.EOL && p.GetOperatorPrecedence(p.PeekToken()) > currentPrecedence {
		p.ConsumeToken() // consume the infix operator

		infixFunction, ok := p.InfixFunctions[p.CurrentToken().Type]

		if !ok {
			log.Fatalf("Could not parse token: %s", p.CurrentToken().Type)
		}

		left = infixFunction(left)
	}

	return left
}

func (p *Parser) ParseNodeDigit() ast.Node {
	val, e := strconv.ParseInt(p.CurrentToken().Value, 10, 64)

	if e != nil {
		log.Fatalf("Parser::Error:%s\n", e)
	}

	return ast.NodeInteger{
		Token: p.CurrentToken(),
		Value: val,
	}
}

func (p *Parser) ParseNodeDouble() ast.Node {
	val, e := strconv.ParseFloat(p.CurrentToken().Value, 64)

	if e != nil {
		log.Fatalf("Parser::Error:%s\n", e)
	}

	return ast.NodeDouble{
		Token: p.CurrentToken(),
		Value: val,
	}
}

func (p *Parser) ParseNodeNil() ast.NodeInteger {
	return ast.NodeInteger{
		Token: p.CurrentToken(),
		Value: 999,
	}
}

func (p *Parser) ParseBinaryExpr(left ast.Node) ast.Node {
	currentToken := p.CurrentToken()

	right := p.ParseExpression(p.GetOperatorPrecedence(currentToken))

	return ast.NodeBinaryExpr{
		Token:    currentToken,
		Left:     left,
		Right:    right,
		Operator: currentToken.Value,
	}
}

func (p *Parser) ParsePrefixExpr() ast.Node {
	currentToken := p.CurrentToken() // the prefix

	right := p.ParseExpression(token.Precedence.PREFIX)

	if currentToken.Value == "-" {
		currentToken.Type = token.NEGATE
	}

	return ast.NodeUnary{
		Token:    currentToken,
		Operator: currentToken.Value,
		Right:    right,
	}
}

func (p *Parser) ParseNodeBoolean() ast.Node {
	val := true

	if p.CurrentToken().Type == token.FALSE {
		val = false
	}

	return ast.NodeBoolean{
		Token: p.CurrentToken(),
		Value: val,
	}
}

func (p *Parser) ParseIdentifier() ast.Node {
	currentToken := p.CurrentToken()
	identifier := currentToken.Value

	if p.PeekToken().Type == token.DOT {
		p.ExpectToken(token.DOT) // consume `.`
		p.ParseExpression(0)

		identifier += "."
		identifier += p.CurrentToken().Value
	}

	return ast.NodeIdentifier{
		Token: currentToken,
		Name:  identifier,
	}
}

func (p *Parser) ParseConditionExpr(left ast.Node) ast.Node {
	currentToken := p.CurrentToken()
	thenArm := p.ParseExpression(p.GetOperatorPrecedence(currentToken))

	p.ConsumeToken() // consume the `:`

	elseArm := p.ParseExpression(p.GetOperatorPrecedence(currentToken) - 1)

	return ast.NodeConditionalExpr{
		Token:     currentToken,
		Condition: left,
		ThenArm:   thenArm,
		ElseArm:   elseArm,
	}
}

func (p *Parser) ParseConstDecl() ast.Node {
	currentToken := p.CurrentToken()
	identifier := p.ParseExpression(0)

	p.ConsumeToken() // consume `=`

	value := p.ParseExpression(0)

	return ast.NodeConstDeclaration{
		Token:      currentToken,
		Identifier: identifier,
		Value:      value,
	}
}

func (p *Parser) ParseImportDecl() ast.Node {
	currentToken := p.CurrentToken()
	p.ConsumeToken() // consume `(`

	path := p.ParseExpression(0)

	p.ConsumeToken() // consume `)`

	return ast.NodeImportDeclaration{
		Token: currentToken,
		Path:  path,
	}

}

func (p *Parser) ParseNodeString() ast.Node {
	return ast.NodeString{
		Token: p.CurrentToken(),
		Value: p.CurrentToken().Value,
	}
}

func (p *Parser) ParseNodeStruct() ast.Node {
	currentToken := p.CurrentToken()

	identifier := p.ParseExpression(0)

	p.ExpectToken(token.LBRACE) // consume `{`

	structProperties := ast.NodeStructProperties{
		Token:      currentToken,
		Properties: make([]ast.Node, 0),
	}

	for p.PeekToken().Type != token.RBRACE {
		if p.PeekToken().Type == token.EOL {
			p.ConsumeToken()
		}

		if p.PeekToken().Type == token.RBRACE {
			break
		}

		identifier := p.ParseExpression(0)

		structProperties.Properties = append(structProperties.Properties, identifier)

		if p.PeekToken().Type == token.COMMA {
			p.ConsumeToken() // consume `,`
		}
	}

	p.ExpectToken(token.RBRACE) // consume `}`

	return ast.NodeStruct{
		Token:      currentToken,
		Identifier: identifier,
		Properties: structProperties,
	}
}

func (p *Parser) ParseNodeFunction() ast.Node {
	currentToken := p.CurrentToken()

	identifier := p.ParseExpression(0)

	p.ExpectToken(token.LPAREN)

	arguements := p.ParseNodeFunctionArguement()

	p.ExpectToken(token.RPAREN)

	p.ExpectToken(token.LBRACE)

	if p.PeekToken().Type == token.EOL {
		p.ExpectToken(token.EOL) // consume `\n`
	}

	body := p.ParseNodeFunctionBody()

	p.ExpectToken(token.RBRACE)

	return ast.NodeFunction{
		Token:      currentToken,
		Identifier: identifier,
		Arguements: arguements,
		Body:       body,
	}
}

func (p *Parser) ParseNodeFunctionArguement() []ast.Node {
	arguements := make([]ast.Node, 0)

	for p.PeekToken().Type != token.RPAREN {
		identifier := p.ParseExpression(0)
		arguements = append(arguements, identifier)

		if p.PeekToken().Type == token.COMMA {
			p.ExpectToken(token.COMMA) // consume `,`
		}
	}

	return arguements
}

func (p *Parser) ParseNodeFunctionBody() []ast.Node {
	body := make([]ast.Node, 0)

	for p.PeekToken().Type != token.RBRACE {
		expression := p.ParseExpression(0)
		p.ExpectToken(token.EOL)
		body = append(body, expression)
	}

	return body
}
