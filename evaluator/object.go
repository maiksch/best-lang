package evaluator

import "fmt"

type ObjectType string

const (
	INTEGER ObjectType = "INTEGER"
	BOOLEAN ObjectType = "BOOLEAN"
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

/**
* Integers
 */

type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType {
	return INTEGER
}

func (i *Integer) Inspect() string {
	return fmt.Sprintf("%d", i.Value)
}

/**
* Boolean
 */

var (
	TRUE  = &Boolean{Value: true}
	FALSE = &Boolean{Value: false}
)

type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType {
	return BOOLEAN
}

func (b *Boolean) Inspect() string {
	return fmt.Sprintf("%t", b.Value)
}
