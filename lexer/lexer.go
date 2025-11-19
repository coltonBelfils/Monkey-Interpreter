package lexer

import (
	"MonkeyInterpreter/token"
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

const STRING_INPUT = "STRING_INPUT"
const REPL_INPUT = "REPL_INPUT"

type Lexer struct {
	input *bufio.Reader
	cur   rune
	next  rune //having cur and next may not be necessary since we can peek with the bufio.Reader

	//file location tracking for error reporting
	line         int
	linePosition int
	fileName     string
}

func New(input io.Reader, fileName string) *Lexer {
	var bufioReader *bufio.Reader
	if br, ok := input.(*bufio.Reader); ok {
		bufioReader = br
	} else {
		bufioReader = bufio.NewReader(input)
	}

	l := &Lexer{
		input:        bufioReader,
		cur:          rune(0),
		next:         rune(0),
		line:         1,
		linePosition: -1, //After the double l.readChar() call below, linePosition will be 1 and cur will also be pointing at the first char
		fileName:     fileName,
	}

	//Priming the pump
	l.readChar() //l.next has data
	l.readChar() //l.next and l.cur have data

	return l
}

func NewFromFile(file *os.File) *Lexer {
	return New(file, file.Name())
}

func NewFromString(input string) *Lexer {
	return New(strings.NewReader(input), STRING_INPUT)
}

func NewFromRepl(input io.Reader) *Lexer {
	return New(input, REPL_INPUT)
}

func (l *Lexer) readChar() {
	newRune, _, err := l.input.ReadRune()
	if err != nil {
		if errors.Is(err, io.EOF) {
			newRune = rune(0)
		} else {
			fmtErr := fmt.Errorf("an error occurred while reading the source input at %s:%d:%d. Error given: %w", l.fileName, l.line, l.linePosition, err)
			panic(fmtErr)
		}
	}

	l.cur = l.next
	l.next = newRune

	if l.cur == '\n' {
		l.line++
		l.linePosition = 0
	} else {
		l.linePosition += 1
	}
}

func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	switch l.cur {
	case '=':
		if l.next == '=' {
			tok = newToken(token.EQ, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else {
			tok = newToken(token.ASSIGN, l.cur, l.fileName, l.line, l.linePosition)
		}
	case ';':
		tok = newToken(token.SEMICOLON, l.cur, l.fileName, l.line, l.linePosition)
	case '\\': //This feels a little weird. "\ a" is different from "\a". White space will be discarded at the parsing level so should \a be encoded as a single token?
		tok = newToken(token.ESCAPE, l.cur, l.fileName, l.line, l.linePosition)
	case '(':
		tok = newToken(token.LPAREN, l.cur, l.fileName, l.line, l.linePosition)
	case ')':
		tok = newToken(token.RPAREN, l.cur, l.fileName, l.line, l.linePosition)
	case ',':
		tok = newToken(token.COMMA, l.cur, l.fileName, l.line, l.linePosition)
	case '+':
		if l.next == '=' {
			tok = newToken(token.PLUS_ASSIGN, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else if l.next == '+' {
			tok = newToken(token.PLUS_PLUS, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else {
			tok = newToken(token.PLUS, l.cur, l.fileName, l.line, l.linePosition)
		}
	case '-':
		if l.next == '=' {
			tok = newToken(token.MINUS_ASSIGN, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else if l.next == '-' {
			tok = newToken(token.MINUS_MINUS, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else {
			tok = newToken(token.MINUS, l.cur, l.fileName, l.line, l.linePosition)
		}
	case '!':
		if l.next == '=' {
			tok = newToken(token.NEQ, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else {
			tok = newToken(token.BANG, l.cur, l.fileName, l.line, l.linePosition)
		}
	case '/':
		if l.next == '=' {
			tok = newToken(token.DIV_ASSIGN, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else {
			tok = newToken(token.SLASH, l.cur, l.fileName, l.line, l.linePosition)
		}
	case '*':
		if l.next == '=' {
			tok = newToken(token.MULT_ASSIGN, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else {
			tok = newToken(token.ASTERISK, l.cur, l.fileName, l.line, l.linePosition)
		}
	case '<':
		if l.next == '=' {
			tok = newToken(token.LT_EQ, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else {
			tok = newToken(token.LT, l.cur, l.fileName, l.line, l.linePosition)
		}
	case '>':
		if l.next == '=' {
			tok = newToken(token.GT_EQ, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else {
			tok = newToken(token.GT, l.cur, l.fileName, l.line, l.linePosition)
		}
	case '{':
		tok = newToken(token.LBRACE, l.cur, l.fileName, l.line, l.linePosition)
	case '}':
		tok = newToken(token.RBRACE, l.cur, l.fileName, l.line, l.linePosition)
	case '&':
		if l.next == '&' {
			tok = newToken(token.AND, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else {
			tok = newToken(token.ILLEGAL, l.cur, l.fileName, l.line, l.linePosition)
		}
	case '|':
		if l.next == '|' {
			tok = newToken(token.OR, l.cur, l.fileName, l.line, l.linePosition)
			l.readChar()
			tok.Literal += string(l.cur)
		} else {
			tok = newToken(token.ILLEGAL, l.cur, l.fileName, l.line, l.linePosition)
		}
	case rune(0):
		tok = newToken(token.EOF, rune(0), l.fileName, l.line, l.linePosition)
	default:
		if isLetter(l.cur) {
			tok.FileName = l.fileName
			tok.LinePosition = l.line
			tok.Position = l.linePosition

			tok.Literal = l.readIdentifier()

			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.cur) {
			tok.FileName = l.fileName
			tok.LinePosition = l.line
			tok.Position = l.linePosition

			tok.Literal = l.readNumber()

			tok.Type = token.INT

			if l.cur == '.' && isDigit(l.next) {
				tok.Literal += "."
				l.readChar()
				tok.Literal += l.readNumber()
				tok.Type = token.DEC
			}

			return tok
		} else {
			tok = newToken(token.ILLEGAL, l.cur, l.fileName, l.line, l.linePosition)
		}
	}

	l.readChar()
	return tok
}

func (l *Lexer) readIdentifier() string {
	var ident string
	for isLetter(l.cur) {
		ident += string(l.cur)
		l.readChar()
	}

	return ident
}

func (l *Lexer) readNumber() string {
	var ident string
	for isDigit(l.cur) {
		ident += string(l.cur)
		l.readChar()
	}

	return ident
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.cur) {
		l.readChar()
	}
}

func newToken(tokenType token.TokenType, r rune, fileName string, line int, linePosition int) token.Token {
	lit := string(r)
	if r == rune(0) {
		lit = ""
	}

	return token.Token{
		Type:         tokenType,
		Literal:      lit,
		FileName:     fileName,
		LinePosition: line,
		Position:     linePosition,
	}
}

func isLetter(r rune) bool {
	return unicode.IsLetter(r) || r == '_'
}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}
