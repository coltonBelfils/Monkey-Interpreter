package token

import "fmt"

type TokenType string

const (
	//meta tokens
	ILLEGAL = "ILLEGAL"
	EOF     = "EOF"

	//identifiers and literals (user generated content)
	IDENT = "IDENT"
	INT   = "INT"
	DEC = "DEC"

	ESCAPE = "ESCAPE"

	//operators
	ASSIGN = "ASSIGN"
	PLUS   = "PLUS"
	MINUS = "MINUS"
	BANG = "BANG"
	ASTERISK = "ASTERISK"
	SLASH = "SLASH"

	PLUS_PLUS = "PLUS_PLUS"
	MINUS_MINUS = "MINUS_MINUS"
	PLUS_ASSIGN = "PLUS_ASSIGN"
	MINUS_ASSIGN = "MINUS_ASSIGN"
	MULT_ASSIGN = "MULT_ASSIGN"
	DIV_ASSIGN = "DIV_ASSIGN"

	LT = "LT"
	GT = "GT"
	EQ = "EQ"
	NEQ = "NEQ"
	LT_EQ = "LT_EQ"
	GT_EQ = "GT_EQ"

	AND = "AND"
	OR = "OR"

	//delimiters
	COMMA     = "COMMA"
	SEMICOLON = "SEMICOLON"

	LPAREN = "LPAREN"
	RPAREN = "RPAREN"
	LBRACE = "LBRACE"
	RBRACE = "RBRACE"

	//keywords
	FUNCTION = "FUNCTION"
	LET      = "LET"
	TRUE = "TRUE"
	FALSE = "FALSE"
	RETURN = "RETURN"
	IF = "IF"
	ELSE = "ELSE"
	QUIT = "QUIT"
)

var keywords = map[string]TokenType{
	"fn": FUNCTION,
	"let": LET,
	"true": TRUE,
	"false": FALSE,
	"return": RETURN,
	"if": IF,
	"else": ELSE,
	"quit": QUIT,
}

func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

type Token struct {
	Type         TokenType
	Literal      string
	FileName     string
	LinePosition int
	Position     int
	PositionEnd  int//maybe?
}

func (t Token) String() string {
	return fmt.Sprintf("%s of type %s at %s:%d:%d", t.Literal, t.Type, t.FileName, t.LinePosition, t.Position)
}
