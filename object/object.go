package object

import (
	"MonkeyInterpreter/ast"
	"fmt"
)

type ObjectType string

const (
	INTEGER_OBJ ObjectType = "INTEGER"
	DECIMAL_OBJ ObjectType = "DECIMAL"
	BOOLEAN_OBJ ObjectType = "BOOLEAN"
	NULL_OBJ ObjectType = "NULL"
	POINTER_OBJ ObjectType = "POINTER"
	RETURN_VALUE_OBJ ObjectType = "RETURN_VALUE"
	FULL_RETURN_VALUE_OBJ ObjectType = "FULL_RETURN_VALUE_OBJ"
	UNWIND_VALUE_OBJ ObjectType = "UNWIND_VALUE_OBJ"
	ERROR_OBJ ObjectType = "ERROR"
)

type Object interface {
	Type() ObjectType
	Inspect() string
	//PassByReference() bool //pass by value or pass by reference flag? Default: by value. This would replace the pointer type.
}

type Integer struct {
	Value int
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

type Decimal struct {
	Value float64
}

func (d *Decimal) Inspect() string {
	return fmt.Sprintf("%g", d.Value)
}

func (d *Decimal) Type() ObjectType {
	return DECIMAL_OBJ
}

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN_OBJ
}

type Null struct {}

func (n *Null) Inspect() string {
	return "[null]"
}

func (n *Null) Type() ObjectType {
	return NULL_OBJ
}

//This strictly turns the object from the default pass by value to pass by reference. That's it.
type Pointer [T Object]struct {
	Value T
}

func (p *Pointer[T]) Inspect() string {
	return fmt.Sprintf("*%s", p.Value.Inspect())
}

func (p *Pointer[T]) Type() ObjectType {
	return ObjectType(fmt.Sprintf("%s of %s", POINTER_OBJ, p.Value.Type()))
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType {
	return RETURN_VALUE_OBJ
}

func (rv *ReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type FullReturnValue struct {
	Value Object
}

func (rv *FullReturnValue) Type() ObjectType {
	return FULL_RETURN_VALUE_OBJ
}

func (rv *FullReturnValue) Inspect() string {
	return rv.Value.Inspect()
}

type UnwindValue struct {
	Value Object
	UnwindTo ast.Node
}

func (uw *UnwindValue) Type() ObjectType {
	return UNWIND_VALUE_OBJ
}

func (uw *UnwindValue) Inspect() string {
	return uw.Value.Inspect()
}

type Error struct {
	Message string
	Position string
}

func (e *Error) Type() ObjectType {
	return ERROR_OBJ
}

func (e *Error) Inspect() string {
	return fmt.Sprintf("ERROR: %s. Near %s", e.Message, e.Position)
}

func NewError(at string, format string, a ...interface{}) *Error {
	return &Error{
		Message:  fmt.Sprintf(format, a...),
		Position: at,
	}
}