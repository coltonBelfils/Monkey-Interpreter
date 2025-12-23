package parser

import (
	"MonkeyInterpreter/ast"
	"MonkeyInterpreter/lexer"
	"os"
	"testing"
)

func TestLetStatement(t *testing.T) {
	//	input := `
	//let x 5;
	//let y = 10;
	//let foobar = 838383;
	//`
	//l := lexer.NewFromString(input)

	f, err := os.Open("testfile.monkey")
	if err != nil {
		t.Fatalf("could not open test file testfile.monkey. Err: %v", err)
	}
	defer f.Close()
	l := lexer.NewFromFile(f)
	p := New(l) //parser.New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statments does not contain 3 statements. Got: %d", len(program.Statements))
	}

	tests := []struct {
		expectedIdent string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range tests {
		stmt := program.Statements[i]
		if !testLetStatement(t, stmt, test.expectedIdent) {
			return
		}
	}
}

func testLetStatement(t *testing.T, s ast.Statement, name string) bool {
	if s.TokenLiteral() != "let" {
		t.Errorf("s.TokenLiteral not 'let'. Got: %q", s.TokenLiteral())
		return false
	}

	letStmt, ok := s.(*ast.LetStatement)
	if !ok {
		t.Errorf("s not *ast.LetStatement. Got: %T", s)
		return false
	}

	if letStmt.Name.Value != name {
		t.Errorf("letStmt.Name.Value not '%s'. Got: %s", name, letStmt.Name.Value)
		return false
	}

	if letStmt.Name.TokenLiteral() != name {
		t.Errorf("letStmt.Name.TokenLiteral() not '%s'. Got: %s", name, letStmt.Name.TokenLiteral())
		return false
	}

	return true
}

func TestReturnStatement(t *testing.T) {
	input := `return 5;'
return 10;
return 993322;`

	l := lexer.NewFromString(input)
	p := New(l) //parser.New(l)

	program := p.ParseProgram()

	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("program.Statments does not contain 3 statements. Got: %d", len(program.Statements))
	}

	for _, stmt := range program.Statements {
		retStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. Got: %T", stmt)
			continue
		}
		if retStmt.TokenLiteral() != "return" {
			t.Errorf("returnStmt.TokenLiteral not 'return', Got: %q", retStmt.TokenLiteral())
		}
	}
}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, e := range errors {
		t.Errorf("parser error: %q", e)
	}
	t.FailNow()
}
