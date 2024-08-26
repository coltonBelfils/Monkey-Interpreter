package lexer

import (
	"MonkeyInterpreter/token"
)

type Lexer struct {
	input         string
	position      int
	readPosition  int
	ch            byte
	chunks        bool
	inMultComment bool
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar()
	return l
}

type ChunkLexer struct {
	Lexer
}

func (l *ChunkLexer) More(input string) {
	if len(input) == 0 {
		return
	}

	preLen := len(l.input)
	l.input += input
	if preLen == 0 {
		l.readChar()
	} else if l.position == preLen && l.position < len(l.input) {
		l.ch = l.input[l.position]
	}
}

func NewChunks() *ChunkLexer {
	l := &ChunkLexer{Lexer{chunks: true}}
	return l
}

func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0
	} else {
		l.ch = l.input[l.readPosition]
	}

	if l.position < len(l.input) {
		l.position = l.readPosition
		l.readPosition++
	}
}

func (l *Lexer) NextToken() token.Token {
	if l.chunks && l.inMultComment {
		if !l.eatToCommentTail() {
			l.inMultComment = true
			return token.Token{
				Type:    token.EOF,
				Literal: "",
			}
		} else {
			l.inMultComment = false
		}
		return l.NextToken()
	}

	l.eatWhiteSpace()

	var tok token.Token

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			a := l.ch
			l.readChar()
			b := l.ch
			tok = token.Token{Type: token.EQ, Literal: string([]byte{a, b})} // Should always produce ==, but taking in the given value anyway to better surface potential errors.
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}

	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '!':
		if l.peekChar() == '=' {
			a := l.ch
			l.readChar()
			b := l.ch
			tok = token.Token{Type: token.NEQ, Literal: string([]byte{a, b})} // Should always produce !=, but taking in the given value anyway to better surface potential errors.
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '/':
		if next := l.peekChar(); next == '/' {
			l.eatLine()
			return l.NextToken()
		} else if next == '*' {
			if !l.eatToCommentTail() {
				l.inMultComment = true
				return token.Token{
					Type:    token.EOF,
					Literal: "",
				}
			} else {
				l.inMultComment = false
			}
			return l.NextToken()
		} else {
			tok = newToken(token.SLASH, l.ch)
		}
	case '*':
		tok = newToken(token.ASTERISK, l.ch)
	case '%':
		tok = newToken(token.MOD, l.ch)
	case '>':
		tok = newToken(token.GT, l.ch)
	case '<':
		tok = newToken(token.LT, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case 0:
		tok = token.Token{
			Type:    token.EOF,
			Literal: "",
		}
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar()
	return tok
}

func newToken(tokenType token.TokenType, ch byte) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: string(ch),
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.position
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && 'z' >= ch || 'A' <= ch && 'Z' >= ch || '_' == ch
}

func (l *Lexer) eatWhiteSpace() {
	for isWhite(l.ch) {
		l.readChar()
	}
}

func (l *Lexer) eatLine() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

/*
eatToCommentTail eats characters in the input string until it finds '*' followed by '/', comment end. If found it'll eat those characters too.

Returns true if comment end was found, false if not found.
*/
func (l *Lexer) eatToCommentTail() bool {
	found := false
	for l.ch != '*' || l.peekChar() != '/' {
		if l.ch == 0 {
			break
		}
		l.readChar()
	}
	if l.ch == '*' && l.peekChar() == '/' {
		found = true
	}
	if found {
		l.readChar()
		l.readChar()
	}
	return found
}

func isWhite(ch byte) bool {
	return ' ' == ch || '\n' == ch || '\t' == ch || '\r' == ch
}

func (l *Lexer) readNumber() string {
	start := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[start:l.position]
}

func isDigit(ch byte) bool {
	return '0' <= ch && '9' >= ch
}

func (l *Lexer) peekChar() byte {
	if l.readPosition >= len(l.input) {
		return 0
	} else {
		return l.input[l.readPosition]
	}
}
