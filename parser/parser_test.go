package parser

import (
	"MonkeyInterpreter/ast"
	"MonkeyInterpreter/lexer"
	"MonkeyInterpreter/token"
	"fmt"
	"math/big"
	"os"
	"strings"
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
	input := `return 5;
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

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.NewFromString(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program had unexpected number of statements. Expected: 1. Got: %d", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. Got: %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("expression not *ast.Identifier. Got: %T", stmt.Expression)
	}
	if ident.Value != "foobar" {
		t.Errorf("ident.Value incorrect. Got: %q. Expected: %q", ident.Value, "foobar")
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral() incorrect. Got: %q. Expected: %q", ident.Value, "foobar")
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.NewFromString(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program had unexpected number of statements. Got: %d. Expected: 1", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. Got: %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("expression not *ast.IntegerLiteral. Got: %T", stmt.Expression)
	}
	if ident.Value.String() != "5" {
		t.Errorf("ident.Value incorrect. Got: %s. Expected: 5", ident.Value.String())
	}
	if ident.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral() incorrect. Got: %q. Expected: %q", ident.Value.String(), "5")
	}
}

func TestDecimalLiteralExpression(t *testing.T) {
	input := "5.5;"

	l := lexer.NewFromString(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program had unexpected number of statements. Got: %d. Expected: 1", len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. Got: %T", program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.DecimalLiteral)
	if !ok {
		t.Fatalf("expression not *ast.DecimalLiteral. Got: %T", stmt.Expression)
	}
	if ident.Value.String() != "5.5" {
		t.Errorf("ident.Value incorrect. Got: %s. Expected: 5.5", ident.Value.String())
	}
	if ident.TokenLiteral() != "5.5" {
		t.Errorf("ident.TokenLiteral() incorrect. Got: %q. Expected: %q", ident.Value.String(), "5.5")
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", big.NewInt(5)},
		{"-15.7;", "-", big.NewFloat(15.7)},
		{"++34;", "++", big.NewInt(34)},
		{"--88.9;", "--", big.NewFloat(88.9)},
		{"--88.9;", "--", big.NewFloat(88.9)},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		lPrint := lexer.NewFromString(tt.input)
		var tokens strings.Builder
		for t := lPrint.NextToken(); t.Type != token.EOF; t = lPrint.NextToken() {
			tokens.WriteString(t.String())
			tokens.WriteString(" - ")
			tokens.WriteString(t.PositionString())
			tokens.WriteString("\n")
		}
		fmt.Printf("\nTokens:\n%s", tokens.String())

		l := lexer.NewFromString(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		fmt.Printf("%s\n", program.String())

		if len(program.Statements) != 1 {
			t.Fatalf("program had unexpected number of statements. Got: %d. Expected: 1", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. Got: %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.value) {
			return
		}
	}
}

func TestParsingPostfixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		left     interface{}
	}{
		{"5++;", "++", big.NewInt(5)},
		{"44.1--;", "--", big.NewFloat(44.1)},
		{"34!;", "!", big.NewInt(34)},
		{"34!;", "!", big.NewInt(34)},
	}

	for _, tt := range prefixTests {
		lPrint := lexer.NewFromString(tt.input)
		var tokens strings.Builder
		for t := lPrint.NextToken(); t.Type != token.EOF; t = lPrint.NextToken() {
			tokens.WriteString(t.String())
			tokens.WriteString(" - ")
			tokens.WriteString(t.PositionString())
			tokens.WriteString("\n")
		}
		fmt.Printf("Tokens:\n%s", tokens.String())

		l := lexer.NewFromString(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		fmt.Printf("%s\n", program.String())

		if len(program.Statements) != 1 {
			t.Fatalf("program had unexpected number of statements. Got: %d. Expected: 1", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. Got: %T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PostfixExpression)
		if !ok {
			t.Fatalf("expression not *ast.PostfixExpression. Got: %T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not %q. Got: %s", tt.operator, exp.Operator)
		}

		testPostfixExpression(t, stmt.Expression, tt.left, tt.operator)
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input    string
		left     interface{}
		operator string
		right    interface{}
	}{
		{"5 + 5;", big.NewInt(5), "+", big.NewInt(5)},
		{"5 - 5;", big.NewInt(5), "-", big.NewInt(5)},
		{"5 * 5;", big.NewInt(5), "*", big.NewInt(5)},
		{"5 / 5;", big.NewInt(5), "/", big.NewInt(5)},
		{"5 % 5;", big.NewInt(5), "%", big.NewInt(5)},
		{"5 > 5;", big.NewInt(5), ">", big.NewInt(5)},
		{"5 < 5;", big.NewInt(5), "<", big.NewInt(5)},
		{"5 == 5;", big.NewInt(5), "==", big.NewInt(5)},
		{"5 != 5;", big.NewInt(5), "!=", big.NewInt(5)},
		{"5 <= 5;", big.NewInt(5), "<=", big.NewInt(5)},
		{"5 >= 5;", big.NewInt(5), ">=", big.NewInt(5)},
		{"5 && 5;", big.NewInt(5), "&&", big.NewInt(5)},
		{"5 || 5;", big.NewInt(5), "||", big.NewInt(5)},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		lPrint := lexer.NewFromString(tt.input)
		var tokens strings.Builder
		for t := lPrint.NextToken(); t.Type != token.EOF; t = lPrint.NextToken() {
			tokens.WriteString(t.String())
			tokens.WriteString(" - ")
			tokens.WriteString(t.PositionString())
			tokens.WriteString("\n")
		}
		fmt.Printf("Tokens:\n%s", tokens.String())

		l := lexer.NewFromString(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		fmt.Printf("%s\n", program.String())

		if len(program.Statements) != 1 {
			t.Fatalf("program had unexpected number of statements. Got: %d. Expected: 1", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not an ast.ExpressionStatement. Got: %T", program.Statements[0])
		}

		if !testInfixExpression(t, stmt.Expression, tt.left,
			tt.operator, tt.right) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"1 + 1",
			"(1 + 1)",
		},
		{
			"1 + a * cow",
			"(1 + (a * cow))",
		},
		{
			"!a++",
			"(!(a++))",
		},
		{
			"!a++ + r",
			"((!(a++)) + r)",
		},
		{
			"3 + 4; -5 * 5++",
			"(3 + 4)((-5) * (5++))",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a--",
			"(!(-(a--)))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b++ * c",
			"((a * (b++)) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}

	for _, tt := range tests {
		l := lexer.NewFromString(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q. got=%q", tt.expected, actual)
		}
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input           string
		expectedBoolean bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		l := lexer.NewFromString(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.Boolean)
		if !ok {
			t.Fatalf("exp not *ast.Boolean. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedBoolean {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedBoolean,
				boolean.Value)
		}
	}
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value *big.Int) bool {
	i, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if i.Value.String() != value.String() {
		t.Errorf("i.value not %s. got=%s", value.String(), i.Value.String())
		return false
	}

	if i.TokenLiteral() != fmt.Sprintf("%s", value.String()) {
		t.Errorf("i.TokenLiteral() not %s. got=%s", value.String(), i.TokenLiteral())
		return false
	}

	return true
}

func testDecimalLiteral(t *testing.T, fl ast.Expression, value *big.Float) bool {
	f, ok := fl.(*ast.DecimalLiteral)
	if !ok {
		t.Errorf("fl not *ast.DecimalLiteral. got=%T", fl)
		return false
	}

	if f.Value.String() != value.String() {
		t.Errorf("f.value not %s. got=%s", value.String(), f.Value.String())
		return false
	}

	if f.TokenLiteral() != fmt.Sprintf("%s", value.String()) {
		t.Errorf("f.TokenLiteral() not %s. got=%s", value.String(), f.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, bl ast.Expression, value bool) bool {
	b, ok := bl.(*ast.Boolean)
	if !ok {
		t.Errorf("bl not *ast.Boolean. got=%T", bl)
		return false
	}

	if b.Value != value {
		t.Errorf("b.value not %t. got=%t", value, b.Value)
		return false
	}

	if b.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("b.TokenLiteral() not %t. got=%s", value, b.TokenLiteral())
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral() not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case *big.Int:
		return testIntegerLiteral(t, exp, v)
	case *big.Float:
		return testDecimalLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}

	t.Errorf("type of exp not handled. expected type=%T. got=%T", expected, exp)
	return false
}

func testPostfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string) bool {
	opExp, ok := exp.(*ast.PostfixExpression)
	if !ok {
		t.Errorf("exp is not ast.PostfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %q. got=%q", operator, opExp.Operator)
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not %q. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
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
