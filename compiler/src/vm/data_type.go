package vm

import (
	"fmt"
)

type Integer struct {
	Value int64
}

func (integer *Integer) Type() string {
	return "integer"
}
func (integer *Integer) Inspect() string {
	return fmt.Sprintf("%d", integer.Value)
}

type Boolean struct {
	Value bool
}

func (boolean *Boolean) Type() string {
	return "boolean"
}
func (boolean *Boolean) Inspect() string {
	return fmt.Sprintf("%t", boolean.Value)
}
