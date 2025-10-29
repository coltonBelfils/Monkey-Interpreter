package token

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

type Token struct {
	Type     TokenType
	Literal  string
	File     string
	Line     int
	Position int
}
