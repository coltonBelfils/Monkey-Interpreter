package object

import (
	"fmt"
	"math/big"
)

type ObjectType string

const (
	INTEGER_OBJ ObjectType = "INTEGER"
	FLOAT_OBJ ObjectType = "FLOAT"
	BOOLEAN_OBJ ObjectType = "BOOLEAN"
	NULL_OBJ ObjectType = "NULL"
	POINTER_OBJ ObjectType = "POINTER"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value big.Int
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%s", i.Value.String())
}

func (i *Integer) Type() ObjectType {
	return INTEGER_OBJ
}

type Float struct {
	Value big.Float
}

func (f *Float) Inspect() string {
	return fmt.Sprintf("%s", f.Value.String())
}

func (f *Float) Type() ObjectType {
	return FLOAT_OBJ
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

type Pointer [T Object]struct {
	Value T
}

func (p *Pointer[T]) Inspect() string {
	return fmt.Sprintf("*%s", p.Value.Inspect())
}

func (p *Pointer[T]) Type() ObjectType {
	return ObjectType(fmt.Sprintf("%s of %s", POINTER_OBJ, p.Value.Type()))
}