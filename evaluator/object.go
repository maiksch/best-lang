package evaluator

import "fmt"

type ObjectType string

const (
	INTEGER ObjectType = "INTEGER"
	BOOLEAN ObjectType = "BOOLEAN"
	RETURN  ObjectType = "RETURN"
	ERROR   ObjectType = "ERROR"
	NOTHING ObjectType = "NOTHING"
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

func (i *Integer) Type() ObjectType { return INTEGER }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }

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

func (b *Boolean) Type() ObjectType { return BOOLEAN }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }

/**
* Return
 */

type ReturnValue struct {
	Value Object
}

func (r *ReturnValue) Type() ObjectType { return RETURN }
func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }

/**
* Error
 */

type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }

/**
* Nothing
 */

var NOTHING_OBJ = &Nothing{}

type Nothing struct{}

func (n *Nothing) Type() ObjectType { return NOTHING }
func (n *Nothing) Inspect() string  { return "nothing" }
