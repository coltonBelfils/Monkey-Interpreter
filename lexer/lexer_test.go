package lexer

import (
	"MonkeyInterpreter/token"
	"os"
	"testing"
)

func TestNextToken(t *testing.T) {
	f, err := os.Open("./testfile.monkey")
	if err != nil {
		t.Fatalf("could not open test file testfile.monkey. Err: %v", err)
	}
	defer f.Close()

	tests := []struct {
		expectedType     token.TokenType
		expectedLiteral  string
		expectedFile     string
		expectedLine     int
		expectedPosition int
	}{
		{token.LET, "let", "./testfile.monkey", 1, 1},
		{token.IDENT, "five", "./testfile.monkey", 1, 5},
		{token.ASSIGN, "=", "./testfile.monkey", 1, 10},
		{token.INT, "5", "./testfile.monkey", 1, 12},
		{token.SEMICOLON, ";", "./testfile.monkey", 1, 13},

		{token.LET, "let", "./testfile.monkey", 2, 1},
		{token.IDENT, "ten", "./testfile.monkey", 2, 5},
		{token.ASSIGN, "=", "./testfile.monkey", 2, 9},
		{token.INT, "10", "./testfile.monkey", 2, 11},
		{token.SEMICOLON, ";", "./testfile.monkey", 2, 13},

		{token.LET, "let", "./testfile.monkey", 4, 1},
		{token.IDENT, "add", "./testfile.monkey", 4, 5},
		{token.ASSIGN, "=", "./testfile.monkey", 4, 9},
		{token.FUNCTION, "fn", "./testfile.monkey", 4, 11},
		{token.LPAREN, "(", "./testfile.monkey", 4, 13},
		{token.IDENT, "x", "./testfile.monkey", 4, 14},
		{token.COMMA, ",", "./testfile.monkey", 4, 15},
		{token.IDENT, "y", "./testfile.monkey", 4, 17},
		{token.RPAREN, ")", "./testfile.monkey", 4, 18},
		{token.LBRACE, "{", "./testfile.monkey", 4, 20},
		{token.IDENT, "x", "./testfile.monkey", 5, 2},
		{token.PLUS, "+", "./testfile.monkey", 5, 4},
		{token.IDENT, "y", "./testfile.monkey", 5, 6},
		{token.SEMICOLON, ";", "./testfile.monkey", 5, 7},

		{token.RBRACE, "}", "./testfile.monkey", 6, 1},
		{token.SEMICOLON, ";", "./testfile.monkey", 6, 2},

		{token.LET, "let", "./testfile.monkey", 8, 1},
		{token.IDENT, "result", "./testfile.monkey", 8, 5},
		{token.ASSIGN, "=", "./testfile.monkey", 8, 12},
		{token.IDENT, "add", "./testfile.monkey", 8, 14},
		{token.LPAREN, "(", "./testfile.monkey", 8, 17},
		{token.IDENT, "five", "./testfile.monkey", 8, 18},
		{token.COMMA, ",", "./testfile.monkey", 8, 22},
		{token.IDENT, "ten", "./testfile.monkey", 8, 24},
		{token.RPAREN, ")", "./testfile.monkey", 8, 27},
		{token.SEMICOLON, ";", "./testfile.monkey", 8, 28},
		{token.EOF, "", "./testfile.monkey", 8, 29},
	}

	l := NewFromFile(f)

	for i, tt := range tests {

		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong. expected=%v, got=%v", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%v, got=%v", i, tt.expectedLiteral, tok.Literal)
		}

		if tok.FileName != tt.expectedFile {
			t.Fatalf("test[%d] - file wrong. expected=%v, got=%v", i, tt.expectedFile, tok.FileName)
		}

		if tok.LinePosition != tt.expectedLine {
			t.Fatalf("test[%d] - line wrong. expected=%v, got=%v", i, tt.expectedLine, tok.LinePosition)
		}

		if tok.Position != tt.expectedPosition {
			t.Fatalf("test[%d] - position wrong. expected=%v, got=%v", i, tt.expectedPosition, tok.Position)
		}
	}
}
