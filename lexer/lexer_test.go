package lexer

import (
	"MonkeyInterpreter/token"
	"os"
	"testing"
)

func TestNextTokenWithPosition(t *testing.T) {
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

		{token.LET, "let", "./testfile.monkey", 10, 1},
		{token.IDENT, "dec", "./testfile.monkey", 10, 5},
		{token.ASSIGN, "=", "./testfile.monkey", 10, 9},
		{token.DEC, "12.4", "./testfile.monkey", 10, 11},
		{token.SEMICOLON, ";", "./testfile.monkey", 10, 15},

		{token.BANG, "!", "./testfile.monkey", 12, 1},
		{token.MINUS, "-", "./testfile.monkey", 12, 2},
		{token.SLASH, "/", "./testfile.monkey", 12, 3},
		{token.ASTERISK, "*", "./testfile.monkey", 12, 4},
		{token.INT, "5", "./testfile.monkey", 12, 5},
		{token.SEMICOLON, ";", "./testfile.monkey", 12, 6},

		{token.INT, "5", "./testfile.monkey", 13, 1},
		{token.LT, "<", "./testfile.monkey", 13, 3},
		{token.INT, "10", "./testfile.monkey", 13, 5},
		{token.GT, ">", "./testfile.monkey", 13, 8},
		{token.INT, "5", "./testfile.monkey", 13, 10},
		{token.SEMICOLON, ";", "./testfile.monkey", 13, 11},

		{token.EOF, "", "./testfile.monkey", 13, 12},
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

func TestNextTokenNoPosition(t *testing.T) {
	input := `let five = 5;
let ten = 10;

let add = fn(x, y) {
	x + y;
};

let result = add(five, ten);

let dec = 12.4;

!-/*5;
5 < 10 > 5;

if (5 < 10) {
    return true;
} else {
    return false;
}

== != >= <= && || ++ -- += -= *= /=`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.LET, "let"},
		{token.IDENT, "five"},
		{token.ASSIGN, "="},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "ten"},
		{token.ASSIGN, "="},
		{token.INT, "10"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNCTION, "fn"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.COMMA, ","},
		{token.IDENT, "y"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "x"},
		{token.PLUS, "+"},
		{token.IDENT, "y"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "result"},
		{token.ASSIGN, "="},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.IDENT, "five"},
		{token.COMMA, ","},
		{token.IDENT, "ten"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},

		{token.LET, "let"},
		{token.IDENT, "dec"},
		{token.ASSIGN, "="},
		{token.DEC, "12.4"},
		{token.SEMICOLON, ";"},

		{token.BANG, "!"},
		{token.MINUS, "-"},
		{token.SLASH, "/"},
		{token.ASTERISK, "*"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.GT, ">"},
		{token.INT, "5"},
		{token.SEMICOLON, ";"},

		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.LT, "<"},
		{token.INT, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},

		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},

		{token.RETURN, "return"},
		{token.FALSE, "false"},
		{token.SEMICOLON, ";"},

		{token.RBRACE, "}"},

		{token.EQ, "=="},
		{token.NEQ, "!="},
		{token.GT_EQ, ">="},
		{token.LT_EQ, "<="},
		{token.AND, "&&"},
		{token.OR, "||"},
		{token.PLUS_PLUS, "++"},
		{token.MINUS_MINUS, "--"},
		{token.PLUS_ASSIGN, "+="},
		{token.MINUS_ASSIGN, "-="},
		{token.MULT_ASSIGN, "*="},
		{token.DIV_ASSIGN, "/="},

		{token.EOF, ""},
	}

	l := NewFromString(input)

	for i, tt := range tests {

		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("test[%d] - tokentype wrong. expected=%v, got=%v. Token: %v", i, tt.expectedType, tok.Type, tok)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("test[%d] - literal wrong. expected=%v, got=%v. Token: %v", i, tt.expectedLiteral, tok.Literal, tok)
		}
	}
}
