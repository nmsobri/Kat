package parser

// when you at the beginning of parsing, and you want to parse the
//whole expression, precedence should be 0

// when you at the beginning of parsing, and you want to parse until
//certain operator, precedence should be of that operator precedence

// if you at the middle of expression and you want to call the parse
//expression function again and make the expression left associative,
//precedence should be : previous operator precedence + 1

// if you at the middle of expression and you want to call the parse
//expression function again and make the expression right associative,
//precedence should be : previous operator precedence - 1

import (
	"kat/ast"
	"kat/lexer"
	"kat/token"
	"log"
	"strconv"
)

type StatementParselet func() ast.Stmt
type PrefixParselet func() ast.Expr
type InfixParselet func(left ast.Expr) ast.Expr

type Parser struct {
	Lex                *lexer.Lexer
	Token              token.Token
	NextToken          token.Token
	PrefixFunctions    map[token.TokenType]PrefixParselet
	InfixFunctions     map[token.TokenType]InfixParselet
	StatementFunctions map[token.TokenType]StatementParselet
}

func New(lex *lexer.Lexer) *Parser {
	p := &Parser{
		Lex:                lex,
		PrefixFunctions:    make(map[token.TokenType]PrefixParselet),
		InfixFunctions:     make(map[token.TokenType]InfixParselet),
		StatementFunctions: make(map[token.TokenType]StatementParselet),
	}

	// Regisgter statement functions
	p.StatementFunctions[token.CONST] = p.ParseConstDecl
	p.StatementFunctions[token.STRUCT] = p.ParseNodeStruct
	p.StatementFunctions[token.FUNCTION] = p.ParseNodeFunction
	p.StatementFunctions[token.LET] = p.ParseLetDecl
	p.StatementFunctions[token.IF] = p.ParseIfStmt
	p.StatementFunctions[token.FOR] = p.parseForStmt
	p.StatementFunctions[token.RETURN] = p.parseReturnStmt

	// Register Prefix functions
	p.PrefixFunctions[token.SELF] = p.ParseSelf
	p.PrefixFunctions[token.INTEGER] = p.ParseNodeDigit
	p.PrefixFunctions[token.DOUBLE] = p.ParseNodeDouble
	p.PrefixFunctions[token.TRUE] = p.ParseNodeBoolean
	p.PrefixFunctions[token.FALSE] = p.ParseNodeBoolean
	p.PrefixFunctions[token.MINUS] = p.ParsePrefixExpr
	p.PrefixFunctions[token.BANG] = p.ParsePrefixExpr
	p.PrefixFunctions[token.STRING] = p.ParseNodeString
	p.PrefixFunctions[token.IDENTIFIER] = p.ParseIdentifier
	p.PrefixFunctions[token.LBRACKET] = p.ParseArrayDecl
	p.PrefixFunctions[token.LBRACE] = p.ParseMapExpr
	p.PrefixFunctions[token.MINUSMINUS] = p.ParsePrefixExpr
	p.PrefixFunctions[token.PLUSPLUS] = p.ParsePrefixExpr
	p.PrefixFunctions[token.IMPORT] = p.ParseImportDecl

	// Register Infix functions
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
	p.InfixFunctions[token.LPAREN] = p.ParseFunctionCall
	p.InfixFunctions[token.LBRACKET] = p.ParseIndexExpr
	p.InfixFunctions[token.LBRACE] = p.ParseStructExpr
	p.InfixFunctions[token.MINUSMINUS] = p.parsePostfixExpr
	p.InfixFunctions[token.PLUSPLUS] = p.parsePostfixExpr
	p.InfixFunctions[token.EQUALEQUAL] = p.ParseBinaryExpr
	p.InfixFunctions[token.NOTEQUAL] = p.ParseBinaryExpr
	p.InfixFunctions[token.DOT] = p.ParseBinaryExpr

	p.NextToken = p.Lex.NextToken()
	p.Token = p.NextToken

	return p
}

func (p *Parser) ConsumeToken() token.Token {
	p.Token = p.NextToken
	p.NextToken = p.Lex.NextToken()
	return p.Token
}

func (p *Parser) ExpectToken(tok token.TokenType) token.Token {
	if p.NextToken.Type != tok {
		log.Fatalf("Expect next token of type: %s `%s`, got: %s `%s` at line: %d, column:%d\n",
			tok, tok.Str(), p.NextToken.Type, p.NextToken.Value,
			p.NextToken.Row+1, p.NextToken.Col+1,
		)

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

func (p *Parser) PeekAhead(count int) token.Token {
	if count == 0 {
		return p.Token
	} else if count == 1 {
		return p.NextToken
	}

	return p.Lex.PeekToken(count - 1)
}

func (*Parser) GetOperatorPrecedence(tok token.Token) int {
	return token.GetPrecedence(tok)
}

func (p *Parser) ParseProgram() *ast.NodeProgram {
	program := &ast.NodeProgram{}

	for p.CurrentToken().Type != token.EOF {
		program.Body = append(program.Body, p.ParseStatement())

		p.skipEOL()

		if p.PeekToken().Type == token.EOF {
			p.ConsumeToken()
		}
	}

	return program
}

func (p *Parser) ParseExpression(currentPrecedence int) ast.Expr {
	p.ConsumeToken()
	prefixFunction, ok := p.PrefixFunctions[p.CurrentToken().Type]

	if !ok {
		log.Fatalf("Could not parse prefix token: %s, value: `%s` at line: %d, column: %d",
			p.CurrentToken().Type, p.CurrentToken().Value, p.CurrentToken().Row+1,
			p.CurrentToken().Col+1,
		)
	}

	left := prefixFunction()

	for p.PeekToken().Type != token.EOL && p.GetOperatorPrecedence(p.PeekToken()) > currentPrecedence {
		p.ConsumeToken() // consume the infix operator

		infixFunction, ok := p.InfixFunctions[p.CurrentToken().Type]

		if !ok {
			log.Fatalf("Could not parse infix token: %s, value: `%s` at line: %d, column: %d",
				p.CurrentToken().Type, p.CurrentToken().Value, p.CurrentToken().Row+1,
				p.CurrentToken().Col+1,
			)
		}

		left = infixFunction(left)
	}

	return left
}

func (p *Parser) ParseNodeDigit() ast.Expr {
	val, e := strconv.ParseInt(p.CurrentToken().Value, 10, 64)

	if e != nil {
		log.Fatalf("Parser::Errors:%s\n", e)
	}

	return &ast.NodeInteger{
		Token: p.CurrentToken(),
		Value: val,
	}
}

func (p *Parser) ParseNodeDouble() ast.Expr {
	val, e := strconv.ParseFloat(p.CurrentToken().Value, 64)

	if e != nil {
		log.Fatalf("Parser::Errors:%s\n", e)
	}

	return &ast.NodeFloat{
		Token: p.CurrentToken(),
		Value: val,
	}
}

func (p *Parser) ParseNodeNil() ast.Expr {
	return &ast.NodeInteger{
		Token: p.CurrentToken(),
		Value: 999,
	}
}

func (p *Parser) ParseBinaryExpr(left ast.Expr) ast.Expr {
	currentToken := p.CurrentToken()

	right := p.ParseExpression(p.GetOperatorPrecedence(currentToken))

	return &ast.NodeBinaryExpr{
		Token:    currentToken,
		Left:     left,
		Right:    right,
		Operator: currentToken.Value,
	}
}

func (p *Parser) ParsePrefixExpr() ast.Expr {
	currentToken := p.CurrentToken() // the prefix

	right := p.ParseExpression(token.Precedence.PREFIX)

	if currentToken.Value == "-" {
		currentToken.Type = token.NEGATE
	}

	return &ast.NodePrefixExpr{
		Token:    currentToken,
		Operator: currentToken.Value,
		Right:    right,
	}
}

func (p *Parser) ParseNodeBoolean() ast.Expr {
	val := true

	if p.CurrentToken().Type == token.FALSE {
		val = false
	}

	return &ast.NodeBoolean{
		Token: p.CurrentToken(),
		Value: val,
	}
}

func (p *Parser) ParseIdentifier() ast.Expr {
	currentToken := p.CurrentToken()
	identifier := currentToken.Value

	return &ast.NodeIdentifier{
		Token: currentToken,
		Name:  identifier,
	}
}

func (p *Parser) ParseSelf() ast.Expr {
	currentToken := p.CurrentToken()
	identifier := currentToken.Value

	return &ast.NodeSelf{
		Token: p.CurrentToken(),
		Name:  identifier,
	}
}

func (p *Parser) ParseConditionExpr(left ast.Expr) ast.Expr {
	currentToken := p.CurrentToken()
	thenArm := p.ParseExpression(p.GetOperatorPrecedence(currentToken))

	p.ExpectToken(token.COLON) // consume the `:`

	elseArm := p.ParseExpression(p.GetOperatorPrecedence(currentToken) - 1)

	return &ast.NodeTernaryExpr{
		Token:     currentToken,
		Condition: left,
		ThenArm:   thenArm,
		ElseArm:   elseArm,
	}
}

func (p *Parser) ParseConstDecl() ast.Stmt {
	currentToken := p.CurrentToken()
	identifier := p.ParseExpression(token.Precedence.ASSIGNMENT)

	p.ExpectToken(token.EQUAL) // consume `=`

	value := p.ParseExpression(token.Precedence.LOWEST)

	return &ast.NodeConstStmt{
		Token:      currentToken,
		Identifier: identifier,
		Value:      value,
	}
}

func (p *Parser) ParseImportDecl() ast.Expr {
	currentToken := p.CurrentToken()
	p.ExpectToken(token.LPAREN) // consume `(`

	path := p.ParseExpression(token.Precedence.LOWEST)

	p.ExpectToken(token.RPAREN) // consume `)`

	return &ast.NodeImportExpr{
		Token: currentToken,
		Path:  path,
	}
}

func (p *Parser) ParseNodeString() ast.Expr {
	v, _ := strconv.Unquote(p.CurrentToken().Value)

	return &ast.NodeString{
		Token: p.CurrentToken(),
		Value: v,
	}
}

func (p *Parser) ParseNodeStruct() ast.Stmt {
	currentToken := p.CurrentToken()

	identifier := p.ParseExpression(token.Precedence.EXPR)

	p.skipEOL()

	p.ExpectToken(token.LBRACE) // consume `{`

	structProperties := make([]ast.Expr, 0)

	for p.PeekToken().Type != token.RBRACE {
		p.skipEOL()

		if p.PeekToken().Type == token.RBRACE {
			break
		}

		identifier := p.ParseExpression(token.Precedence.LOWEST)

		structProperties = append(structProperties, identifier)

		if p.PeekToken().Type == token.COMMA {
			p.ConsumeToken() // consume `,`
		}
	}

	p.ExpectToken(token.RBRACE) // consume `}`

	return &ast.NodeStructStmt{
		Token:      currentToken,
		Identifier: identifier,
		Properties: structProperties,
	}
}

func (p *Parser) ParseNodeFunction() ast.Stmt {
	currentToken := p.CurrentToken()
	_identifier := p.ParseExpression(token.Precedence.EXPR)
	var identifier ast.Expr = _identifier.(*ast.NodeIdentifier)

	if p.PeekToken().Type == token.DOT {
		p.ExpectToken(token.DOT)        // consume `.`
		p.ExpectToken(token.IDENTIFIER) // consume the identifier

		identifier = &ast.NodeBinaryExpr{
			Token:    currentToken,
			Left:     identifier,
			Right:    &ast.NodeIdentifier{Token: p.CurrentToken(), Name: p.CurrentToken().Value},
			Operator: ".",
		}
	}

	p.ExpectToken(token.LPAREN)
	arguements := p.ParseNodeFunctionArguement()
	p.ExpectToken(token.RPAREN)

	body := p.parseBlockStmt()

	return &ast.NodeFunctionStmt{
		Token:      currentToken,
		Identifier: identifier,
		Arguements: arguements,
		Body:       body,
	}
}

func (p *Parser) ParseNodeFunctionArguement() []ast.Expr {
	arguements := make([]ast.Expr, 0)

	for p.PeekToken().Type != token.RPAREN {
		identifier := p.ParseExpression(token.Precedence.LOWEST)
		arguements = append(arguements, identifier)

		if p.PeekToken().Type == token.COMMA {
			p.ExpectToken(token.COMMA) // consume `,`
		}
	}

	return arguements
}

func (p *Parser) ParseFunctionCall(left ast.Expr) ast.Expr {
	currentToken := p.CurrentToken()
	functionArgs := p.ParseNodeFunctionArguement()
	p.ExpectToken(token.RPAREN)

	return &ast.NodeFunctionCall{
		Token:      currentToken,
		Identifer:  left,
		Parameters: functionArgs,
	}
}

func (p *Parser) ParseLetDecl() ast.Stmt {
	currentToken := p.CurrentToken()
	ident := p.ParseExpression(token.Precedence.ASSIGNMENT)
	p.ExpectToken(token.EQUAL)
	value := p.ParseExpression(token.Precedence.LOWEST)

	return &ast.NodeLetStmt{
		Token:      currentToken,
		Identifier: ident,
		Value:      value,
	}
}

func (p *Parser) ParseArrayDecl() ast.Expr {
	currentToken := p.CurrentToken()

	values := make([]ast.Expr, 0)

	for p.PeekToken().Type != token.RBRACKET {
		values = append(values, p.ParseExpression(token.Precedence.LOWEST))

		if p.PeekToken().Type == token.COMMA {
			p.ExpectToken(token.COMMA)
		}
	}

	p.ExpectToken(token.RBRACKET)

	return &ast.NodeArrayExpr{
		Token: currentToken,
		Value: values,
	}
}

func (p *Parser) ParseIndexExpr(left ast.Expr) ast.Expr {
	currentToken := p.CurrentToken()
	index := p.ParseExpression(token.Precedence.LOWEST)
	p.ExpectToken(token.RBRACKET)

	return &ast.NodeIndexExpr{
		Token:      currentToken,
		Identifier: left,
		Index:      index,
	}
}

func (p *Parser) ParseMapExpr() ast.Expr {
	currentToken := p.CurrentToken()
	values := make(map[ast.Expr]ast.Expr, 0)

	for p.PeekToken().Type != token.RBRACE {
		ident := p.ParseExpression(token.Precedence.LOWEST)
		p.ExpectToken(token.COLON)
		value := p.ParseExpression(token.Precedence.LOWEST)
		values[ident] = value

		if p.PeekToken().Type == token.COMMA {
			p.ExpectToken(token.COMMA)
		}
	}

	p.ExpectToken(token.RBRACE)

	return &ast.NodeMapExpr{
		Token: currentToken,
		Map:   values,
	}
}

func (p *Parser) ParseStructExpr(left ast.Expr) ast.Expr {
	currentToken := p.CurrentToken()
	values := p.ParseMapExpr()

	return &ast.NodeStructExpr{
		Token:  currentToken,
		Name:   left,
		Values: values,
	}
}

func (p *Parser) ParseIfStmt() ast.Stmt {
	currentToken := p.CurrentToken()

	condition := p.ParseExpression(token.Precedence.LOWEST + 1)

	p.skipEOL()

	nodeIf := &ast.NodeConditionalStmt{
		Token:     currentToken,
		Condition: condition,
		ThenArm:   p.parseBlockStmt(),
	}

	if p.PeekAhead(1).Type == token.ELSE && p.PeekAhead(2).Type != token.IF {
		p.ExpectToken(token.ELSE)
		nodeIf.ElseArm = p.parseBlockStmt()
		return nodeIf
	}

	if p.PeekAhead(1).Type == token.ELSE && p.PeekAhead(2).Type == token.IF {
		p.ExpectToken(token.ELSE)
		p.ExpectToken(token.IF)

		nodeIf.ElseArm = p.ParseIfStmt()

		return nodeIf
	}

	return nodeIf
}

func (p *Parser) skipEOL() {
	for p.PeekToken().Type == token.EOL {
		p.ConsumeToken() // consume EOL
	}
}

func (p *Parser) parseBlockStmt() ast.Stmt {
	p.ExpectToken(token.LBRACE)
	p.skipEOL()

	body := &ast.NodeBlockStmt{}

	for p.PeekToken().Type != token.RBRACE {
		body.Body = append(body.Body, p.ParseStatement())
		p.skipEOL()
	}

	p.ExpectToken(token.RBRACE)

	return body
}

func (p *Parser) parseForStmt() ast.Stmt {
	if p.PeekToken().Type == token.LET {
		return p.parseClassicForStmt()
	}

	return p.parseModernForStmt()
}

func (p *Parser) parsePostfixExpr(left ast.Expr) ast.Expr {
	return &ast.NodePostfixExpr{
		Token:    p.CurrentToken(),
		Left:     left,
		Operator: p.CurrentToken().Value,
	}
}

func (p *Parser) ParseStatement() ast.Stmt {
	if p.StatementFunctions[p.PeekToken().Type] != nil {
		p.ConsumeToken()
		return p.StatementFunctions[p.CurrentToken().Type]()
	}

	return p.ParseExpressionStatement()
}

func (p *Parser) ParseExpressionStatement() ast.Stmt {
	return &ast.NodeExprStmt{Expr: p.ParseExpression(token.Precedence.LOWEST)}
}

func (p *Parser) parseModernForStmt() ast.Stmt {
	currentToken := p.CurrentToken()
	condition := p.ParseExpression(token.Precedence.LOWEST + 1)
	body := p.parseBlockStmt()

	return &ast.NodeModernForStmt{
		Token:     currentToken,
		Condition: condition,
		Body:      body,
	}
}

func (p *Parser) parseClassicForStmt() ast.Stmt {
	currentToken := p.CurrentToken()
	preExpr := p.ParseStatement()
	p.ExpectToken(token.SEMICOLON)

	condition := p.ParseExpression(token.Precedence.LOWEST)
	p.ExpectToken(token.SEMICOLON)

	postExpr := p.ParseExpression(token.Precedence.LOWEST + 1)

	body := p.parseBlockStmt()

	return &ast.NodeClassicForStmt{
		Token:     currentToken,
		Condition: condition,
		PreExpr:   preExpr,
		PostExpr:  postExpr,
		Body:      body,
	}
}

func (p *Parser) parseReturnStmt() ast.Stmt {
	currentToken := p.CurrentToken()
	expr := p.ParseExpression(token.Precedence.LOWEST)

	return &ast.NodeReturnStmt{
		Token: currentToken,
		Value: expr,
	}
}
