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

	//operators
	ASSIGN = "ASSIGN"
	PLUS   = "PLUS"

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
)

var keywords = map[string]TokenType{
	"fn": FUNCTION,
	"let": LET,
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
