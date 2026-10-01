package ast

import (
	"MonkeyInterpreter/token"
	"math/big"
	"strings"
)

type Node interface {
	TokenLiteral() string
	String() string
}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	} else {
		return ""
	}
}

func (p *Program) String() string {
	var out strings.Builder

	for _, s := range p.Statements {
		out.WriteString(s.String())
	}

	return strings.TrimSpace(out.String())
}

type LetStatement struct {
	Token token.Token
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode() {}
func (ls *LetStatement) TokenLiteral() string {
	return ls.Token.Literal
}

func (ls *LetStatement) String() string {
	var out strings.Builder

	out.WriteString(ls.TokenLiteral())
	out.WriteString(" ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")

	/* //??? Why is this what he chose?
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	} else {
		out.WriteString("___")
	}
	*/

	out.WriteString(ls.Value.String())

	out.WriteString(";\n")

	return out.String()
}

type ReturnStatement struct {
	Token       token.Token
	ReturnValue Expression
}

func (rs *ReturnStatement) statementNode() {}
func (rs *ReturnStatement) TokenLiteral() string {
	return rs.Token.Literal
}

func (rs *ReturnStatement) String() string {
	var out strings.Builder

	out.WriteString(rs.TokenLiteral())
	out.WriteString(" ")

	if rs.ReturnValue != nil {
		out.WriteString(rs.ReturnValue.String())
	} else {
		out.WriteString("___")
	}

	out.WriteString(";\n")

	return out.String()
}

type ExpressionStatement struct {
	Token      token.Token
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}
func (es *ExpressionStatement) TokenLiteral() string {
	return es.Token.Literal
}

func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String() + ";\n"
	}

	return ""
}

type Identifier struct {
	Token token.Token
	Value string
}

func (i *Identifier) expressionNode() {}
func (i *Identifier) TokenLiteral() string {
	return i.Token.Literal
}

func (i *Identifier) String() string {
	return i.Value
}

type IntegerLiteral struct {
	Token token.Token
	Value big.Int
}

func (il *IntegerLiteral) expressionNode() {}
func (il *IntegerLiteral) TokenLiteral() string {
	return il.Token.Literal
}
func (il *IntegerLiteral) String() string {
	return il.Token.Literal
}

type DecimalLiteral struct {
	Token token.Token
	Value big.Float
}

func (dl *DecimalLiteral) expressionNode() {}
func (dl *DecimalLiteral) TokenLiteral() string {
	return dl.Token.Literal
}
func (dl *DecimalLiteral) String() string {
	return dl.Token.Literal
}

type Boolean struct {
	Token token.Token
	Value bool
}

func (b *Boolean) expressionNode() {}
func (b *Boolean) TokenLiteral() string {
	return b.Token.Literal
}
func (b *Boolean) String() string {
	return b.Token.Literal
}

type PrefixExpression struct {
	Token    token.Token
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode() {}
func (pe *PrefixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PrefixExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")

	return out.String()
}

type InfixExpression struct {
	Token    token.Token
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}
func (ie *InfixExpression) TokenLiteral() string {
	return ie.Token.Literal
}
func (ie *InfixExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" ")
	out.WriteString(ie.Operator)
	out.WriteString(" ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")

	return out.String()
}

type PostfixExpression struct {
	Token    token.Token
	Operator string
	Left     Expression
}

func (pe *PostfixExpression) expressionNode() {}
func (pe *PostfixExpression) TokenLiteral() string {
	return pe.Token.Literal
}
func (pe *PostfixExpression) String() string {
	var out strings.Builder

	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(pe.Operator)
	out.WriteString(")")

	return out.String()
}

type IfExpression struct {
	Token       token.Token
	Condition   Expression
	Consequence Expression
	Alternative Expression
}

func (ie *IfExpression) expressionNode() {}

func (ie *IfExpression) TokenLiteral() string {
	return ie.Token.Literal
}

func (ie *IfExpression) String() string {
	var out strings.Builder

	out.WriteString("if(")
	out.WriteString(ie.Condition.String())
	out.WriteString(")")
	if ie.Consequence != nil {
		out.WriteString(" ")
		out.WriteString(ie.Consequence.String())

		if ie.Alternative != nil {
			out.WriteString(" else ")
			out.WriteString(ie.Alternative.String())
		}
	}

	return out.String()
}

/*
type BlockStatement struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

func (bs *BlockStatement) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockStatement) String() string {
	var out strings.Builder

	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}

	return out.String()
}
*/

type BlockExpression struct {
	Token      token.Token
	Statements []Statement
}

func (bs *BlockExpression) expressionNode() {}

func (bs *BlockExpression) TokenLiteral() string {
	return bs.Token.Literal
}

func (bs *BlockExpression) String() string {
	var out strings.Builder

	if len(bs.Statements) == 0 {
		return "{}"
	}

	out.WriteString("{\n\t")
	for i, s := range bs.Statements {
		str := strings.ReplaceAll(s.String(), "\n", "\n\t")
		if i == len(bs.Statements)-1 {
			str = strings.TrimSuffix(str, "\t")
		}
		out.WriteString(str)
	}
	out.WriteString("}")

	return out.String()
}

type FunctionLiteral struct {
	Token      token.Token
	Parameters []Identifier //expressions or identifiers?
	Body       Expression
}

func (fl *FunctionLiteral) expressionNode() {}

func (fl *FunctionLiteral) TokenLiteral() string {
	return fl.Token.Literal
}

func (fl *FunctionLiteral) String() string {
	var out strings.Builder

	out.WriteString("fn(")
	for i, param := range fl.Parameters {
		out.WriteString(param.String())
		if i != len(fl.Parameters)-1 {
			out.WriteString(", ")
		}
	}
	out.WriteString(") ")
	out.WriteString(fl.Body.String())

	return out.String()
}

type CallExpression struct {
	Token token.Token
	Function Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

func (ce *CallExpression) TokenLiteral() string {
	return ce.Token.Literal
}

func (ce *CallExpression) String() string {
	out := strings.Builder{}

	out.WriteString(ce.Function.String())

	out.WriteString("(")

	for i, a := range ce.Arguments {
		out.WriteString(a.String())
		if i != len(ce.Arguments) - 1 {
			out.WriteString(", ")
		}
	}

	out.WriteString(")")

	return out.String()
}