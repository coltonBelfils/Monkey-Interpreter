package parser

import (
	"MonkeyInterpreter/ast"
	"MonkeyInterpreter/lexer"
	"MonkeyInterpreter/token"
	"fmt"
	"math/big"
)

const (
	_ int = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // < > <= >= //Should things like <= be here or with ==?
	SUM         // + -
	PRODUCT     // * / %
	PREFIX      // -x --x !x ++x
	POSTFIX     // x++ x-- x! // postfix is second so -x++ renders as (-x)++ and not -(x++).
	CALL        // x()
)

var precedences = map[token.TokenType]int{
	token.EQ:          EQUALS,
	token.NEQ:         EQUALS,
	token.LT_EQ:       LESSGREATER, // Or should it be EQUALS, or an intermediary value?
	token.GT_EQ:       LESSGREATER, // Or should it be EQUALS, or an intermediary value?
	token.LT:          LESSGREATER,
	token.GT:          LESSGREATER,
	token.PLUS:        SUM,
	token.MINUS:       SUM,
	token.ASTERISK:    PRODUCT,
	token.SLASH:       PRODUCT,
	token.PERCENT:     PRODUCT,
	token.AND:         SUM,
	token.OR:          SUM,
	token.PLUS_PLUS:   POSTFIX,
	token.MINUS_MINUS: POSTFIX,
	token.BANG:        POSTFIX,

	token.RIGHT_SHIFT: PRODUCT,
	token.LEFT_SHIFT:  PRODUCT,
	token.AMPERSAND:   PRODUCT,
	token.CAROT:       SUM,
	token.PIPE:        SUM,
}

type (
	prefixParseFn  func() ast.Expression
	infixParseFn   func(expression ast.Expression) ast.Expression
	postfixParseFn func(expression ast.Expression) ast.Expression
)

type Parser struct {
	l      *lexer.Lexer
	errors []string

	curToken  token.Token
	peekToken token.Token

	prefixParseFns  map[token.TokenType]prefixParseFn
	infixParseFns   map[token.TokenType]infixParseFn
	postfixParseFns map[token.TokenType]postfixParseFn
}

func postfixWrapper() ast.Expression {
	return nil
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.prefixParseFns = make(map[token.TokenType]prefixParseFn)
	p.registerPrefixFn(token.IDENT, p.parseIdentifier)
	p.registerPrefixFn(token.INT, p.parseIntegerLiteral)
	p.registerPrefixFn(token.DEC, p.parseDecimalLiteral)
	p.registerPrefixFn(token.BANG, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.PLUS_PLUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.MINUS_MINUS, p.parsePrefixExpression)
	p.registerPrefixFn(token.TILDE, p.parsePrefixExpression)
	p.registerPrefixFn(token.TRUE, p.parseBoolean)
	p.registerPrefixFn(token.FALSE, p.parseBoolean)
	p.registerPrefixFn(token.LPAREN, p.parseGroupedExpression)

	p.infixParseFns = make(map[token.TokenType]infixParseFn)
	p.registerInfixFn(token.PLUS, p.parseInfixExpression)
	p.registerInfixFn(token.MINUS, p.parseInfixExpression)
	p.registerInfixFn(token.SLASH, p.parseInfixExpression)
	p.registerInfixFn(token.ASTERISK, p.parseInfixExpression)
	p.registerInfixFn(token.PERCENT, p.parseInfixExpression)
	p.registerInfixFn(token.EQ, p.parseInfixExpression)
	p.registerInfixFn(token.NEQ, p.parseInfixExpression)
	p.registerInfixFn(token.LT_EQ, p.parseInfixExpression)
	p.registerInfixFn(token.GT_EQ, p.parseInfixExpression)
	p.registerInfixFn(token.LT, p.parseInfixExpression)
	p.registerInfixFn(token.GT, p.parseInfixExpression)
	p.registerInfixFn(token.AND, p.parseInfixExpression)
	p.registerInfixFn(token.OR, p.parseInfixExpression)
	p.registerInfixFn(token.AMPERSAND, p.parseInfixExpression)
	p.registerInfixFn(token.PIPE, p.parseInfixExpression)
	p.registerInfixFn(token.CAROT, p.parseInfixExpression)
	p.registerInfixFn(token.LEFT_SHIFT, p.parseInfixExpression)
	p.registerInfixFn(token.RIGHT_SHIFT, p.parseInfixExpression)

	// postfix is infix with the left ignored
	p.registerInfixFn(token.BANG, p.parsePostfixExpression)
	p.registerInfixFn(token.PLUS_PLUS, p.parsePostfixExpression)
	p.registerInfixFn(token.MINUS_MINUS, p.parsePostfixExpression)

	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) registerPrefixFn(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

func (p *Parser) registerInfixFn(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func (p *Parser) registerPostfixFn(tokenType token.TokenType, fn postfixParseFn) {
	p.postfixParseFns[tokenType] = fn
}

func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}

	return LOWEST
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	stmt := &ast.LetStatement{Token: p.curToken}

	if !p.assertNext(token.IDENT) {
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.assertNext(token.ASSIGN) {
		return nil
	}

	//TODO: skipping the expressions for now. Jump to the semicolon
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{Token: p.curToken}

	p.nextToken()

	//TODO: skipping the expressions for now. Jump to the semicolon
	for p.curToken.Type != token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekToken.Type == token.SEMICOLON {
		p.nextToken()
	}

	return stmt
}

func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for (p.peekToken.Type != token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()

		leftExp = infix(leftExp)
	}

	return leftExp
}

func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	tempInt := big.Int{}

	_, ok := tempInt.SetString(p.curToken.Literal, 10)
	if !ok {
		msg := fmt.Sprintf("integer token %q at %s is an invalid int", p.curToken.Literal, p.peekToken.PositionString())
		p.errors = append(p.errors, msg)
	}

	lit.Value = tempInt

	return lit
}

func (p *Parser) parseDecimalLiteral() ast.Expression {
	lit := &ast.DecimalLiteral{Token: p.curToken}

	tempFloat := big.Float{}

	_, ok := tempFloat.SetString(p.curToken.Literal)
	if !ok {
		msg := fmt.Sprintf("decimal token %q at %s is an invalid decimal", p.curToken.Literal, p.peekToken.PositionString())
		p.errors = append(p.errors, msg)
	}

	lit.Value = tempFloat

	return lit
}

func (p *Parser) parseBoolean() ast.Expression {
	return &ast.Boolean{
		Token: p.curToken,
		Value: p.curToken.Type == token.TRUE,
	}
}

func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.assertNext(token.RPAREN) {
		return nil
	}

	return exp
}

func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

func (p *Parser) parsePostfixExpression(left ast.Expression) ast.Expression {
	postfix := p.postfixParseFns[p.peekToken.Type]
	if postfix == nil {
		return &ast.PostfixExpression{
			Token:    p.curToken,
			Operator: p.curToken.Literal,
			Left:     left,
		}
	}

	p.nextToken()

	left = postfix(left)

	return left
}

func (p *Parser) assertNext(t token.TokenType) bool {
	if p.peekToken.Type == t {
		p.nextToken()
		return true
	} else {
		p.peekError(t)
		return false
	}
}

func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token at %s to be a %s, got %s instead", p.peekToken.PositionString(), t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no prefix parse function found for %s", t)
	p.errors = append(p.errors, msg)
}

func (p *Parser) noPostfixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("no postfix parse function found for %s", t)
	p.errors = append(p.errors, msg)
}
